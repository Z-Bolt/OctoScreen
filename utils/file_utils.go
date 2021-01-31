package utils

import (
	"os"
	// "strings"
)

// FileExists checks if a file exists and is not a directory.
// From https://golangcode.com/check-if-a-file-exists/
func FileExists(filename string) bool {
	info := exists(filename)
	if info == nil {
		return false
	}

	return !info.IsDir()
}

func FolderExists(filename string) bool {
	info := exists(filename)
	if info == nil {
		return false
	}

	return info.IsDir()
}

func exists(filename string) os.FileInfo {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil
	}

	return info
}
