package configfile

import "testing"

func TestGenerateTagGroupsSingles(t *testing.T) {
	singles := []EnumeratedTrack{
		{
			Number: 1,
			Name:   "Hi",
			Tags:   []string{"hi", "bye"},
		},
		{
			Number: 2,
			Name: "Sorry",
			Tags: []string{"nick", "of", "time"},
		},
	}
	resultSingles, resultGroups := generateTagGroups(singles)
	if len(resultGroups) > 0 {
		t.Error("there should not have been any groups")
	}
	if len(resultSingles) != 5 {
		t.Error("there should be five singles in the results")
	}
}

func TestGenerateTagGroupsGroups(t *testing.T) {
	tracks := []EnumeratedTrack{
		{
			Number: 1,
			Name: "Hi",
			Tags: []string{"hi", "bye"},
		},
		{
			Number: 2,
			Name: "Bye",
			Tags: []string{"hi", "bye"},
		},
	}
	resultSingles, resultGroups := generateTagGroups(tracks)
	if len(resultGroups) != 2 {
		t.Error("there should not have been any groups")
	}
	if len(resultSingles) > 0 {
		t.Error("there should be five singles in the results")
	}
}
