package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"

	"github.com/fatih/color"
	"golang.org/x/term"
	"gopkg.in/yaml.v2"
)

var whiteColor *color.Color
var redColor *color.Color
var yellowColor *color.Color

func init() {
	whiteColor = color.New(color.FgWhite)
	redColor = color.New(color.FgRed)
	yellowColor = color.New(color.FgYellow)
}

// DieIf invoke os.Exit if there err is not nil
func DieIf(err error, exitFunc func()) {
	if err != nil {
		if exitFunc != nil {
			exitFunc()
		}
		Fatalf("\nCommand failed: %v", err)
	}
}

// RenderAsYaml convert value as yaml and render into writer
func RenderAsYaml(writer io.Writer, value interface{}) {
	bytes, err := yaml.Marshal(value)
	if err != nil {
		return
	}
	content := string(bytes)
	_, _ = fmt.Fprint(writer, content)
}

// NewTableRenderer returns tablewriter.Table
// without any default settings
func NewTableRenderer(writer io.Writer) *tablewriter.Table {
	table := tablewriter.NewWriter(writer)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(false)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(false)
	table.SetNoWhiteSpace(true)
	table.SetTablePadding("\t")
	table.SetRowSeparator("")
	return table
}

// PrintYaml convert value as yaml and render into OutputStream()
func PrintYaml(value interface{}) {
	Info("---")
	RenderAsYaml(OutputStream(), value)
}

// FPrintYaml convert value as yaml and render into given writer
func FPrintYaml(writer io.Writer, value interface{}) {
	FInfo(writer, "---")
	RenderAsYaml(writer, value)
}

// FInfo print values to given writer in white color
func FInfo(writer io.Writer, value ...interface{}) {
	_, _ = whiteColor.Add(color.Bold).Fprint(writer, value...)
	_, _ = fmt.Fprintf(writer, "\n")
}

// Info print values to OutputStream() in white color
func Info(value ...interface{}) {
	stream := OutputStream()
	_, _ = whiteColor.Add(color.Bold).Fprint(stream, value...)
	_, _ = fmt.Fprintf(stream, "\n")
}

// Error print values to os.Stderr in red color
func Error(value ...interface{}) {
	stream := os.Stderr
	_, _ = redColor.Fprint(stream, value...)
	_, _ = fmt.Fprintf(stream, "\n")
}

// Warn print values to OutputStream() in yellow color
func Warn(value ...interface{}) {
	stream := OutputStream()
	_, _ = yellowColor.Fprint(stream, value...)
	_, _ = fmt.Fprintf(stream, "\n")
}

// Fatal print values to os.Stderr in red color
// and invoke os.Exit with non zero status code
func Fatal(value ...interface{}) {
	stream := os.Stderr
	_, _ = redColor.Add(color.Bold).Fprint(stream, value...)
	_, _ = fmt.Fprintf(stream, "\n")
	os.Exit(1)
}

// Infof format and print values to OutputStream() in white color
func Infof(format string, value ...interface{}) {
	Info(fmt.Sprintf(format, value...))
}

// Errorf format and print values to os.Stderr in red color
func Errorf(format string, value ...interface{}) {
	Error(fmt.Sprintf(format, value...))
}

// Warnf format and print values to OutputStream() in yellow color
func Warnf(format string, value ...interface{}) {
	Warn(fmt.Sprintf(format, value...))
}

// Fatalf format and print values to os.Stderr in red color
// and invoke os.Exit with non zero status code
func Fatalf(format string, value ...interface{}) {
	Fatal(fmt.Sprintf(format, value...))
}

// IsTTY returns true if the term is an interactive tty
// and false if not
func IsTTY() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// ColorForced returns true if "FORCE_COLOR" is enabled
// Eg:
// export FORCE_COLOR=true
// (or)
// export FORCE_COLOR=1
func ColorForced() bool {
	colorForced := os.Getenv("FORCE_COLOR")
	if colorForced == "true" || colorForced == "1" {
		return true
	}

	return false
}

// OutputStream gives io.Writer based on interactive term
func OutputStream() io.Writer {
	if IsTTY() {
		return os.Stdout
	}
	return os.Stderr
}
