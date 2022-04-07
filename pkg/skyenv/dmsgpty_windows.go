//go:build windows
// +build windows

package skyenv

import (
	"os"
	"path/filepath"
)

// CLIAddr gets the default cli address
func CLIAddr() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		homedir = os.TempDir()
	}
	return filepath.Join(homedir, "dmsgpty.sock")
}
