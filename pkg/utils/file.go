package utils

import (
	"github.com/spf13/afero"
	"os"
)

// FileUtils is wrapper around fs to give additional functions
type FileUtils struct {
	afero.Fs
}

// NewOsFileUtils constructs new file utils with OS file system
func NewOsFileUtils() FileUtils {
	return FileUtils{Fs: afero.NewOsFs()}
}

// TempDir creates a new temporary directory in the directory dir
// with a name beginning with prefix and returns the path of the
// new directory.
func (f FileUtils) TempDir(prefix string) (name string, err error) {
	return afero.TempDir(f.Fs, os.TempDir(), prefix)
}

// WriteFile writes data to a file named by filename.
func (f FileUtils) WriteFile(filename string, data []byte, perm os.FileMode) error {
	return afero.WriteFile(f.Fs, filename, data, perm)
}
