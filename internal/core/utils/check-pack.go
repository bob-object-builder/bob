package utils

import (
	"os"
)

type PathType int

const (
	PathNotExist PathType = iota
	PathIsFile
	PathIsDir
)

func CheckPath(path string) PathType {
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
