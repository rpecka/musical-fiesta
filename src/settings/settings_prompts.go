package settings

import (
	"bufio"
	"github.com/mitchellh/go-homedir"
	"strings"
)

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

func executePromptForPathSetting(promptLogic func() (string, error), setting *string, pathHandler func(*string) error) error {
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

func executeIfNil(setting *string, logic func() error, updated *bool) error {
	if setting != nil {
		return nil
	}
	err := logic()
	if err != nil {
		return err
	}
	if updated != nil {
		*updated = true
	}
	return nil
}

func promptForPathSettingIfNeeded(printer Printer, reader *bufio.Reader, prompt string, setting *string, pathHandler func(*string) error, updated *bool) error {
	return executeIfNil(setting, func() error {
		return executePromptForPathSetting(func() (string, error) {
			return promptForPath(printer, reader, prompt)
		}, setting, pathHandler)
	}, updated)
}

func promptForPathSettingDefaultIfNeeded(printer Printer, reader *bufio.Reader, prompt string, defaultValue string, setting *string, pathHandler func(*string) error, updated *bool) error {
	return executeIfNil(setting, func() error {
		return executePromptForPathSetting(func() (string, error) {
			return promptForPathDefault(printer, reader, prompt, defaultValue)
		}, setting, pathHandler)
	}, updated)
}

func promptForKeySettingDefaultIfNeeded(printer Printer, reader *bufio.Reader, prompt string, defaultValue string, setting *string, updated *bool) error {
	return executeIfNil(setting, func() error {
		key, err := promptForKeyDefault(printer, reader, prompt, defaultValue)
		if err != nil {
			return err
		}
		setting = &key
		return nil
	}, updated)
}
