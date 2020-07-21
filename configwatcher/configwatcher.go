package configwatcher

import (
	"errors"
	"fiesta/csgo"
	"fiesta/util"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"path/filepath"
)

func addWatchers(watcher *fsnotify.Watcher, userdataDir string) error {
	contents, err := ioutil.ReadDir(userdataDir)
	if err != nil {
		return err
	}
	foundOne := false
	for _, file := range contents {
		if !file.IsDir() {
			continue
		}
		// We are looking for ../userdata/USER_ID/730/local/ since we don't know which steam account the user is logged
		// into but we should be able to just watch all the account directories that have CSGO (730) in them
		localConfigDirPath := filepath.Join(userdataDir, file.Name(), csgo.GameID, csgo.LocalConfigDirName)
		if !util.Exists(localConfigDirPath) {
			continue
		}
		err = watcher.Add(localConfigDirPath)
		if err != nil {
			continue
		}
		foundOne = true
	}
	if !foundOne {
		return errors.New("failed to find a valid CSGO local config directory to watch")
	}
	return nil
}

func Start(userdataDir string, relayFileName string, started chan error, fileChanged chan string, stop chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		started <- fmt.Errorf("failed to start file system watcher: %v", err)
		return
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					done <- true
					return
				}
				if (event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create) &&
					filepath.Base(event.Name) == relayFileName {
					fileChanged <- event.Name
				}
			case _, ok := <-watcher.Errors:
				if !ok {
					done <- true
					return
				}
			case <-stop:
				done <- true
				return
			}
		}
	}()

	err = addWatchers(watcher, userdataDir)
	if err != nil {
		started <- err
		return
	}
	started <- nil
	<-done
	close(fileChanged)
	return
}
