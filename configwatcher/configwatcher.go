package configwatcher

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
)

func Start(path string, started chan error, fileChanged chan bool, stop chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		started <- fmt.Errorf("failed to start file system watcher: %v", err)
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
				if event.Op&fsnotify.Write == fsnotify.Write {
					fileChanged <- true
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

	err = watcher.Add(path)
	if err != nil {
		started <- err
	}
	started <- nil
	<-done
	close(fileChanged)
	return
}
