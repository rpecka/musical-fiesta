// +build windows

package crossplatform

import (
	"path/filepath"
)

const (
	WhichCommand = "where"
)

var (
	steamDir    = filepath.Join("C:\\", "Program Files (x86)", "Steam")
	userdataDir = filepath.Join(windowsSteamDir, "userdata")
)
