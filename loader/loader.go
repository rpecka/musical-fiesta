package loader

import (
	"fiesta/configwatcher"
	"fiesta/csgo/configfile"
	"fiesta/library"
	"fiesta/util"
	"fmt"
)

func load(trackNumber int, destination string, library library.Library) error {
	track, err := library.GetTrack(trackNumber)
	if err != nil {
		return err
	}
	err = util.CopyFile(track.Path, destination)
	return err
}

func Start(userdataDir string, relayKey string, stop chan bool, destination string, library library.Library) error {
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
			result, err := configfile.Parse(path, relayKey)
			if err != nil {
				fmt.Print(err)
				continue
			}
			err = load(result.CurrentTrack, destination, library)
			if err != nil {
				fmt.Printf("failed to load track because: %v\n", err)
			}
		}
	}()
	return nil
}
