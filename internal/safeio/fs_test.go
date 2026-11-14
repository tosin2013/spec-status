package safeio

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestRejectsSiblingPrefixAndDotDot(t *testing.T) {
	parent := t.TempDir()
	root := filepath.Join(parent, "proj")
	evil := filepath.Join(parent, "proj-evil")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(evil, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(evil, "secret.txt"), []byte("TOPSECRET-VALUE"), 0o644); err != nil {
		t.Fatal(err)
	}
	fs, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := fs.Relativize(filepath.Join(evil, "secret.txt")); err == nil {
		t.Fatal("sibling path was accepted")
	}
	if _, err := fs.ReadRel("../proj-evil/secret.txt"); err == nil {
		t.Fatal("dot-dot read was accepted")
	}
}

func TestSymlinkComponentIsNotRead(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	secret := filepath.Join(outside, "payload.txt")
	if err := os.WriteFile(secret, []byte("TOPSECRET-VALUE"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(root, "specs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "specs", "leak")); err != nil {
		t.Fatal(err)
	}
	fs, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	_, err = fs.ReadRel("specs/leak/payload.txt")
	if err == nil || err.Error() != CodeUnsafe {
		t.Fatalf("symlink read err=%v", err)
	}
	_, err = fs.StatRel("specs/leak")
	if err == nil || err.Error() != CodeUnsafe {
		t.Fatalf("symlink stat err=%v", err)
	}
}

func TestSensitiveNameIsNotOpened(t *testing.T) {
	root := t.TempDir()
	fifo := filepath.Join(root, ".env")
	if err := syscall.Mkfifo(fifo, 0o644); err != nil {
		t.Fatal(err)
	}
	fs, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		_, err := fs.ReadRel(".env")
		done <- err
	}()
	select {
	case err := <-done:
		if err == nil || err.Error() != CodeSensitive {
			t.Fatalf("sensitive read err=%v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("sensitive read blocked, which means the FIFO was opened")
	}
}

func TestNonSensitiveFIFOIsNotOpened(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "specs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(root, "specs", "pipe"), 0o644); err != nil {
		t.Fatal(err)
	}
	fs, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		ent, err := fs.StatRel("specs/pipe")
		if err != nil {
			t.Errorf("stat fifo: %v", err)
			return
		}
		if !ent.Exists || ent.Regular {
			t.Errorf("fifo metadata %+v", ent)
		}
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("metadata check blocked on a FIFO")
	}
}

func TestSourceChangeAndLimit(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "notes.txt")
	if err := os.WriteFile(path, []byte("stable"), 0o644); err != nil {
		t.Fatal(err)
	}
	fs, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	fs.SetAfterReadHook(func(p string) {
		if err := os.WriteFile(p, []byte("changed-contents"), 0o644); err != nil {
			t.Error(err)
		}
	})
	if _, err := fs.ReadRel("notes.txt"); err == nil || err.Error() != CodeChanged {
		t.Fatalf("change err=%v", err)
	}
	fs.SetAfterReadHook(nil)
	big := filepath.Join(root, "big.txt")
	f, err := os.Create(big)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Truncate(MaxFileBytes + 1); err != nil {
		t.Fatal(err)
	}
	f.Close()
	if _, err := fs.ReadRel("big.txt"); err == nil || err.Error() != CodeLimit {
		t.Fatalf("limit err=%v", err)
	}
}

// A FIFO swapped in after the Lstat check must not block the open (review D4).
func TestOpenNoFollowDoesNotBlockOnFIFO(t *testing.T) {
	root := t.TempDir()
	pipe := filepath.Join(root, "tasks.md")
	if err := syscall.Mkfifo(pipe, 0o644); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		f, err := openNoFollow(pipe)
		if err == nil {
			info, statErr := f.Stat()
			f.Close()
			if statErr == nil && info.Mode().IsRegular() {
				err = os.ErrInvalid
			}
		}
		done <- err
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("openNoFollow blocked on a FIFO")
	}
}

// --- ListDir (feature 006, research R5) ---

