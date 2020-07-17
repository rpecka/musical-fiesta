// +build windows

package crossplatform

import (
	"path/filepath"
)

const (
	WhichCommand = "where"
)

var (
	defaultSteamDir    = filepath.Join("C:\\", "Program Files (x86)", "Steam")
	defaultUserdataDir = filepath.Join(defaultSteamDir, "userdata")
)

func DefaultUserdataDir() string {
	return defaultUserdataDir
}
