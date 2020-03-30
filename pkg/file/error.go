package file

import (
	"bytes"
	"fmt"
)

// Error represents error related to file
type Error struct {
	Reason   error
	Filename string
}

// Error wraps custom message on the underlying error
func (err Error) Error() string {
	return fmt.Sprintf("Error while processing %s\n%s", err.Filename, err.Reason)
}

// Errors represents error related to files
type Errors []Error

// Error wraps custom message on the underlying file error
func (errors Errors) Error() string {
	buff := bytes.NewBufferString(fmt.Sprintf("Failed with %d Errors\n", len(errors)))
	for i, err := range errors {
		buff.WriteString(fmt.Sprintf("\n%d. %s\n", i+1, err.Error()))
	}
	return buff.String()
}
