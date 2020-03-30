package yaml_test

import (
	"fmt"
	"github.com/gojek/stevedore/client/internal/mocks"
	"github.com/gojek/stevedore/client/yaml"
	"github.com/gojek/stevedore/pkg/stevedore"
	"github.com/golang/mock/gomock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestNewYamlFile(t *testing.T) {
	t.Run("should create new file", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file"
		content := `
---
kind: StevedoreManifest
version: "2"
`
		_ = afero.WriteFile(memFs, filename, []byte(content), 0644)

		file, err := yaml.NewYamlFile(memFs, filename)
		assert.NotNil(t, file)
		assert.NoError(t, err)

		actualKind, ok := file.Kind()
		assert.True(t, ok)
		assert.Equal(t, "StevedoreManifest", actualKind)

		actualVersion, ok := file.Version()
		assert.True(t, ok)
		assert.Equal(t, "2", actualVersion)

		data, err := ioutil.ReadAll(file.Reader())
		assert.NoError(t, err)
		assert.Equal(t, content, string(data))
	})

	t.Run("should return false if kind and version info are available", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file"
		content := `
---
`
		_ = afero.WriteFile(memFs, filename, []byte(content), 0644)

		file, err := yaml.NewYamlFile(memFs, filename)
		assert.NotNil(t, file)
		assert.NoError(t, err)

		actualKind, ok := file.Kind()
		assert.False(t, ok)
		assert.Equal(t, "", actualKind)

		actualVersion, ok := file.Version()
		assert.False(t, ok)
		assert.Equal(t, "", actualVersion)

		data, err := ioutil.ReadAll(file.Reader())
		assert.NoError(t, err)
		assert.Equal(t, content, string(data))
	})

	t.Run("should return error", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file"

		file, err := yaml.NewYamlFile(memFs, filename)

		assert.Equal(t, yaml.File{}, file)
		if assert.Error(t, err) {
			assert.Equal(t, "Error while processing /mock/file\nopen /mock/file: file does not exist", err.Error())
		}
	})
}

func TestFilesNames(t *testing.T) {
	t.Run("it should return the file names", func(t *testing.T) {
		files := yaml.Files{
			{Name: "one.yaml"},
			{Name: "two.json"},
			{Name: "three.txt"},
		}

		actual := files.Names()

		assert.Equal(t, []string{"one.yaml", "two.json", "three.txt"}, actual)
	})
}

func TestNewYamlFiles(t *testing.T) {
	t.Run("should read all the yaml files", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file.yaml"
		anotherFile := "/mock/anotherFile.yaml"

		_ = afero.WriteFile(memFs, filename, []byte("Name: x-stevedore"), 0644)
		_ = afero.WriteFile(memFs, anotherFile, []byte("Name: y-stevedore"), 0644)

		files, err := yaml.NewYamlFiles(memFs, "/mock")

		assert.NoError(t, err)
		assert.NotNil(t, files)
		assert.Equal(t, 2, len(files))
	})

	t.Run("should the file if its not dir", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file.yaml"

		_ = afero.WriteFile(memFs, filename, []byte("Name: x-stevedore"), 0644)

		files, err := yaml.NewYamlFiles(memFs, "/mock/file.yaml")

		assert.NoError(t, err)
		assert.NotNil(t, files)
		assert.Equal(t, 1, len(files))
	})

	t.Run("should handle if the folder doesn't exists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFs := mocks.NewMockFs(ctrl)

		mockFs.EXPECT().Stat("/mock").Return(nil, fmt.Errorf("error while opening the dir"))

		files, err := yaml.NewYamlFiles(mockFs, "/mock")

		if assert.Error(t, err) {
			assert.Equal(t, "unable to get file(s) from \"/mock\", failed with error: error while opening the dir", err.Error())
		}
		assert.Equal(t, yaml.Files(nil), files)
	})

	t.Run("should handle if a file in a folder has error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFs := mocks.NewMockFs(ctrl)
		mockFileInfo := mocks.NewMockInfo(ctrl)
		mockFile := mocks.NewMockFile(ctrl)

		mockFs.EXPECT().Stat("/mock").Return(mockFileInfo, nil).Times(2)
		mockFileInfo.EXPECT().IsDir().Return(true).Times(2)
		mockFs.EXPECT().Open("/mock").Return(mockFile, nil)
		mockFs.EXPECT().Open("/mock/one.yaml").Return(nil, fmt.Errorf("unable to open the file"))

		mockFile.EXPECT().Readdirnames(-1).Return([]string{"one.yaml"}, nil)
		mockFile.EXPECT().Close()

		files, err := yaml.NewYamlFiles(mockFs, "/mock")

		if assert.Error(t, err) {
			assert.Equal(t, "Failed with 1 Errors\n\n1. Error while processing /mock/one.yaml\nError while processing /mock/one.yaml\nunable to open the file\n", err.Error())
		}
		assert.Equal(t, yaml.Files(nil), files)
	})

	t.Run("should return error if no yaml files exists in the given dir", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock"

		_ = memFs.MkdirAll(filename, 0644)

		file, err := yaml.NewYamlFiles(memFs, filename)

		assert.Equal(t, yaml.Files(nil), file)
		if assert.Error(t, err) {
			assert.Equal(t, "folder /mock does not have any yamls to apply", err.Error())
		}
	})
}

