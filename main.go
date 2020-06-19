package main

import (
	"bufio"
	"errors"
	"fiesta/audio"
	"fiesta/configwatcher"
	"fiesta/library"
	"fiesta/settings"
	"fmt"
	"github.com/desertbit/grumble"
	"os"
	"strings"
)

func main() {
	var app = grumble.New(&grumble.Config{
		Name:        "fiesta",
		Description: "short app description",

		Flags: func(f *grumble.Flags) {
			f.String("d", "directory", "DEFAULT", "set an alternative directory path")
			f.Bool("v", "verbose", false, "enable verbose mode")
		},
	})

	config, err := settings.InitializeSettings()
	if err != nil {
		panic(err)
	}

	libPath, err := config.LibraryPath()
	if err != nil {
		panic(err)
	}

	manipulator, err := audio.InitializeAudioManipulator()
	if err != nil {
		panic(err)
	}

	lib, err := library.InitializeLibrary(libPath, manipulator)
	if err != nil {
		panic(err)
	}

	app.AddCommand(&grumble.Command{
		Name:      "import",
		Help:      "import a track",
		Usage:     "import [path]",
		AllowArgs: true,
		Run: func(c *grumble.Context) error {
			if len(c.Args) != 1 {
				return errors.New("incorrect number of arguments passed. import expects one argument")
			}
			return lib.Import(c.Args[0])
		},
	})

	app.AddCommand(&grumble.Command{
		Name:      "list",
		Help:      "list tracks",
		Usage:     "list",
		AllowArgs: false,
		Run: func(c *grumble.Context) error {
			tracks, err := lib.Tracks()
			if err != nil {
				return fmt.Errorf("failed to list tracks: %v", err)
			}
			output := ""
			for index, track := range tracks {
				output += fmt.Sprintf("%v: %v\t%v\n", index+1, track.Name, "["+strings.Join(track.Tags, ", ")+"]")
			}
			_, _ = c.App.Printf(output)
			return nil
		},
	})

	app.AddCommand(&grumble.Command{
		Name:      "start",
		Help:      "start listening for commands from CSGO",
		Usage:     "start",
		AllowArgs: false,
		Run: func(c *grumble.Context) error {
			started := make(chan error)
			fileChanged := make(chan bool)
			stop := make(chan bool)
			go func() {
				configwatcher.Start("/Users/rpecka/Desktop/test.txt", started, fileChanged, stop)
			}()
			err := <-started
			if err != nil {
				return fmt.Errorf("failed to start config watcher: %v", err)
			}
			go func() {
				for _ = range fileChanged {
					// Do something when the file changes
				}
			}()
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Press return to exit…")
			_, err = reader.ReadString('\n')
			if err != nil {
				stop <- true
				return err
			}
			stop <- true
			return nil
		},
	})

	grumble.Main(app)
}
