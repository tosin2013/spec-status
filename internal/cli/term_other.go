//go:build !linux && !darwin

package cli

import (
	"io"
	"os"
	"strconv"
)

func writerIsTerminal(io.Writer) bool { return false }

func textColumns() int {
	if raw := os.Getenv("COLUMNS"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err == nil {
			return n
		}
	}
	return 80
}
