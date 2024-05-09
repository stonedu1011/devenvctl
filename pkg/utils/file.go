package utils

import (
	"fmt"
	cp "github.com/otiai10/copy"
	"path/filepath"
)

func AbsPath[T any](path string, base T) string {
	var basePath string
	switch v := any(base).(type) {
	case string:
		basePath = v
	case fmt.Stringer:
		basePath = v.String()
	default:
		basePath = fmt.Sprintf("%v", base)
	}
	path = filepath.Clean(path)
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(basePath, path)
}

func CopyDir(src, dst string) error {
	return cp.Copy(src, dst)
}