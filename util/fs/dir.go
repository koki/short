package fs

import (
	"os"
	"path/filepath"
)

func MkdirP(path string) error {
	dir, _ := filepath.Split(path)
	return os.MkdirAll(dir, 0777)
}
