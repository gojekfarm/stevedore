//go:build !windows && !darwin && !linux
// +build !windows,!darwin,!linux

package credentials

func defaultCredentialsStore() string {
	return ""
}
