package settings

import (
	"bufio"
	"encoding/json"
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
	settingsDirPath  = filepath.Join(forceUserHomeDir(), settingsDirName)
	settingsFilePath = filepath.Join(settingsDirPath, settingsFileName)
)

type Settings interface {
	LibraryPath() (string, error)
}

// Settings : Object to manage the configuration for the app
type realSettings struct {
	path string
}

type settingsFile struct {
	LibraryPath *string `json:"libraryPath,omitempty"`
}

func InitializeSettings() (Settings, error) {
	if !exists(settingsDirPath) {
		err := os.Mkdir(settingsDirPath, 0755)
		if err != nil {
			return nil, err
		}
	}

	if !exists(settingsFilePath) {
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

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func forceUserHomeDir() string {
	dir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return dir
}
