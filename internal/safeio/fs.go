// Package safeio reads regular files inside one selected project root.
// It is the only filesystem boundary for Spec Status.
//
// Limits and the sensitive-name matcher are PROFILE-R1 implementation
// proposals, not owner-approved policy. The matcher is intentionally
// over-inclusive and is not a security certification.
package safeio

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	MaxFileBytes  = 2 << 20
	MaxTotalBytes = 32 << 20
	MaxFiles      = 1024
	MaxPathBytes  = 4096

	CodeUnsafe      = "unsafe"
	CodeSensitive   = "sensitive"
	CodeLimit       = "limit"
	CodeChanged     = "changed"
	CodeNotFound    = "notfound"
	CodeUnavailable = "unavailable"
	CodeUnsupported = "unsupported"
)

// Error is a stable confinement failure. Error text is the code only;
// callers must not surface raw operating-system errors.
type Error struct {
	Code string
}

func (e *Error) Error() string { return e.Code }

func errCode(code string) error { return &Error{Code: code} }

// Entry is metadata from a confined lstat. Contents are not read.
type Entry struct {
	Exists  bool
	Dir     bool
	Regular bool
	// ID identifies the file by device and inode where the platform reports
	// them, so two path spellings of one file (for example on a
	// case-insensitive filesystem) can be recognized. Empty if unknown.
	ID string
}

// FS is a root-confined reader. It does not write, follow linked-path
// symlinks, or cross a detectable device boundary.
type FS struct {
	root      string
	dev       uint64
	files     int
	bytes     int64
	afterRead func(path string)
}

// SetAfterReadHook is a test seam. Production leaves it unset.
// The hook runs after a successful open and before the second stat,
// so tests can observe source-change detection.
func (f *FS) SetAfterReadHook(fn func(path string)) { f.afterRead = fn }

// New canonicalizes root once. Later reads must stay under that directory.
func New(root string) (*FS, error) {
	if root == "" {
		return nil, errCode(CodeUnavailable)
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, errCode(CodeUnavailable)
	}
	canon, err := filepath.EvalSymlinks(abs)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errCode(CodeNotFound)
		}
		return nil, errCode(CodeUnavailable)
	}
	info, err := os.Stat(canon)
	if err != nil || !info.IsDir() {
		return nil, errCode(CodeUnavailable)
	}
	dev, ok := deviceOf(info)
	if !ok {
		return nil, errCode(CodeUnsupported)
	}
	return &FS{root: canon, dev: dev}, nil
}

// Root is the canonical absolute root. It is for confinement checks, not for reports.
func (f *FS) Root() string { return f.root }

// Relativize converts an absolute path to a root-relative slash path.
// Absolute paths outside the root are unsafe.
func (f *FS) Relativize(path string) (string, error) {
	if path == "" {
		return "", errCode(CodeUnsafe)
	}
	if filepath.IsAbs(path) {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", errCode(CodeUnsafe)
		}
		rel, err := filepath.Rel(f.root, abs)
		if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return "", errCode(CodeUnsafe)
		}
		return NormalizeRel(filepath.ToSlash(rel))
	}
	return NormalizeRel(filepath.ToSlash(path))
}

// NormalizeRel rejects traversal, including ".." that would re-enter the root.
func NormalizeRel(rel string) (string, error) {
	if rel == "" || rel == "." {
		return "", errCode(CodeUnsafe)
	}
	rel = strings.TrimPrefix(rel, "./")
	if strings.Contains(rel, "\\") || strings.ContainsRune(rel, 0) || strings.Contains(rel, "%") {
		return "", errCode(CodeUnsafe)
	}
	if filepath.IsAbs(rel) || strings.HasPrefix(rel, "/") {
		return "", errCode(CodeUnsafe)
	}
	parts := strings.Split(rel, "/")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p == "" || p == "." {
			continue
		}
		if p == ".." {
			return "", errCode(CodeUnsafe)
		}
		out = append(out, p)
	}
	if len(out) == 0 {
		return "", errCode(CodeUnsafe)
	}
	clean := strings.Join(out, "/")
	if len(clean) > MaxPathBytes {
		return "", errCode(CodeLimit)
	}
	return clean, nil
}

// Sensitive reports whether a normalized relative path matches the
// candidate deny matcher. The file is not opened.
func Sensitive(rel string) bool {
	parts := strings.Split(rel, "/")
	for _, part := range parts {
		if deniedComponent(part) || deniedBase(part) {
			return true
		}
	}
	return false
}

func deniedComponent(part string) bool {
	switch strings.ToLower(part) {
	case ".git", ".ssh", ".aws", ".gnupg":
		return true
	default:
		return false
	}
}

