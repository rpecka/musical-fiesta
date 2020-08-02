package loader

import (
	"errors"
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
			switch result.ResultType {
			case configfile.LoadNumberResult:
				err = library.Load(result.TrackNumber, result.Offset, destination)
			case configfile.TagResult:
				err = errors.New("string commands are not yet supported")
			}
			if err != nil {
				fmt.Printf("failed to load track because: %v\n", err)
			}
		}
	}()
	return nil
}
