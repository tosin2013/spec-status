//go:build linux || darwin

package safeio

import (
	"os"
	"runtime"
	"strings"
	"syscall"
)

func deviceOf(fi os.FileInfo) (uint64, bool) {
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok || st == nil {
		return 0, false
	}
	return uint64(st.Dev), true
}

// openNoFollow opens without following a final symlink. O_NONBLOCK keeps the
// open from waiting on a FIFO that replaced a regular file after the Lstat
// check; the caller's fstat then rejects anything that is not a regular file.
// Regular-file reads are unaffected by O_NONBLOCK.
func openNoFollow(path string) (*os.File, error) {
	fd, err := syscall.Open(path, syscall.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK|syscall.O_CLOEXEC, 0)
	if err != nil {
		return nil, mapOpenErr(err)
	}
	return os.NewFile(uintptr(fd), path), nil
}

func confirmOpened(f *os.File, root string) error {
	if runtime.GOOS != "linux" {
		// Darwin has no /proc/self/fd. Reads there rely on O_NOFOLLOW, the
		// symlink-component check, the device check, and the before/after
		// identity check. A parent directory swapped between those checks
		// and the open is a documented residual limit on Darwin.
		return nil
	}
	link, err := os.Readlink("/proc/self/fd/" + itoa(int(f.Fd())))
	if err != nil {
		return errCode(CodeUnsupported)
	}
	if strings.HasSuffix(link, " (deleted)") {
		return errCode(CodeChanged)
	}
	if !within(root, link) {
		return errCode(CodeUnsafe)
	}
	return nil
}

func mapOpenErr(err error) error {
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		return errCode(CodeNotFound)
	}
	if err == syscall.ELOOP {
		return errCode(CodeUnsafe)
	}
	return errCode(CodeUnavailable)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [32]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
