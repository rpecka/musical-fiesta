// +build windows

package defaults

import (
	"path/filepath"
)

const (
	WhichCommand = "where"
)

var (
	defaultSteamDir    = filepath.Join("C:\\", "Program Files (x86)", "Steam")
	defaultUserdataDir = filepath.Join(defaultSteamDir, "userdata")
	defaultCSGODir     = filepath.Join(defaultSteamDir, "steamapps", "common", "Counter-Strike Global Offensive")
)
