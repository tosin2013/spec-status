//go:build !linux && !darwin

package safeio

import "os"

func deviceOf(os.FileInfo) (uint64, bool) { return 0, false }

func fileID(os.FileInfo) string { return "" }

func openNoFollow(string) (*os.File, error) {
	return nil, errCode(CodeUnsupported)
}

func confirmOpened(*os.File, string) error {
	return errCode(CodeUnsupported)
}
