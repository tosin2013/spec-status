//go:build linux

package cli

import "golang.org/x/sys/unix"

// ioctlGetTermios is the request that reads terminal attributes on Linux.
const ioctlGetTermios = unix.TCGETS
