package helm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetHelmVersion(t *testing.T) {
	version := "2.7"
	originalVersion := GetHelmVersion()
	defer func() {
		SetHelmVersion(originalVersion)
	}()

	SetHelmVersion(version)

	actual := GetHelmVersion()

	assert.Equal(t, version, actual)
}

func TestSetBuildMetadata(t *testing.T) {
	meta := "go-1.11.2"
	originalMetadata := GetBuildMetadata()
	defer func() {
		SetBuildMetadata(originalMetadata)
	}()

	SetBuildMetadata(meta)

	actual := GetBuildMetadata()

	assert.Equal(t, meta, actual)
}
