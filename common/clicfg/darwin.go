//go:build darwin

// Copyright (c) "Neo4j"
// Neo4j Sweden AB [http://neo4j.com]

package clicfg

import (
	"os/user"
	"path/filepath"
)

func init() {
	currentUser, _ := user.Current()
	homeDir := currentUser.HomeDir

	ConfigPrefix = filepath.Join(homeDir, "Library/Preferences")
}
