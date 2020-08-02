package settings

import (
	"bufio"
	"encoding/json"
	"fiesta/src/defaults"
	"fiesta/src/util"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	settingsDirName  = ".fiesta"
	settingsFileName = "settings.json"
	libraryDirName   = "library"
)

var (
	settingsDirPath  = filepath.Join(util.UserHomeDir(), settingsDirName)
	settingsFilePath = filepath.Join(settingsDirPath, settingsFileName)
)

type Settings interface {
	LibraryPath() (string, error)
	UserdataDirPath() (string, error)
	CSGODirPath() (string, error)
	TrackRelayKey() (string, error)
	OffsetRelayKey() (string, error)
}

type Printer interface {
	Printf(format string, args ...interface{}) (int, error)
}

// Settings : Object to manage the configuration for the app
type realSettings struct {
	path string
}

type settingsFile struct {
	LibraryPath    *string `json:"libraryPath,omitempty"`
	UserdataDir    *string `json:"userdataDir,omitempty"`
	CSGODir        *string `json:"CSGODir,omitempty"`
	TrackRelayKey  *string `json:"trackRelayKey,omitempty"`
	OffsetRelayKey *string `json:"offsetRelayKey,omitempty"`
}

func InitializeSettings(printer Printer) (Settings, error) {
	if !util.Exists(settingsDirPath) {
		err := os.Mkdir(settingsDirPath, 0755)
		if err != nil {
			return nil, err
		}
	}

	if !util.Exists(settingsFilePath) {
		file, err := os.Create(settingsFilePath)
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	settings := realSettings{settingsFilePath}

	settingsFile, err := settings.parseSettings()
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(os.Stdin)
	needsWrite := false

	err = promptForPathSettingIfNeeded(printer, reader, "Please provide the path where your library will be created: ", settingsFile.LibraryPath, func(path *string) error {
		*path = filepath.Join(*path, libraryDirName)
		fmt.Println("Your library will be created at: " + *path)
		return os.MkdirAll(*path, 0755)
	}, &needsWrite)
	if err != nil {
		return nil, err
	}

	defaultUserdataDir := defaults.DefaultUserdataDir()
	err = promptForPathSettingDefaultIfNeeded(printer, reader, "Please provide the path to your Steam userdata directory. "+
		"Press return to use the default ("+defaultUserdataDir+"): ", defaultUserdataDir, settingsFile.UserdataDir, nil, &needsWrite)
	if err != nil {
		return nil, err
	}

	defaultCSGODir := defaults.DefaultCSGODir()
	err = promptForPathSettingDefaultIfNeeded(printer, reader, "Please provide the path to your CSGO game files directory. Press return to use the default "+
		"("+defaultCSGODir+"): ", defaultCSGODir, settingsFile.CSGODir, nil, &needsWrite)
	if err != nil {
		return nil, err
	}

	defaultTrackRelayKey := defaults.DefaultTrackRelayKey
	err = promptForKeySettingDefaultIfNeeded(printer, reader, "Please enter the bind code for a key you do not use in CSGO. Press return to use the default "+
		"("+defaultTrackRelayKey+"): ", defaultTrackRelayKey, settingsFile.TrackRelayKey, &needsWrite)
	if err != nil {
		return nil, err
	}

	defaultKey := defaults.DefaultOffsetRelayKey
	err = promptForKeySettingDefaultIfNeeded(printer, reader, "Please enter the bind code for a key you do not use in CSGO. Press return to use the default "+
		"("+defaultKey+"): ", defaultKey, settingsFile.OffsetRelayKey, &needsWrite)
	if err != nil {
		return nil, err
	}

	if needsWrite {
		err = settings.writeSettings(settingsFile)
		if err != nil {
			return nil, err
		}
	}

	return settings, nil
}

func (s realSettings) parseSettings() (*settingsFile, error) {
	file, err := os.Open(s.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	byteValue, err := ioutil.ReadAll(file)
	var settingsFile settingsFile
	err = json.Unmarshal(byteValue, &settingsFile)
	if err != nil {
		return nil, err
	}
	return &settingsFile, nil
}

func (s realSettings) writeSettings(f *settingsFile) error {
	byteValue, err := json.MarshalIndent(f, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(settingsFilePath, byteValue, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (s realSettings) LibraryPath() (string, error) {
	settings, err := s.parseSettings()
	if err != nil {
		return "", err
	}
	return *settings.LibraryPath, nil
}

func (s realSettings) UserdataDirPath() (string, error) {
	settings, err := s.parseSettings()
	if err != nil {
		return "", err
	}
	return *settings.UserdataDir, nil
}

func (s realSettings) CSGODirPath() (string, error) {
	settings, err := s.parseSettings()
	if err != nil {
		return "", err
	}
	return *settings.CSGODir, nil
}

func (s realSettings) TrackRelayKey() (string, error) {
	settings, err := s.parseSettings()
	if err != nil {
		return "", err
	}
	return *settings.TrackRelayKey, nil
}

func (s realSettings) OffsetRelayKey() (string, error) {
	settings, err := s.parseSettings()
	if err != nil {
		return "", err
	}
	return *settings.OffsetRelayKey, nil
}
