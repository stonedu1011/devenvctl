package utils

import (
	"embed"
	"fmt"
	cp "github.com/otiai10/copy"
	"io/fs"
	"os"
	"path/filepath"
)

func AbsPath[T any](path string, base T) string {
	var basePath string
	switch v := any(base).(type) {
	case string:
		basePath = v
	case embed.FS:
		basePath = "presets"
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

func CopyDir(srcFS fs.FS, src, dst string) error {
	switch fsys := srcFS.(type) {
	case embed.FS:
		return copyEmbedDir(fsys, src, dst)
	case fs.StatFS:
		return cp.Copy(src, dst)
	default:
		return fmt.Errorf(`cannot copy directory - supported source file system`)
	}
}

func copyEmbedDir(srcFS embed.FS, src, dst string) error {
	return fs.WalkDir(srcFS, src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, e := filepath.Rel(src, path)
		if e != nil {
			return e
		}
		dstPath := filepath.Join(dst, relPath)
		if d.IsDir() {
			if e := os.MkdirAll(dstPath, 0755); e != nil {
				return e
			}
		} else {
			data, e := fs.ReadFile(srcFS, path)
			if e != nil {
				return e
			}
			if e := os.WriteFile(dstPath, data, d.Type()); e != nil {
				return e
			}
		}
		return nil
	})
}
