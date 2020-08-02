package settings

import (
	"bufio"
	"encoding/json"
	"fiesta/src/defaults"
	"fiesta/src/util"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

func InitializeSettings() (Settings, error) {
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

	if settingsFile.LibraryPath == nil {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Please provide the path where your library will be created: ")
		path, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		cleanPath := strings.TrimSpace(path)
		libraryPath, _ := homedir.Expand(filepath.Join(cleanPath, libraryDirName))
		fmt.Println("Your library will be created at: " + libraryPath)
		err = os.MkdirAll(libraryPath, 0755)
		if err != nil {
			return nil, err
		}
		settingsFile.LibraryPath = &libraryPath
		err = settings.writeSettings(settingsFile)
		if err != nil {
			return nil, err
		}
	}

	if settingsFile.UserdataDir == nil {
		reader := bufio.NewReader(os.Stdin)
		defaultDir := defaults.DefaultUserdataDir()
		fmt.Print("Please provide the path to your Steam userdata directory. Press return to use the default " +
			"(" + defaultDir + "): ")
		path, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		cleanPath := strings.TrimSpace(path)
		userdataPath, _ := homedir.Expand(cleanPath)
		if userdataPath == "" {
			settingsFile.UserdataDir = &defaultDir
		} else {
			settingsFile.UserdataDir = &userdataPath
		}
		err = settings.writeSettings(settingsFile)
		if err != nil {
			return nil, err
		}
	}

	if settingsFile.CSGODir == nil {
		reader := bufio.NewReader(os.Stdin)
		defaultDir := defaults.DefaultCSGODir()
		fmt.Print("Please provide the path to your CSGO game files directory. Press return to use the default " +
			"(" + defaultDir + "): ")
		path, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		cleanPath := strings.TrimSpace(path)
		csgoPath, _ := homedir.Expand(cleanPath)
		if csgoPath == "" {
			settingsFile.CSGODir = &defaultDir
		} else {
			settingsFile.CSGODir = &csgoPath
		}
		err = settings.writeSettings(settingsFile)
		if err != nil {
			return nil, err
		}
	}

	if settingsFile.TrackRelayKey == nil {
		reader := bufio.NewReader(os.Stdin)
		defaultKey := defaults.DefaultTrackRelayKey
		fmt.Print("Please enter the bind code for a key you do not use in CSGO. Press return to use the default " +
			"(" + defaultKey + "): ")
		key, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		cleanKey := strings.TrimSpace(key)
		if cleanKey == "" {
			settingsFile.TrackRelayKey = &defaultKey
		} else {
			settingsFile.OffsetRelayKey = &cleanKey
		}
		err = settings.writeSettings(settingsFile)
		if err != nil {
			return nil, err
		}
	}

	if settingsFile.OffsetRelayKey == nil {
		reader := bufio.NewReader(os.Stdin)
		defaultKey := defaults.DefaultOffsetRelayKey
		fmt.Print("Please enter the bind code for a key you do not use in CSGO. Press return to use the default " +
			"(" + defaultKey + "): ")
		key, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		cleanKey := strings.TrimSpace(key)
		if cleanKey == "" {
			settingsFile.OffsetRelayKey = &defaultKey
		} else {
			settingsFile.OffsetRelayKey = &cleanKey
		}
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
	json.Unmarshal(byteValue, &settingsFile)
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
