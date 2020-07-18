// +build darwin

package crossplatform

import (
	"fiesta/util"
	"path/filepath"
)

const (
	WhichCommand = "which"
)

var (
	defaultSteamDir    = filepath.Join(util.UserHomeDir(), "Library", "Application Support", "Steam")
	defaultUserdataDir = filepath.Join(defaultSteamDir, "userdata")
	defaultCSGODir     = filepath.Join(defaultSteamDir, "steamapps", "common", "Counter-Strike Global Offensive")
)
