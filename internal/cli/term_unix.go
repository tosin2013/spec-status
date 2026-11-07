//go:build linux || darwin

package cli

import (
	"io"
	"os"
	"strconv"

	"golang.org/x/sys/unix"
)

// writerIsTerminal reports whether w is the terminal being written to.
// It uses an ioctl on that file descriptor only. It does not open /dev/tty
// and it does not send a terminal query.
func writerIsTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok || f == nil {
		return false
	}
	_, err := unix.IoctlGetTermios(int(f.Fd()), unix.TCGETS)
	return err == nil
}

func textColumns() int {
	if raw := os.Getenv("COLUMNS"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err == nil {
			return n
		}
	}
	return 80
}
