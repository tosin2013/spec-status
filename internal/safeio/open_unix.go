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

func openNoFollow(path string) (*os.File, error) {
	fd, err := syscall.Open(path, syscall.O_RDONLY|syscall.O_NOFOLLOW, 0)
	if err != nil {
		return nil, mapOpenErr(err)
	}
	return os.NewFile(uintptr(fd), path), nil
}

func confirmOpened(f *os.File, root string) error {
	if runtime.GOOS != "linux" {
		// Darwin builds keep O_NOFOLLOW and the device check.
		// This environment did not execute a Darwin binary.
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
