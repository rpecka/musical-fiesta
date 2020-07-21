package commands

import (
	"errors"
	"fiesta/src/library"
	"fmt"
	"github.com/desertbit/grumble"
	"strconv"
)

const (
	addFlag = "add"
	addShort = "a"
	deleteFlag = "delete"
	deleteShort = "d"
)

func listTags(app *grumble.App, name string, tags []string) {
	output := fmt.Sprintf("tags for \"%s\":\n", name)
	for idx, tag := range tags {
		output += fmt.Sprintf("\t%d. %s\n", idx + 1, tag)
	}
	app.Printf(output)
}

func addTag(app *grumble.App, library library.Library) {
	app.AddCommand(&grumble.Command{
		Name:      "tag",
		Help:      "view or edit the tags of a track",
		Usage:     "tag [track-number]\n" +
			"\tor: tag [track-number] -a [new-tag]\n" +
			"\tor: tag [track-number] -d [tag-number]",
		Flags: func(f *grumble.Flags) {
			f.String(addShort, addFlag, "", "the tag to add to the track")

			// This is a string so that we can support deleting by the tag instead of tag number in the future
			f.String(deleteShort, deleteFlag, "", "the tag number to delete")
		},
		AllowArgs: true,
		Run: func(c *grumble.Context) error {
			if len(c.Args) != 1 {
				return errors.New("incorrect number of arguments provided")
			}
			trackNumber, err := strconv.Atoi(c.Args[0])
			if err != nil {
				return fmt.Errorf("the track number must be an integer: `%v`", c.Args[0])
			}

			addInput := c.Flags.String(addFlag)
			addRequested := addInput != ""

			deleteInput := c.Flags.String(deleteFlag)
			deleteRequested := deleteInput != ""

			if addRequested && deleteRequested {
				return fmt.Errorf("cannot simultaneously add and delete a tag")
			}

			if !addRequested && !deleteRequested {
				track, err := library.GetTrack(trackNumber)
				if err != nil {
					return fmt.Errorf("failed to get track: %v", err)
				}
				listTags(app, track.Name, track.Tags)
			} else if addRequested {
				err = library.AddTag(trackNumber, addInput)
				if err != nil {
					return err
				}
			} else {  // deleteRequested

			}
			return nil
		},
	})
}
