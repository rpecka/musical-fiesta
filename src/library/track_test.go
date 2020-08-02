package library

import (
	"testing"
	"time"
)

func TestResolveStartEnd(t *testing.T) {
	test := track{Trim: nil}
	start, end, err := test.resolveStartEnd(nil)
	if err != nil {
		t.Error(err)
	}
	if start != nil || end != nil {
		t.Error("a track with a nil trim and no offset should result in no override start and end time")
	}

	thirtySeconds, err := time.ParseDuration("30s")
	if err != nil {
		t.Error(err)
	}
	test = track{Trim: &trackTrim{End: &thirtySeconds}}
	start, end, err = test.resolveStartEnd(nil)
	if err != nil {
		t.Error(err)
	}
	if start != nil || end == nil || *end != 30 {
		t.Error("track trimmed to 30s end should have 30s end when there is no offset applied")
	}

	test = track{Trim: &trackTrim{End: &thirtySeconds}}
	offset := 50
	start, end, err = test.resolveStartEnd(&offset)
	if err != nil {
		t.Error(err)
	}
	if start == nil || *start != 15 || end == nil || *end != 30 {
		t.Errorf("track trimmed to 30s end should have 15s when a 50%% offset is applied")
	}
}
