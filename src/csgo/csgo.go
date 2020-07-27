package csgo

import "path/filepath"

const (
	// Userdata
	VoiceInputFileName = "voice_input.wav"
	GameID             = "730"
	LocalConfigDirName = "local"
	CFGDirName         = "cfg"

	// Steamapps
	csgo = "csgo"
	cfg  = "cfg"
)

func PathToCFG(csgoDir string) string {
	return filepath.Join(csgoDir, csgo, cfg)
}
