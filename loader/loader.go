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

func Start(relayFilePath string, relayKey string, stop chan bool, destination string, library library.Library) error {
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
			err = load(result.CurrentTrack, destination, library)
			if err != nil {
				fmt.Printf("failed to load track because: %v\n", err)
			}
		}
	}()
	return nil
}
