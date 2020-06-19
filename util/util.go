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

func Unique(slice []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueSlice []string
	for _, entry := range slice {
		if _, value := uniqueMap[entry]; !value {
			uniqueMap[entry] = true
			uniqueSlice = append(uniqueSlice, entry)
		}
	}
	return uniqueSlice
}

func Filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
