package configfile

import (
	"fmt"
	"testing"
)

const (
	testTrackKey  = "="
	testOffsetKey = "ins"
)

func TestParse(t *testing.T) {
	ten := 10
	cases := map[string]*ParseResult{
		fmt.Sprintf("bind \"%s\" \"10\"", testTrackKey): &ParseResult{
			TrackNumber: 10,
		},
		fmt.Sprintf("bind \"%s\" \"10\"\nbind \"%s\" \"10%%\"", testTrackKey, testOffsetKey): &ParseResult{
			TrackNumber: 10,
			Offset:      &ten,
		},
		fmt.Sprintf("bind \"%s\" \"10%%\"", testOffsetKey): nil,
		"":                 nil,
		"bind \"3\" \"4\"": nil,
	}

	for relayString, expectedResult := range cases {
		result, err := parse(relayString, testTrackKey, testOffsetKey)
		if err != nil {
			t.Error(err)
		}
		if expectedResult == nil {
			if result != nil {
				t.Error("expected result to be nil but it was not")
			}
		} else if result == nil {
			t.Error("expected result was not nil but the result was")
		} else {
			if expectedResult.TrackNumber != result.TrackNumber {
				t.Errorf("expected track number %v but got %v", expectedResult.TrackNumber, result.TrackNumber)
			} else if expectedResult.Offset != nil {
				if result.Offset == nil {
					t.Errorf("expected the offset to be %v but it was nil", *(expectedResult.Offset))
				} else if *(expectedResult.Offset) != *(result.Offset) {
					t.Errorf("expected offset %v but got %v", *(expectedResult.Offset), *(result.Offset))
				}
			}
		}
	}
}
