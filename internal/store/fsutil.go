package store

import (
	"os"
	"path/filepath"
)

func mkdirAllFor(path string) error {
	return os.MkdirAll(filepath.Dir(path), 0o755)
}

func writeFileBytes(path string, data []byte) error {
	return os.WriteFile(path, data, 0o644)
}
