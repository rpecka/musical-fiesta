package loader

import (
	"fiesta/configwatcher"
	"fiesta/csgo/configfile"
	"fmt"
)

func Start(relayFilePath string, relayKey string, stop chan bool) error {
	started := make(chan error)
	fileChanged := make(chan bool)
	go func() {
		configwatcher.Start(relayFilePath, started, fileChanged, stop)
	}()
	err := <-started
	if err != nil {
		return fmt.Errorf("failed to start config watcher: %v", err)
	}
	go func() {
		for _ = range fileChanged {
			result, err := configfile.Parse(relayFilePath, relayKey)
			if err != nil {
				fmt.Print(err)
				continue
			}
			fmt.Printf("loading: %v", result.CurrentTrack)
		}
	}()
	return nil
}
