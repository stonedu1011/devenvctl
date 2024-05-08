package utils

import "path/filepath"

func AbsPath(path string, base string) string {
	path = filepath.Clean(path)
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(base, path)
}
