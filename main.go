package main

import (
	"errors"
	"fiesta/library"
	"fiesta/settings"
	"github.com/desertbit/grumble"
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

	lib, err := library.InitializeLibrary(libPath)
	if err != nil {
		panic(err)
	}

	app.AddCommand(&grumble.Command{
		Name: "import",
		Help: "import a track",
		Usage: "import [path]",
		AllowArgs: true,
		Run: func(c *grumble.Context) error {
			if len(c.Args) != 1 {
				return errors.New("incorrect number of arguments passed. import expects one argument")
			}
			return lib.Import(c.Args[0])
		},
	})

	grumble.Main(app)
}
