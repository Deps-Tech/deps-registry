package filesystem

import (
	"os"
	"path/filepath"
)

func ListFiles(dir string, extensions ...string) ([]string, error) {
	var files []string
	extMap := make(map[string]bool)
	for _, ext := range extensions {
		extMap[ext] = true
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if len(extMap) == 0 || extMap[filepath.Ext(path)] {
			relPath, _ := filepath.Rel(dir, path)
			files = append(files, relPath)
		}

		return nil
	})

	return files, err
}
