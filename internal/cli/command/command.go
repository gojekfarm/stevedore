package command

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Executable represents cli command which can be executed
type Executable interface {
	Command() (string, []string)
	Print() string
	Environ() []string
}

// Command holds a executable command, stdout, stderr
type Command struct {
	executable Executable
}

// Execute executes the underlying command, returns error if any
func (command Command) Execute() ([]byte, error) {
	envs := os.Environ()
	envs = append(envs, command.executable.Environ()...)

	name, args := command.executable.Command()
	cmd := exec.Command(name, args...)
	cmd.Env = envs
	return cmd.CombinedOutput()
}

// Print returns executable command as string
func (command Command) Print() string {
	return command.executable.Print()
}

// NewCommand creates a new instance of executable command
func NewCommand(executable Executable) Command {
	return Command{
		executable: executable,
	}
}

// Print the underlying command
func Print(executable Executable) string {
	cmd, args := executable.Command()
	buffer := strings.Builder{}
	buffer.WriteString(fmt.Sprintf("%s ", cmd))
	buffer.WriteString(strings.Join(args, " "))
	return buffer.String()
}
