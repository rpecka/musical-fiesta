package library

import (
	"github.com/faiface/beep/wav"
	"os"
	"time"
)

type track struct {
	Name string     `json:"name"`
	Path string     `json:"path"`
	Tags []string   `json:"tags"`
	Trim *trackTrim `json:"trim,omitempty"`
}

type trackTrim struct {
	Start *time.Duration `json:"start"`
	End   *time.Duration `json:"end"`
}

func (t track) resolveStartEnd(offset *int) (start, end *float64, err error) {
	if t.Trim != nil {
		if t.Trim.Start != nil {
			r := t.Trim.Start.Seconds()
			start = &r
		}
		if t.Trim.End != nil {
			r := t.Trim.End.Seconds()
			end = &r
		}
	}
	if offset != nil {
		percent := float64(*offset) / 100.0
		if end != nil {
			if start != nil {
				*start = (*end-*start)*percent + *start
			} else {
				offsetStart := *end * percent
				start = &offsetStart
			}
		} else {
			f, err := os.Open(t.Path)
			if err != nil {
				return nil, nil, err
			}
			defer f.Close()
			streamer, format, err := wav.Decode(f)
			if err != nil {
				return nil, nil, err
			}
			duration := format.SampleRate.D(streamer.Len())
			if start != nil {
				*start = (duration.Seconds()-*start)*percent + *start
			} else {
				offsetStart := duration.Seconds() * percent
				start = &offsetStart
			}
		}
	}
	return
}
