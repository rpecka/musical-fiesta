// +build linux

package crossplatform

import (
	"fiesta/src/util"
	"path/filepath"
)

const (
	WhichCommand = "which"
)

// These have not actually been validated -- I just needed them to exist so that we can build and test on Linux in CI
var (
	defaultSteamDir    = filepath.Join(util.UserHomeDir(), ".steam", "steam")
	defaultUserdataDir = filepath.Join(defaultSteamDir, "userdata")
	defaultCSGODir     = filepath.Join(defaultSteamDir, "SteamApps", "common", "Counter-Strike Global Offensive")
)
