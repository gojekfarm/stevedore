package stevedore

import (
	"github.com/spf13/afero"
	"os"
)

// FileUtils to combine Fs and other package methods
type FileUtils interface {
	afero.Fs
	TempDir(prefix string) (name string, err error)
	WriteFile(filename string, data []byte, perm os.FileMode) error
}
