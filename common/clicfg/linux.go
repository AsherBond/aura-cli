//go:build linux

package clicfg

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func init() {
	xdgConfigHome, found := os.LookupEnv("XDG_CONFIG_HOME")

	if !found || xdgConfigHome == "" {
		xdgConfigHome = "~/.config"
	}

	if strings.HasPrefix(xdgConfigHome, "~/") {
		currentUser, _ := user.Current()
		homeDir := currentUser.HomeDir
		ConfigPrefix = filepath.Join(homeDir, xdgConfigHome[2:])
	} else {
		ConfigPrefix = xdgConfigHome
	}
}
