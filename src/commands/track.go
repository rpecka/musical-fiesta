package commands

import (
	"errors"
	"fiesta/src/library"
	"fmt"
	"github.com/desertbit/grumble"
	"strconv"
)

func extractTrackNumberArgument(ctx *grumble.Context) (int, error) {
	if len(ctx.Args) < 1 {
		return -1, errors.New("no track number provided")
	}
	trackString := (ctx.Args)[0]
	trackNumber, err := strconv.Atoi(trackString)
	if err != nil {
		return -1, errors.New("track number must be a number")
	}
	ctx.Args = ctx.Args[1:]
	return trackNumber, nil
}

func listTags(app *grumble.App, name string, tags []string) {
	output := fmt.Sprintf("tags for \"%s\":\n", name)
	for idx, tag := range tags {
		output += fmt.Sprintf("\t%d. %s", idx + 1, tag)
	}
	app.Printf(output)
}

func addTrack(app *grumble.App, library *library.Library) {
	trackCommand := grumble.Command{
		Name:      "track",
		Help:      "commands associated with manipulating tracks",
		Usage:     "track [command...]",
		AllowArgs: true,
		Run: nil,
		Completer: nil,
	}

	trackCommand.AddCommand(&grumble.Command{
		Name:      "list-tags",
		Help:      "list all of a track's tags",
		Usage:     "track list-tags [track-number]",
		AllowArgs: true,
		Run: func(c *grumble.Context) error {
			trackNumber, err := extractTrackNumberArgument(c)
			if err != nil {
				return err
			}
			track, err := (*library).GetTrack(trackNumber)
			if err != nil {
				return fmt.Errorf("failed to get track: %v", err)
			}
			listTags(app, track.Name, track.Tags)
			return nil
		},
		Completer: nil,
	})

	trackCommand.AddCommand(&grumble.Command{
		Name:      "add-tag",
		Help:      "add a tag to the track",
		Usage:     "track add-tag [track-number] [tag]",
		AllowArgs: true,
		Run: func(c *grumble.Context) error {
			trackNumber, err := extractTrackNumberArgument(c)
			if err != nil {
				return err
			}
			if len(c.Args) != 1 {
				return errors.New("incorrect number of arguments provided")
			}
			tagString := c.Args[0]
			err = (*library).AddTag(trackNumber, tagString)
			return err
		},
		Completer: nil,
	})

	trackCommand.AddCommand(&grumble.Command{
		Name:      "delete-tag",
		Help:      "delete a tag on a track",
		Usage:     "track delete-tag [track-number] [tag-number]",
		AllowArgs: true,
		Run: func(c *grumble.Context) error {
			trackNumber, err := extractTrackNumberArgument(c)
			if err != nil {
				return err
			}
			if len(c.Args) != 1 {
				return errors.New("incorrect number of arguments provided")
			}
			tagString := c.Args[0]
			tagNumber, err := strconv.Atoi(tagString)
			if err != nil {
				return errors.New("tag number must be a number")
			}
			err = (*library).DeleteTag(trackNumber, tagNumber)
			return err
		},
		Completer: nil,
	})
	app.AddCommand(&trackCommand)
}
