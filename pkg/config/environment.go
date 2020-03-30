package config

// Environment represents local configuration
type Environment interface {
	Fetch() map[string]interface{}
	Cwd() (string, error)
}