func deniedBase(part string) bool {
	base := strings.ToLower(part)
	if base == ".env" || strings.HasPrefix(base, ".env.") {
		return true
	}
	for _, s := range []string{".pem", ".key", ".p12", ".pfx"} {
		if strings.HasSuffix(base, s) {
			return true
		}
	}
	for _, token := range []string{"secret", "credentials", "private-key", "answer-key"} {
		if strings.Contains(base, token) {
			return true
		}
	}
	return false
}

// StatRel is a confined metadata check. It does not read file contents.
func (f *FS) StatRel(rel string) (Entry, error) {
	full, err := f.prepare(rel)
	if err != nil {
		return Entry{}, err
	}
	info, err := os.Lstat(full)
	if err != nil {
		if os.IsNotExist(err) {
			return Entry{}, nil
		}
		return Entry{}, errCode(CodeUnavailable)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return Entry{}, errCode(CodeUnsafe)
	}
	if dev, ok := deviceOf(info); ok && dev != f.dev {
		return Entry{}, errCode(CodeUnsafe)
	}
	return Entry{Exists: true, Dir: info.IsDir(), Regular: info.Mode().IsRegular(), ID: fileID(info)}, nil
}

// ReadRel reads one regular file under the root, counting it against the budgets.
func (f *FS) ReadRel(rel string) ([]byte, error) {
	full, err := f.prepare(rel)
	if err != nil {
		return nil, err
	}
	before, err := os.Lstat(full)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errCode(CodeNotFound)
		}
		return nil, errCode(CodeUnavailable)
	}
	if before.Mode()&os.ModeSymlink != 0 || !before.Mode().IsRegular() {
		return nil, errCode(CodeUnsafe)
	}
	if before.Size() > MaxFileBytes {
		return nil, errCode(CodeLimit)
	}
	if f.files+1 > MaxFiles || f.bytes+before.Size() > MaxTotalBytes {
		return nil, errCode(CodeLimit)
	}
	file, err := openNoFollow(full)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil {
		return nil, errCode(CodeUnavailable)
	}
	if !opened.Mode().IsRegular() {
		return nil, errCode(CodeUnsafe)
	}
	if dev, ok := deviceOf(opened); !ok || dev != f.dev {
		return nil, errCode(CodeUnsafe)
	}
	if err := confirmOpened(file, f.root); err != nil {
		return nil, err
	}
	if opened.Size() > MaxFileBytes {
		return nil, errCode(CodeLimit)
	}
	buf := make([]byte, opened.Size())
	n, err := readFull(file, buf)
	if err != nil {
		return nil, errCode(CodeUnavailable)
	}
	if int64(n) != opened.Size() || int64(n) > MaxFileBytes || f.bytes+int64(n) > MaxTotalBytes {
		if int64(n) != opened.Size() {
			return nil, errCode(CodeChanged)
		}
		return nil, errCode(CodeLimit)
	}
	if f.afterRead != nil {
		f.afterRead(full)
	}
	after, err := os.Lstat(full)
	if err != nil || !sameIdentity(before, after) || !sameIdentity(before, opened) {
		return nil, errCode(CodeChanged)
	}
	f.files++
	f.bytes += int64(n)
	return buf[:n], nil
}

func (f *FS) prepare(rel string) (string, error) {
	clean, err := NormalizeRel(rel)
	if err != nil {
		return "", err
	}
	if Sensitive(clean) {
		return "", errCode(CodeSensitive)
	}
	full := filepath.Join(f.root, filepath.FromSlash(clean))
	if !within(f.root, full) {
		return "", errCode(CodeUnsafe)
	}
	if err := rejectSymlinkParents(f.root, clean); err != nil {
		return "", err
	}
	return full, nil
}

func rejectSymlinkParents(root, rel string) error {
	cur := root
	for _, part := range strings.Split(rel, "/") {
		cur = filepath.Join(cur, part)
		info, err := os.Lstat(cur)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return errCode(CodeUnavailable)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return errCode(CodeUnsafe)
		}
	}
	return nil
}

func within(root, full string) bool {
	rel, err := filepath.Rel(root, full)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func sameIdentity(a, b os.FileInfo) bool {
	if a.Size() != b.Size() || !a.ModTime().Equal(b.ModTime()) || a.Mode() != b.Mode() {
		return false
	}
	da, aok := deviceOf(a)
	db, bok := deviceOf(b)
	return aok && bok && da == db
}

func readFull(f *os.File, buf []byte) (int, error) {
	n, err := io.ReadFull(f, buf)
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		err = nil
	}
	if err != nil {
		return n, err
	}
	var extra [1]byte
	got, rerr := f.Read(extra[:])
	if got > 0 {
		return n + got, nil
	}
	if rerr != nil && rerr != io.EOF {
		return n, rerr
	}
	return n, nil
}
