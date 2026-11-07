//go:build !linux

package cli

import "errors"

func syscallMkfifo(string) error {
	return errors.New("mkfifo is only exercised on linux")
}
