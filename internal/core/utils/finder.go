package utils

import (
	"os"
	"path/filepath"
)

func FindBobFiles(directory string) ([]string, error) {
	var files []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".bob" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
