package main

import (
	"errors"
	"fiesta/src/audio"
	"fiesta/src/commands"
	"fiesta/src/library"
	"fiesta/src/settings"
	"fmt"
	"github.com/desertbit/grumble"
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

	config, err := settings.InitializeSettings(app)
	if err != nil {
		panic(err)
	}

	libPath, err := config.LibraryPath()
	if err != nil {
		panic(err)
	}

	manipulator, err := audio.InitializeManipulator()
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

	commands.SetUpCommands(app, config, lib)

	grumble.Main(app)
}