func TestFileKind(t *testing.T) {
	t.Run("should return the kind if specified", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file"
		content := `
---
kind: SomeKind
`
		_ = afero.WriteFile(memFs, filename, []byte(content), 0644)

		file, err := yaml.NewYamlFile(memFs, filename)

		assert.NoError(t, err)
		actualKind, ok := file.Kind()
		assert.True(t, ok)
		assert.Equal(t, "SomeKind", actualKind)
	})

	t.Run("should not return the kind if not specified", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file"
		content := `
---
version: "2"
`
		_ = afero.WriteFile(memFs, filename, []byte(content), 0644)

		file, err := yaml.NewYamlFile(memFs, filename)

		assert.NoError(t, err)
		actualKind, ok := file.Kind()
		assert.False(t, ok)
		assert.Equal(t, "", actualKind)
	})
}

func TestFileVersion(t *testing.T) {
	t.Run("should return the version if specified", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file"
		content := `
---
version: "2"
`
		_ = afero.WriteFile(memFs, filename, []byte(content), 0644)

		file, err := yaml.NewYamlFile(memFs, filename)

		assert.NoError(t, err)
		actualKind, ok := file.Version()
		assert.True(t, ok)
		assert.Equal(t, "2", actualKind)
	})

	t.Run("should not return the version if not specified", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file"
		content := `
---
kind: SomeKind
`
		_ = afero.WriteFile(memFs, filename, []byte(content), 0644)

		file, err := yaml.NewYamlFile(memFs, filename)

		assert.NoError(t, err)
		actualKind, ok := file.Version()
		assert.False(t, ok)
		assert.Equal(t, "", actualKind)
	})
}

func TestFileReader(t *testing.T) {
	t.Run("should return io.Reader if data is not null", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file"
		content := `some content in file`
		_ = afero.WriteFile(memFs, filename, []byte(content), 0644)

		file, err := yaml.NewYamlFile(memFs, filename)
		if assert.NoError(t, err) {
			actual := file.Reader()
			data := make([]byte, len(content))
			n, err := actual.Read(data)
			assert.NoError(t, err)
			assert.Equal(t, "some content in file", string(data))
			assert.Equal(t, len(content), n)
		}
	})

	t.Run("should return null if data", func(t *testing.T) {
		file := yaml.File{
			Name: "memory",
		}

		actual := file.Reader()
		assert.Nil(t, actual)
	})
}

func TestFileCheck(t *testing.T) {
	t.Run("should return true if the kind and version matches", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file"
		content := `
---
kind: StevedoreManifest
version: "2"
`
		_ = afero.WriteFile(memFs, filename, []byte(content), 0644)

		file, err := yaml.NewYamlFile(memFs, filename)
		assert.NoError(t, err)

		result, err := file.Check(stevedore.KindStevedoreManifest, "2")
		assert.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("should return false if the kind matches and version is less than the supported version", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file"
		content := `
---
kind: StevedoreManifest
version: "1"
`
		_ = afero.WriteFile(memFs, filename, []byte(content), 0644)

		file, err := yaml.NewYamlFile(memFs, filename)
		assert.NoError(t, err)

		result, err := file.Check(stevedore.KindStevedoreManifest, "2")
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return false if the kind matches and version is greater than the supported version", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file"
		content := `
---
kind: StevedoreManifest
version: "3.0.0"
`
		_ = afero.WriteFile(memFs, filename, []byte(content), 0644)

		file, err := yaml.NewYamlFile(memFs, filename)
		assert.NoError(t, err)

		result, err := file.Check(stevedore.KindStevedoreManifest, "2")
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return false if the kind is different and version matches", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file"
		content := `
---
kind: StevedoreEnv
version: "2"
`
		_ = afero.WriteFile(memFs, filename, []byte(content), 0644)

		file, err := yaml.NewYamlFile(memFs, filename)
		assert.NoError(t, err)

		result, err := file.Check(stevedore.KindStevedoreManifest, "2")
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return false if the kind is not provided and version matches", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file"
		content := `
---
version: "2"
`
		_ = afero.WriteFile(memFs, filename, []byte(content), 0644)

		file, err := yaml.NewYamlFile(memFs, filename)
		assert.NoError(t, err)

		result, err := file.Check(stevedore.KindStevedoreManifest, "2")
		assert.NoError(t, err)
		assert.False(t, result)
	})

	t.Run("should return false if the kind is provided and version is not provided", func(t *testing.T) {
		memFs := afero.NewMemMapFs()
		filename := "/mock/file"
		content := `
---
kind: StevedoreEnv
`
		_ = afero.WriteFile(memFs, filename, []byte(content), 0644)

		file, err := yaml.NewYamlFile(memFs, filename)
		assert.NoError(t, err)

		result, err := file.Check(stevedore.KindStevedoreManifest, "2")
		assert.NoError(t, err)
		assert.False(t, result)
	})
}
