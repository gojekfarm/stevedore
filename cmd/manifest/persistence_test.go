package manifest

import (
	"os"
	"testing"

	"github.com/gojek/stevedore/cmd/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDiskPersistenceWrite(t *testing.T) {
	t.Run("should write", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		data := []byte("content")

		mockFs := mocks.NewMockFs(ctrl)
		mockFile := mocks.NewMockFile(ctrl)

		mockFs.EXPECT().OpenFile("/mock/temp.yaml", gomock.Any(), os.FileMode(0666)).Return(mockFile, nil)
		mockFile.EXPECT().Write(data).Return(len(data), nil)
		mockFile.EXPECT().Close().Return(nil)

		persistence := DiskPersistence{fs: mockFs}

		err := persistence.write("/mock/temp.yaml", data)

		assert.Nil(t, err)
	})
}
