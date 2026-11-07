//go:build darwin

package cli

import "golang.org/x/sys/unix"

// ioctlGetTermios is the request that reads terminal attributes on Darwin.
const ioctlGetTermios = unix.TIOCGETA
