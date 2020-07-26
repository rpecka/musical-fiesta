package configfile

import (
	"fmt"
	"strings"
)

func formatTagList(tags []string) string {
	return fmt.Sprintf("[%s]", strings.Join(tags, ", "))
}

func writeTrackList(dest string, tracks []EnumeratedTrack) error {
	writer, err := newWriter(dest)
	if err != nil {
		return err
	}
	defer writer.close()

	err = writer.writeEchoHeader("Tracks")
	if err != nil {
		return err
	}

	for _, track := range tracks {
		_ = writer.writeEcho(fmt.Sprintf("%d. %s %v", track.Number, track.Name, formatTagList(track.Tags)))
	}

	return nil
}
