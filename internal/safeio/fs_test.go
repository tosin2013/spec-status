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
