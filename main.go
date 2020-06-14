package main

import (
	"fiesta/settings"
	"fmt"

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

	fmt.Println(config.LibraryPath())

	app.AddCommand(&grumble.Command{
		Name: "import",
		Help: "import a track",
		Run: func(c *grumble.Context) error {
			c.App.Println("Not implemented")
			return nil
		},
	})

	grumble.Main(app)
}
