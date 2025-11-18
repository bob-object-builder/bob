package checker

import (
	"os"
	"path/filepath"
	"strings"
)

type PathType int

const (
	PathNotExist PathType = iota
	PathIsFile
	PathIsDir
)

func ValidatePath(path string) PathType {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return PathNotExist
	} else if err != nil {
		return PathNotExist
	}

	if info.IsDir() {
		return PathIsDir
	}

	return PathIsFile
}

func CheckPath(path string) PathType {
	if path == "" {
		return PathNotExist
	}

	if path == "." || path == ".." {
		return PathIsDir
	}

	cleaned := filepath.Clean(path)

	if strings.HasSuffix(path, "/") {
		return PathIsDir
	}

	if filepath.Ext(cleaned) != "" {
		return PathIsFile
	}

	return PathIsDir
}
