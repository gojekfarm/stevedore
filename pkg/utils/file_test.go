package utils_test

import (
	"os"
	"path"
	"testing"

	"github.com/gojek/stevedore/pkg/utils"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestFileUtilsTempDir(t *testing.T) {
	t.Run("should be able to create temp dir", func(t *testing.T) {
		fileUtils := utils.NewOsFileUtils()

		tempDir, err := fileUtils.TempDir("stevedore")
		assert.NoError(t, err)

		exists, err := afero.Exists(afero.NewOsFs(), tempDir)
		assert.True(t, exists)
		assert.NoError(t, err)

		err = os.Remove(tempDir)
		assert.NoError(t, err)
	})

	t.Run("should be able to write file", func(t *testing.T) {
		fileUtils := utils.NewOsFileUtils()
		contents := []byte("random-string")
		tempDir, _ := fileUtils.TempDir("stevedore")
		fileName := path.Join(tempDir, "test-file.txt")

		err := fileUtils.WriteFile(fileName, contents, os.ModeTemporary)
		assert.NoError(t, err)

		exists, err := afero.Exists(afero.NewOsFs(), fileName)
		assert.True(t, exists)
		assert.NoError(t, err)

		err = os.Remove(fileName)
		assert.NoError(t, err)

		err = os.Remove(tempDir)
		assert.NoError(t, err)
	})
}