func mkListRoot(t *testing.T) (string, *FS) {
	t.Helper()
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "specs"), 0o755); err != nil {
		t.Fatal(err)
	}
	fs, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	return root, fs
}

func TestListDirKindsAndByteOrder(t *testing.T) {
	root, fs := mkListRoot(t)
	specs := filepath.Join(root, "specs")
	for _, d := range []string{"b-dir", "A-dir", "001-alpha"} {
		if err := os.Mkdir(filepath.Join(specs, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(specs, "README.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(specs, "link")); err != nil {
		t.Fatal(err)
	}
	if err := syscall.Mkfifo(filepath.Join(specs, "pipe"), 0o644); err != nil {
		t.Skip("mkfifo unavailable")
	}
	// A sensitive-named FIFO: opening it would block, so it must never be opened or stat'ed for content.
	if err := syscall.Mkfifo(filepath.Join(specs, "my-secret"), 0o644); err != nil {
		t.Fatal(err)
	}
	beforeFiles, beforeBytes := fs.files, fs.bytes
	got, err := fs.ListDir("specs")
	if err != nil {
		t.Fatal(err)
	}
	want := []struct{ name, kind string }{
		{"001-alpha", KindDir}, {"A-dir", KindDir}, {"README.md", KindFile}, {"b-dir", KindDir},
		{"link", KindSymlink}, {"my-secret", KindSensitive}, {"pipe", KindOther},
	}
	if len(got) != len(want) {
		t.Fatalf("entries %+v", got)
	}
	for i, w := range want {
		if got[i].Name != w.name || got[i].Kind != w.kind {
			t.Errorf("entry %d = %s/%s, want %s/%s", i, got[i].Name, got[i].Kind, w.name, w.kind)
		}
	}
	if fs.files != beforeFiles || fs.bytes != beforeBytes {
		t.Error("listing must not consume read budgets")
	}
}

func TestListDirRejectsUnsafeTargets(t *testing.T) {
	root, fs := mkListRoot(t)
	if _, err := fs.ListDir("missing"); err == nil || err.Error() != CodeNotFound {
		t.Fatalf("missing err=%v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "file"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := fs.ListDir("file"); err == nil || err.Error() != CodeUnsafe {
		t.Fatalf("file err=%v", err)
	}
	if err := os.Symlink(filepath.Join(root, "specs"), filepath.Join(root, "alias")); err != nil {
		t.Fatal(err)
	}
	if _, err := fs.ListDir("alias"); err == nil || err.Error() != CodeUnsafe {
		t.Fatalf("symlinked dir err=%v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "specs", "real"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := fs.ListDir("alias/real"); err == nil || err.Error() != CodeUnsafe {
		t.Fatalf("symlink component err=%v", err)
	}
	if _, err := fs.ListDir("../x"); err == nil || err.Error() != CodeUnsafe {
		t.Fatalf("dot-dot err=%v", err)
	}
	if _, err := fs.ListDir("specs/.git"); err == nil || err.Error() != CodeSensitive {
		t.Fatalf("sensitive dir err=%v", err)
	}
}

func TestListDirLimit(t *testing.T) {
	root, fs := mkListRoot(t)
	for i := 0; i <= MaxDirEntries; i++ {
		if err := os.Mkdir(filepath.Join(root, "specs", itoa(1000+i)), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := fs.ListDir("specs"); err == nil || err.Error() != CodeLimit {
		t.Fatalf("limit err=%v", err)
	}
}

func TestListDirDetectsChange(t *testing.T) {
	root, fs := mkListRoot(t)
	if err := os.Mkdir(filepath.Join(root, "specs", "001-a"), 0o755); err != nil {
		t.Fatal(err)
	}
	fs.SetAfterReadHook(func(string) {
		time.Sleep(10 * time.Millisecond)
		_ = os.Mkdir(filepath.Join(root, "specs", "002-b"), 0o755)
	})
	if _, err := fs.ListDir("specs"); err == nil || err.Error() != CodeChanged {
		t.Fatalf("change err=%v", err)
	}
}
