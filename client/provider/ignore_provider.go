package provider

import (
	"bytes"
	"github.com/gojek/stevedore/pkg/config"
	"github.com/gojek/stevedore/pkg/file"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/spf13/afero"
	"path"
	"path/filepath"
)

const stevedoreIgnoreFileName = ".stevedoreignore"

// IgnoreProvider is the Ignore ProviderImpl interface
type IgnoreProvider interface {
	Files() ([]string, error)
	Ignores() (stevedore.Ignores, error)
}

// DefaultIgnoreProvider represents the default ignore provider
// which reads the context from file
type DefaultIgnoreProvider struct {
	fs           afero.Fs
	manifestPath string
	cwd          string
}

// Files returns all the stevedoreignore files
func (provider DefaultIgnoreProvider) Files() ([]string, error) {
	files := make([]string, 0)
	manifestPath, err := dir(provider.fs, provider.manifestPath)
	if err != nil {
		return nil, err
	}

	knownPaths := []string{manifestPath, provider.cwd}
	for _, knownPath := range knownPaths {
		if ignoreFile, ok, err := ignoreFile(provider.fs, knownPath); ok {
			files = append(files, ignoreFile)
		} else if err != nil {
			return files, err
		}
	}
	return files, nil
}

// Ignores returns the ignores if any by reading .stevedoreignore file
func (provider DefaultIgnoreProvider) Ignores() (stevedore.Ignores, error) {
	ignoreFiles, err := provider.Files()
	fileErrors := file.Errors{}
	result := make(stevedore.Ignores, 0, len(ignoreFiles))
	if err != nil {
		return nil, err
	}

	for _, ignoreFile := range ignoreFiles {
		ignores, err := readIgnoresFrom(provider.fs, ignoreFile)
		if err != nil {
			fileErrors = append(fileErrors, file.Error{
				Reason:   err,
				Filename: ignoreFile,
			})
		}
		result = append(result, ignores...)
	}
	if len(fileErrors) != 0 {
		return result, fileErrors
	}
	return result, nil
}

func ignoreFile(fs afero.Fs, dir string) (string, bool, error) {
	stevedoreIgnoreFilePath := path.Join(dir, stevedoreIgnoreFileName)
	exists, err := afero.Exists(fs, stevedoreIgnoreFilePath)
	if err != nil {
		return "", false, err
	}
	return stevedoreIgnoreFilePath, exists, nil
}

func readIgnoresFrom(fs afero.Fs, stevedoreIgnoreFilePath string) (stevedore.Ignores, error) {
	if exists, err := afero.Exists(fs, stevedoreIgnoreFilePath); err != nil {
		return nil, err
	} else if !exists {
		return stevedore.Ignores{}, nil
	}

	data, err := afero.ReadFile(fs, stevedoreIgnoreFilePath)
	if err != nil {
		return nil, err
	}

	return stevedore.NewIgnores(bytes.NewReader(data))
}

func dir(fs afero.Fs, path string) (string, error) {
	isDir, err := afero.DirExists(fs, path)
	if err != nil {
		return "", err
	}

	if isDir {
		return path, nil
	}

	return filepath.Dir(path), nil
}

// NewIgnoreProvider returns an instance of override.IgnoreProvider
func NewIgnoreProvider(fs afero.Fs, manifestPath string, environment config.Environment) (IgnoreProvider, error) {
	dir, err := environment.Cwd()
	return DefaultIgnoreProvider{fs: fs, manifestPath: manifestPath, cwd: dir}, err
}
