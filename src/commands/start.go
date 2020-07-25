package commands

import (
	"bufio"
	"fiesta/src/csgo"
	"fiesta/src/csgo/configfile"
	"fiesta/src/library"
	"fiesta/src/loader"
	"fiesta/src/settings"
	"fmt"
	"github.com/desertbit/grumble"
	"os"
	"path/filepath"
)

func enumerateTracks(l library.Library) ([]configfile.EnumeratedTrack, error) {
	tracks, err := l.Tracks()
	if err != nil {
		return nil, err
	}
	enumeratedTracks := make([]configfile.EnumeratedTrack, 0, len(tracks))
	for index, track := range tracks {
		enumeratedTracks = append(enumeratedTracks,
			configfile.EnumeratedTrack{Number: library.TrackIndexToNumber(index), Name: track.Name, Tags: track.Tags})
	}
	return enumeratedTracks, nil
}

func addStart(app *grumble.App, settings *settings.Settings, library *library.Library) {
	app.AddCommand(&grumble.Command{
		Name:      "start",
		Help:      "start listening for commands from CSGO",
		Usage:     "start",
		AllowArgs: false,
		Run: func(c *grumble.Context) error {
			stop := make(chan bool)
			userdataDir, err := (*settings).UserdataDirPath()
			if err != nil {
				return err
			}
			csgoDir, err := (*settings).CSGODirPath()
			if err != nil {
				return err
			}
			enumeratedTracks, err := enumerateTracks(*library)
			if err != nil {
				return err
			}
			cfgPath := csgo.PathToCFG(csgoDir)
			err = configfile.WriteConfigFiles(cfgPath, "z", "=", enumeratedTracks)
			if err != nil {
				return err
			}
			destination := filepath.Join(csgoDir, csgo.VoiceInputFileName)
			err = loader.Start(userdataDir, "z", stop, destination, library)
			if err != nil {
				return err
			}
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Press return to exit…")
			_, err = reader.ReadString('\n')
			configfile.DeleteConfigFiles(cfgPath)
			stop <- true
			return err
		},
	})
}
