package manifest

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
)

var readWritePermission = os.FileMode(0666)

// Persistence is an interface for writing data
type Persistence interface {
	write(filepath string, data []byte) error
}

// DiskPersistence is an interface for writing data into disk
type DiskPersistence struct {
	fs afero.Fs
}

func (persistence DiskPersistence) write(filepath string, data []byte) error {
	return afero.WriteFile(persistence.fs, filepath, data, readWritePermission)
}

// ConsolePersistence is an interface for writing data into console
type ConsolePersistence struct{}

func (persistence ConsolePersistence) write(_ string, data []byte) error {
	_, err := fmt.Fprintf(os.Stdout, "---\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(os.Stdout, "%s\n", string(data))
	if err != nil {
		return err
	}
	return nil
}
