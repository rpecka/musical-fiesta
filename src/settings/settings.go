package settings

import (
	"bufio"
	"encoding/json"
	"fiesta/src/crossplatform"
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
}

// Settings : Object to manage the configuration for the app
type realSettings struct {
	path string
}

type settingsFile struct {
	LibraryPath *string `json:"libraryPath,omitempty"`
	UserdataDir *string `json:"userdataDir,omitempty"`
	CSGODir     *string `json:"CSGODir,omitempty"`
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
		defaultDir := crossplatform.DefaultUserdataDir()
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
		defaultDir := crossplatform.DefaultCSGODir()
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
	byteValue, err := json.Marshal(f)
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
