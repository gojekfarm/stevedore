package cmd

import (
	"os"
	"strings"
)

// Thanks to `docker-machine` for the idea
// Reference: https://github.com/docker/machine/blob/master/commands/create.go
func flagHackLookup(flagName string) string {
	// e.g. "-d" for "--driver"
	flagPrefix := flagName[1:3]

	for i, arg := range os.Args {
		if strings.Contains(arg, flagPrefix) {
			// format '--driver foo' or '-d foo'
			if arg == flagPrefix || arg == flagName {
				if i+1 < len(os.Args) {
					return os.Args[i+1]
				}
			}

			// format '--driver=foo' or '-d=foo'
			if strings.HasPrefix(arg, flagPrefix+"=") || strings.HasPrefix(arg, flagName+"=") {
				return strings.Split(arg, "=")[1]
			}
		}
	}

	return ""
}
