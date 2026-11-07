//go:build linux

package cli

import "syscall"

func syscallMkfifo(path string) error {
	return syscall.Mkfifo(path, 0o644)
}
