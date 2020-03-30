package stevedore

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
	"io"
)

type marshal = func(in interface{}) ([]byte, error)

func formatAsJSONOrYaml(f fmt.State, c rune, in interface{}) {
	switch c {
	case 'y':
		formatAs(yaml.Marshal, in, f)
	case 'j':
		if f.Flag('#') {
			formatAs(jsonMarshalPrettier, in, f)
		} else {
			formatAs(jsoniter.Marshal, in, f)
		}
	}
}

func formatAs(m marshal, in interface{}, out io.Writer) {
	bytes, err := m(in)
	if err != nil {
		_, _ = out.Write([]byte(fmt.Sprintf("unable to format, reason: %v", err.Error())))
		return
	}
	_, _ = out.Write(bytes)
}

func jsonMarshalPrettier(in interface{}) ([]byte, error) {
	return jsoniter.MarshalIndent(in, "", "  ")
}
