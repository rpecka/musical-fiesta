package library

import "time"

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

func (t track) NeedsModification() bool {
	return t.Trim != nil
}
