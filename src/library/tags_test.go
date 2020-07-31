package library

import "testing"

func TestIsValidTag(t *testing.T) {
	invalid := []string{
		"drop",
		"DROP",
		"reload",
		"1",
	}

	valid := []string{
		"hi",
		"yoda",
		"foo",
	}

	for _, tag := range invalid {
		if isValidTag(tag) {
			t.Errorf("%s is not a valid tag", tag)
		}
	}
	for _, tag := range valid {
		if !isValidTag(tag) {
			t.Errorf("%s should be a valid tag", tag)
		}
	}
}

func TestGenerateTagsFromFilename(t *testing.T) {
	files := map[string]map[string]bool {
		"DROP IT [Dubstep] Sporty O HitDrop-Mix": map[string]bool{
			"it": true, "[dubstep]": true, "sporty": true, "o": true, "hitdrop-mix": true,
		},
	}
	for filename, expectedTags := range files {
		gen := generateTagsFromFilename(filename)
		for _, tag := range gen {
			_, ok := expectedTags[tag]
			if !ok {
				t.Errorf("unexpected tag from %s: %s", filename, tag)
			} else {
				delete(expectedTags, tag)
			}
		}
		if len(expectedTags) > 0 {
			t.Errorf("some tags missed: %v", expectedTags)
		}
	}
}
