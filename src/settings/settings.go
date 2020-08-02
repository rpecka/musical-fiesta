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

func promptForInput(printer Printer, reader *bufio.Reader, prompt string) (string, error) {
	printer.Printf(prompt)
	return reader.ReadString('\n')
}

func promptForPath(printer Printer, reader *bufio.Reader, prompt string) (string, error) {
	path, err := promptForInput(printer, reader, prompt)
	if err != nil {
		return "", err
	}
	trimmedPath := strings.TrimSpace(path)
	expandedPath, err := homedir.Expand(trimmedPath)
	if err != nil {
		return "", err
	}
	return expandedPath, nil
}

func promptForPathDefault(printer Printer, reader *bufio.Reader, prompt string, defaultValue string) (string, error) {
	path, err := promptForPath(printer, reader, prompt)
	if err != nil {
		return "", err
	}
	if path == "" {
		return defaultValue, nil
	} else {
		return path, nil
	}
}

func promptForKeyDefault(printer Printer, reader *bufio.Reader, prompt string, defaultValue string) (string, error) {
	input, err := promptForInput(printer, reader, prompt)
	input = strings.TrimSpace(input)
	if err != nil {
		return "", err
	}
	if input == "" {
		return defaultValue, nil
	} else {
		return input, nil
	}
}

func executePromptForPathSetting(promptLogic func() (string, error), setting *string, pathHandler func (*string) error) error {
	path, err := promptLogic()
	if err != nil {
		return err
	}
	if pathHandler != nil {
		err = pathHandler(&path)
		if err != nil {
			return err
		}
	}
	setting = &path
	return nil
}

func promptForPathSetting(printer Printer, reader *bufio.Reader, prompt string, setting *string, pathHandler func (*string) error) error {
	return executePromptForPathSetting(func() (string, error) {
		return promptForPath(printer, reader, prompt)
	}, setting, pathHandler)
}

func promptForPathSettingDefault(printer Printer, reader *bufio.Reader, prompt string, defaultValue string, setting *string, pathHandler func (*string) error) error {
	return executePromptForPathSetting(func() (string, error) {
		return promptForPathDefault(printer, reader, prompt, defaultValue)
	}, setting, pathHandler)
}

func promptForKeySettingDefault(printer Printer, reader *bufio.Reader, prompt string, defaultValue string, setting *string) error {
	key, err := promptForKeyDefault(printer, reader, prompt, defaultValue)
	if err != nil {
		return err
	}
	setting = &key
	return nil
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
	madeChange := false
	if settingsFile.LibraryPath == nil {
		madeChange = true
		err = promptForPathSetting(printer, reader, "Please provide the path where your library will be created: ", settingsFile.LibraryPath, func(path *string) error {
			*path = filepath.Join(*path, libraryDirName)
			fmt.Println("Your library will be created at: " + *path)
			return os.MkdirAll(*path, 0755)
		})
		if err != nil {
			return nil, err
		}
	}

	if settingsFile.UserdataDir == nil {
		madeChange = true
		defaultDir := defaults.DefaultUserdataDir()
		err = promptForPathSettingDefault(printer, reader, "Please provide the path to your Steam userdata directory. "+
			"Press return to use the default ("+defaultDir+"): ", defaultDir, settingsFile.UserdataDir, nil)
		if err != nil {
			return nil, err
		}
	}

	if settingsFile.CSGODir == nil {
		madeChange = true
		defaultDir := defaults.DefaultCSGODir()
		err = promptForPathSettingDefault(printer, reader, "Please provide the path to your CSGO game files directory. Press return to use the default "+
			"("+defaultDir+"): ", defaultDir, settingsFile.CSGODir, nil)
		if err != nil {
			return nil, err
		}
	}

	if settingsFile.TrackRelayKey == nil {
		madeChange = true
		defaultKey := defaults.DefaultTrackRelayKey
		err = promptForKeySettingDefault(printer, reader, "Please enter the bind code for a key you do not use in CSGO. Press return to use the default "+
			"("+defaultKey+"): ", defaultKey, settingsFile.TrackRelayKey)
		if err != nil {
			return nil, err
		}
	}

	if settingsFile.OffsetRelayKey == nil {
		madeChange = true
		defaultKey := defaults.DefaultOffsetRelayKey
		err = promptForKeySettingDefault(printer, reader, "Please enter the bind code for a key you do not use in CSGO. Press return to use the default "+
			"("+defaultKey+"): ", defaultKey, settingsFile.OffsetRelayKey)
		if err != nil {
			return nil, err
		}
	}

	if madeChange {
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
