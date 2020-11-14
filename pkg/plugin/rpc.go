package plugin

import (
	"fmt"
	"net/rpc"
	"time"

	"github.com/hashicorp/go-plugin"
)

func rpcCallWithTimeout(client *rpc.Client, rpcCallFunc func() error) error {
	rpcTimeout := time.Second * 100
	ch := make(chan error, 1)
	go func() { ch <- rpcCallFunc() }()
	select {
	case err := <-ch:
		return err
	case <-time.After(rpcTimeout):
		_ = client.Close()
		return fmt.Errorf("timeout while invoking plugin")
	}
}

// CreateHandshakeConfig creates a handshake configuration
func CreateHandshakeConfig() plugin.HandshakeConfig {
	return plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "BASIC_PLUGIN",
		MagicCookieValue: "hello",
	}
}

// Flags get flags from given interface and populate the resp
func Flags(plugin Interface, resp *[]Flag) error {
	flags, err := plugin.Flags()
	if err != nil {
		return err
	}

	*resp = flags
	return nil
}

// GetType get types from given interface and populate the resp
func GetType(plugin Interface, resp *Type) error {
	flags, err := plugin.Type()
	if err != nil {
		return err
	}

	*resp = flags
	return nil
}

// Version get version from given interface and populate the resp
func Version(plugin Interface, resp *string) error {
	version, err := plugin.Version()
	if err != nil {
		return err
	}

	*resp = version
	return nil
}

// Help get help from given interface and populate the resp
func Help(plugin Interface, resp *string) error {
	help, err := plugin.Help()
	if err != nil {
		return err
	}

	*resp = help
	return nil
}
