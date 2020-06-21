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
)

func DefaultUserdataDir() string {
	return defaultUserdataDir
}
