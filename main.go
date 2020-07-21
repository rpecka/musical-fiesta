package main

import (
	"bufio"
	"errors"
	"fiesta/audio"
	"fiesta/commands"
	"fiesta/csgo"
	"fiesta/library"
	"fiesta/loader"
	"fiesta/settings"
	"fmt"
	"github.com/desertbit/grumble"
	"os"
	"path/filepath"
	"strings"
)

const (
	directoryFlag = "directory"
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
		Name: "import",
		Help: "import a track",
		Usage: "import [path]\n" +
			"\tor: -d [directory-path]",
		AllowArgs: true,
		Flags: func(f *grumble.Flags) {
			f.Bool("d", directoryFlag, false, "import a directory of tracks")
		},
		Run: func(c *grumble.Context) error {
			if len(c.Args) != 1 {
				return errors.New("incorrect number of arguments passed. import expects one argument")
			}
			path := c.Args[0]
			if c.Flags.Bool(directoryFlag) {
				failures, err := lib.ImportDir(path)
				if err != nil {
					return err
				}
				if len(failures) > 0 {
					_, _ = c.App.Println("Failures:")
					for _, failure := range failures {
						_, _ = c.App.Println("\t" + failure)
					}
				}
				return nil
			} else {
				return lib.Import(path)
			}

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
			stop := make(chan bool)
			userdataDir, err := config.UserdataDirPath()
			if err != nil {
				return err
			}
			csgoDir, err := config.CSGODirPath()
			if err != nil {
				return err
			}
			destination := filepath.Join(csgoDir, csgo.VoiceInputFileName)
			err = loader.Start(userdataDir, "z", stop, destination, lib)
			if err != nil {
				return err
			}
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

	commands.SetUpCommands(app, lib)

	grumble.Main(app)
}
