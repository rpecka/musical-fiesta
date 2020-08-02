package loader

import (
	"fiesta/src/configwatcher"
	"fiesta/src/csgo/configfile"
	"fiesta/src/library"
	"fmt"
)

func Start(userdataDir string, trackRelayKey, offsetRelayKey string, stop chan bool, destination string, library library.Library) error {
	started := make(chan error)
	fileChanged := make(chan string)
	go func() {
		configwatcher.Start(userdataDir, configfile.RelayFileName, started, fileChanged, stop)
	}()
	err := <-started
	if err != nil {
		return fmt.Errorf("failed to start config watcher: %v", err)
	}
	go func() {
		for path := range fileChanged {
			result, err := configfile.Parse(path, trackRelayKey, offsetRelayKey)
			if err != nil {
				fmt.Print(err)
				continue
			}
			if result == nil {
				continue
			}
			err = library.Load(result.TrackNumber, result.Offset, destination)
			if err != nil {
				fmt.Printf("failed to load track because: %v\n", err)
			}
		}
	}()
	return nil
}
