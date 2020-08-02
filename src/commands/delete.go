package commands

import (
	"errors"
	"fiesta/src/library"
	"fmt"
	"github.com/desertbit/grumble"
	"strconv"
)

func addDeleteTrack(app *grumble.App, library library.Library) {
	app.AddCommand(&grumble.Command{
		Name:      "delete",
		Help:      "delete a track",
		Usage:     "delete [track-number]",
		AllowArgs: true,
		Run: func(c *grumble.Context) error {
			if len(c.Args) != 1 {
				return errors.New("incorrect number of arguments passed to delete: " + string(len(c.Args)))
			}
			trackNumber, err := strconv.Atoi(c.Args[0])
			if err != nil {
				return fmt.Errorf("invalid track number: %v", err)
			}
			err = library.DeleteTrack(trackNumber)
			return err
		},
	})
}
