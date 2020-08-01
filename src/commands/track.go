package commands

import (
	"bufio"
	"errors"
	"fiesta/src/library"
	"fiesta/src/loader"
	"fmt"
	"github.com/desertbit/grumble"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	startTimeFlag      = "start"
	startTimeFlagShort = "s"
	endTimeFlag        = "end"
	endTimeFlagShort   = "e"
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
		output += fmt.Sprintf("\t%d. %s", idx+1, tag)
	}
	app.Printf(output)
}

func addTrack(app *grumble.App, library *library.Library) {
	trackCommand := grumble.Command{
		Name:      "track",
		Help:      "commands associated with manipulating tracks",
		Usage:     "track [command...]",
		AllowArgs: true,
		Run:       nil,
		Completer: nil,
	}

	trackCommand.AddCommand(&grumble.Command{
		Name:      "info",
		Help:      "show information for a track",
		Usage:     "track info [track-number]",
		AllowArgs: true,
		Run: func(c *grumble.Context) error {
			trackNumber, err := extractTrackNumberArgument(c)
			if err != nil {
				return err
			}
			track, err := (*library).GetTrack(trackNumber)
			if err != nil {
				return err
			}
			output := fmt.Sprintf("Name:\t%s\nPath:\t%s\n", track.Name, track.Path)
			output += "Tags:\t[" + strings.Join(track.Tags, ", ") + "]\n"
			trim := track.Trim
			output += "Trim:\t"
			if trim != nil {
				output += fmt.Sprintf("start: %v\tend: %v\n", trim.Start, trim.End)
			} else {
				output += "None\n"
			}
			app.Printf(output)
			return nil
		},
		Completer: nil,
	})

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

	trackCommand.AddCommand(&grumble.Command{
		Name: "trim",
		Help: "trim a track's start and end times",
		Usage: "track trim --start [start-seconds] --end [end-seconds] [track-number] \n" +
			"or:\ttrack trim --start [start-seconds] [track-number] \n" +
			"or:\ttrack trim --end [end-seconds] [track-number]",
		Flags: func(f *grumble.Flags) {
			f.Duration(startTimeFlagShort, startTimeFlag, -1, "the start time in seconds. "+
				"Negative values indicate the start of the track")
			f.Duration(endTimeFlagShort, endTimeFlag, -1, "the end time in seconds. "+
				"Negative values indicate the end of the track")
		},
		AllowArgs: true,
		Run: func(c *grumble.Context) error {
			trackNumber, err := extractTrackNumberArgument(c)
			if err != nil {
				return err
			}
			startInput := c.Flags.Duration(startTimeFlag)
			endInput := c.Flags.Duration(endTimeFlag)

			var startTime *time.Duration
			var endTime *time.Duration

			if startInput < 0 {
				startTime = nil
			} else {
				startTime = &startInput
			}

			if endInput < 0 {
				endTime = nil
			} else {
				endTime = &endInput
			}

			if startTime == nil && endTime == nil {
				return errors.New("either a start time or an end time or both must be provided\n" +
					"use track clear-trim to reset trim settings")
			}

			return (*library).TrimTrack(trackNumber, startTime, endTime)
		},
		Completer: nil,
	})

	trackCommand.AddCommand(&grumble.Command{
		Name:      "clear-trim",
		Help:      "erase trim settings for a track",
		Usage:     "track clear-trim [track-number]",
		AllowArgs: true,
		Run: func(c *grumble.Context) error {
			trackNumber, err := extractTrackNumberArgument(c)
			if err != nil {
				return err
			}
			return (*library).ClearTrim(trackNumber)
		},
		Completer: nil,
	})

	trackCommand.AddCommand(&grumble.Command{
		Name:      "test",
		Help:      "listen to a track to make sure that the volume and trimming is correct",
		Usage:     "track test [track-number]",
		AllowArgs: true,
		Run: func(c *grumble.Context) error {
			trackNumber, err := extractTrackNumberArgument(c)
			if err != nil {
				return err
			}
			destination := filepath.Join(os.TempDir(), "test-track.wav")
			err = loader.Load(trackNumber, destination, library)
			if err != nil {
				return err
			}
			f, err := os.Open(destination)
			if err != nil {
				return err
			}
			streamer, format, err := wav.Decode(f)
			if err != nil {
				return err
			}
			defer streamer.Close()
			_ = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
			speaker.Play(streamer)
			app.Println("Press return to exit…")
			reader := bufio.NewReader(os.Stdin)
			_, err = reader.ReadString('\n')
			speaker.Clear()
			return nil
		},
		Completer: nil,
	})

	app.AddCommand(&trackCommand)
}
