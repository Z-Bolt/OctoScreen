package utils

import (
	"os"
)

// FileExists checks if a file exists and is not a directory.
// From https://golangcode.com/check-if-a-file-exists/
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}
