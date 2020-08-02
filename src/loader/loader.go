package loader

import (
	"errors"
	"fiesta/src/configwatcher"
	"fiesta/src/csgo/configfile"
	"fiesta/src/library"
	"fiesta/src/util"
	"fmt"
	"github.com/faiface/beep/wav"
	"os"
)

func Load(trackNumber int, offset *int, destination string, library library.Library) error {
	track, err := library.GetTrack(trackNumber)
	if err != nil {
		return err
	}

	var start *float64
	var end *float64
	if track.Trim != nil {
		if track.Trim.Start != nil {
			r := track.Trim.Start.Seconds()
			start = &r
		}
		if track.Trim.End != nil {
			r := track.Trim.End.Seconds()
			end = &r
		}
	}
	if offset != nil {
		percent := float64(*offset) / 100.0
		if end != nil {
			if start != nil {
				*start = (*end - *start) * percent + *start
			} else {
				offsetStart := *end * percent
				start = &offsetStart
			}
		} else {
			f, err := os.Open(track.Path)
			if err != nil {
				return err
			}
			defer f.Close()
			streamer, format, err := wav.Decode(f)
			if err != nil {
				return err
			}
			duration := format.SampleRate.D(streamer.Len())
			if start != nil {
				*start = (duration.Seconds() - *start) * percent + *start
			} else {
				offsetStart := duration.Seconds() * percent
				start = &offsetStart
			}
		}
	}

	if start == nil && end == nil {
		return util.CopyFile(track.Path, destination)
	} else {
		return library.Manipulator().ApplyTransformations(track.Path, destination, start, end)
	}
}

func Start(userdataDir string, trackRelayKey, offsetRelayKey string, stop chan bool, destination string, library library.Library) error {
	started := make(chan error)
	fileChanged := make(chan string)
	go func() {
		configwatcher.Start(userdataDir, configfile.RelayFileName, started, fileChanged, stop)
	}()
	err := <-started
	if err != nil {
		return fmt.Errorf("failed to start config watcher: %v", err)
	}
	go func() {
		for path := range fileChanged {
			result, err := configfile.Parse(path, trackRelayKey, offsetRelayKey)
			if err != nil {
				fmt.Print(err)
				continue
			}
			switch result.ResultType {
			case configfile.LoadNumberResult:
				err = Load(result.TrackNumber, result.Offset, destination, library)
			case configfile.TagResult:
				err = errors.New("string commands are not yet supported")
			}
			if err != nil {
				fmt.Printf("failed to load track because: %v\n", err)
			}
		}
	}()
	return nil
}
