//go:build !linux && !darwin

package cli

import "errors"

func syscallMkfifo(string) error {
	return errors.New("mkfifo is only exercised on linux and darwin")
}
