package util

import (
	"os"
	"path/filepath"
	"strings"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func RemoveFileExtension(basename string) string {
	return strings.TrimSuffix(basename, filepath.Ext(basename))
}
