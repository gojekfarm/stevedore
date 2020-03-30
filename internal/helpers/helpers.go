package helpers

import (
	"fmt"
	"github.com/gojek/stevedore/cmd/cli"
	"github.com/gojek/stevedore/internal/cli/command"
	"github.com/spf13/afero"
	"strings"
)

var testFileDir string

func testDir() string {
	if testFileDir == "" {
		dir, err := afero.TempDir(afero.NewOsFs(), "", "stevedore")
		if err != nil {
			cli.Fatal(fmt.Sprintf("unable to create temp directory '%s', for storing files, reason :%v", testFileDir, err))
		}
		testFileDir = dir
	}
	return testFileDir
}

func execute(command command.Command) (string, error) {
	out, err := command.Execute()
	if err != nil {
		return "", fmt.Errorf("got %v error while executing %s, reason: %s", err, command.Print(), out)
	}
	return string(out), nil
}

func splitByNewLine(in string) []string {
	return strings.Split(strings.Replace(in, "\r\n", "\n", -1), "\n")
}
