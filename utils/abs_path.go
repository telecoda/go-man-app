package utils

import (
	"path"
	"runtime"
)

func GetAbsolutePathOfCurrentPackage(appendPath string) string {
	_, filename, _, _ := runtime.Caller(1)

	absPath := path.Join(path.Dir(filename), appendPath)

	return absPath
}
