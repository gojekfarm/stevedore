package plugin

import (
	"github.com/gojek/stevedore/pkg/config"
	"github.com/gojek/stevedore/pkg/manifest"
)

// Flag represents the plugin flag
type Flag struct {
	Name      string
	Default   string
	Shorthand string
	Required  bool
	Usage     string
}

// Type represents the plugin type
type Type int

const (
	// TypeConfig represents the configuration plugin type
	TypeConfig Type = iota
	// TypeManifest represents the manifest plugin type
	TypeManifest
	// TypeContext represents the context plugin type
	TypeContext
)

// String returns the Type's string representation.
func (t Type) String() string {
	switch t {
	case TypeConfig:
		return "config"
	case TypeManifest:
		return "manifest"
	case TypeContext:
		return "context"
	}

	return ""
}

// Interface represents the plugin Interface
type Interface interface {
	Version() (string, error)
	Flags() ([]Flag, error)
	Type() (Type, error)
	Help() (string, error)
	Close() error
}

// ConfigInterface represents the configuration plugin
type ConfigInterface interface {
	Interface
	config.Provider
}

// ManifestInterface represents the configuration plugin
type ManifestInterface interface {
	Interface
	manifest.Provider
}
