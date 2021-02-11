package store

import (
	"os"
	"strings"
)

// Local store to get local configuration
type Local struct{}

// Fetch the local configuration
func (ls Local) Fetch() map[string]interface{} {
	envs := map[string]interface{}{}
	for _, env := range os.Environ() {
		pair := strings.Split(env, "=")
		envs[pair[0]] = pair[1]
	}
	return envs
}

// Cwd Returns the current working directory and error if unable to retrieve it
func (ls Local) Cwd() (string, error) {
	return os.Getwd()
}
