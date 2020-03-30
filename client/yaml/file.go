package yaml

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/gojek/stevedore/pkg/file"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	"io"
	"path"
)

// Files represents collection of file
type Files []File

const (
	kindMeta    = "kind"
	versionMeta = "version"
)

// Names return the file name(s)
func (files Files) Names() []string {
	result := make([]string, 0, len(files))
	for _, yamlFile := range files {
		result = append(result, yamlFile.Name)
	}
	return result
}

// File represents file along with the file content
type File struct {
	Name string
	data []byte
	meta map[string]string
}

// Kind returns the yaml kind and true if meta contains
// kind information
func (file File) Kind() (string, bool) {
	kind, ok := file.meta[kindMeta]
	return kind, ok
}

// Version returns the yaml version and true if meta contains
// version information
func (file File) Version() (string, bool) {
	kind, ok := file.meta[versionMeta]
	return kind, ok
}

// Check checks if the yaml file has the specified kind and version
// returns true if it matches the kind and specified version
// else returns false
func (file File) Check(kind, version string) (bool, error) {
	if fileKind, ok := file.Kind(); !ok || fileKind != kind {
		return false, nil
	}

	fileVersion, ok := file.Version()
	if !ok {
		return false, nil
	}

	semVersion, err := semver.NewVersion(fileVersion)
	if err != nil {
		return false, fmt.Errorf("unable to detect version information from file %s", file.Name)
	}

	constraints, err := semver.NewConstraint(fmt.Sprintf("= %s", version))
	if err != nil {
		return false, fmt.Errorf("unable to check version information from file %s", file.Name)
	}

	return constraints.Check(semVersion), nil
}

// NewYamlFile returns new File
func NewYamlFile(fs afero.Fs, filename string) (File, error) {
	data, err := afero.ReadFile(fs, filename)
	if err != nil {
		return File{}, file.Error{Filename: filename, Reason: err}
	}

	meta := map[string]string{}
	_ = yaml.Unmarshal(data, &meta)
	return File{Name: filename, data: data, meta: meta}, nil
}

// Reader return data as io.Reader
func (file File) Reader() io.Reader {
	if file.data != nil {
		return bytes.NewReader(file.data)
	}
	return nil
}

// EmptyFolderError represents empty folder error
type EmptyFolderError string

func (err EmptyFolderError) Error() string {
	return string(err)
}

func fileNames(fs afero.Fs, filePath string) ([]string, error) {
	isDir, err := afero.IsDir(fs, filePath)
	if err != nil {
		formattedErr := fmt.Errorf("unable to get file(s) from \"%s\", failed with error: %v", filePath, err)
		return nil, formattedErr
	}

	if isDir {
		files, err := afero.Glob(fs, path.Join(filePath, "/*.yaml"))
		if err != nil {
			formattedErr := fmt.Errorf("unable to get files from folder %s, failed with formattedErr: %v", filePath, err)
			return nil, formattedErr
		}

		if len(files) == 0 {
			err := EmptyFolderError(fmt.Sprintf("folder %s does not have any yamls to apply", filePath))
			return nil, err
		}
		return files, nil
	}
	return []string{filePath}, nil
}

// NewYamlFiles returns new Files
func NewYamlFiles(fs afero.Fs, filePath string) (Files, error) {
	fileNames, err := fileNames(fs, filePath)
	if err != nil {
		return nil, err
	}

	errors := file.Errors{}
	yamlFiles := Files{}

	for _, filename := range fileNames {
		yamlFile, err := NewYamlFile(fs, filename)
		if err != nil {
			errors = append(errors, file.Error{Filename: filename, Reason: err})
			continue
		}
		yamlFiles = append(yamlFiles, yamlFile)
	}

	if len(errors) != 0 {
		return nil, errors
	}

	return yamlFiles, nil
}
