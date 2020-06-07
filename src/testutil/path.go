package testutil

import (
	"path/filepath"
	"runtime"
)

// SourceFilePath returns absolute path of current source file
func SourceFilePath() string {
	_, fileName, _, _ := runtime.Caller(1)
	return fileName
}

// ProjectPath returns absolute path of source root
func ProjectPath() string {
	_, fileName, _, _ := runtime.Caller(0)
	return filepath.Join(fileName, "../../..")
}
