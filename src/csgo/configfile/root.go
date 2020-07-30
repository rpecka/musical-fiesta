package configfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const (
	// CFG Hierarchy
	rootCFGName      = "fiesta.cfg"
	cfgDirName       = "fiesta"
	trackListCFGName = "fiesta_tracklist.cfg"

	// Listing Tracks
	listTracksCommand     = "fiesta_list"
	listTracksAliasList   = "list"
	listTracksAliasTracks = "tracks"
	listTracksAliasLA     = "la"

	// Audio Controls
	toggleCommand = "fiesta_toggle"
	playCommand   = "fiesta_play"
	stopCommand   = "fiesta_stop"

	// CFG Updates
	updateCommand = "fiesta_updatecfg"
)

type EnumeratedTrack struct {
	Number int
	Name   string
	Tags   []string
}

func rootCFGPath(rootDir string) string {
	return filepath.Join(rootDir, rootCFGName)
}

func cfgDirPath(rootDir string) string {
	return filepath.Join(rootDir, cfgDirName)
}

func makeLoadLogic(relayKey string, track EnumeratedTrack) string {
	return chainCommands([]string{
		makeBindCommand(relayKey, strconv.Itoa(track.Number)),
		updateCommand,
		makeEcho(fmt.Sprintf("Loaded %s", track.Name)),
	})
}

func generateTagGroups(tracks []EnumeratedTrack) (map[string]*EnumeratedTrack, map[string][]*EnumeratedTrack) {
	singles := make(map[string]*EnumeratedTrack)
	groups := make(map[string][]*EnumeratedTrack)
	for _, track := range tracks {
		for _, tag := range track.Tags {
			existing, inSingles := singles[tag]
			if inSingles {
				groups[tag] = []*EnumeratedTrack{existing, &track}
				delete(singles, tag)
			} else {
				group, inGroups := groups[tag]
				if inGroups {
					groups[tag] = append(group, &track)
				} else {
					singles[tag] = &track
				}
			}
		}
	}
	return singles, groups
}

func WriteConfigFiles(rootDir string, playKey string, relayKey string, enumeratedTracks []EnumeratedTrack) error {
	writer, err := newWriter(rootCFGPath(rootDir))
	if err != nil {
		return err
	}
	defer writer.close()

	cfgDirPath := cfgDirPath(rootDir)
	err = os.RemoveAll(cfgDirPath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(cfgDirPath, 0755)
	if err != nil {
		return err
	}

	// Listing Tracks
	err = writeTrackList(filepath.Join(rootDir, cfgDirName, trackListCFGName), enumeratedTracks)
	if err != nil {
		return err
	}
	err = writer.writeHeader("Listing Tracks")
	if err != nil {
		return err
	}
	// Use / no matter what here because we're talking in .cfg
	trackListRelativePath := cfgDirName + "/" + trackListCFGName
	err = writer.writeAlias(listTracksCommand, makeExecCommand(trackListRelativePath))
	if err != nil {
		return err
	}

	for _, alias := range []string{listTracksAliasList, listTracksAliasTracks, listTracksAliasLA} {
		err = writer.writeAlias(alias, listTracksCommand)
		if err != nil {
			return err
		}
	}

	// Audio Controls
	err = writer.writeHeader("Audio Controls")
	if err != nil {
		return err
	}
	_ = writer.writeAlias(toggleCommand, playCommand)
	_ = writer.writeAlias(playCommand, chainCommands([]string{
		makeAliasCommand(toggleCommand, stopCommand), // Set the toggle to stop
		makeVoiceInputFromFile(true),                 // Start pointing voice input to a file
		makeVoiceLoopBack(true),                      // Loop voice back to the user so they can hear too
		makeVoiceRecord(true),                        // Start voice input to the game
	}))
	_ = writer.writeAlias(stopCommand, chainCommands([]string{
		makeVoiceRecord(false),                       // Stop voice input
		makeVoiceInputFromFile(false),                // Stop redirecting file to voice output
		makeVoiceLoopBack(false),                     // Stop playing the user's voice output back to them
		makeAliasCommand(toggleCommand, playCommand), // Set the toggle start again
	}))
	_ = writer.writeBind(playKey, toggleCommand)

	// CFG Updates
	_ = writer.writeHeader("CFG Updates")
	_ = writer.writeAlias(updateCommand, makeHostWriteconfig(RelayFileName))

	// Loading Tracks by Index
	_ = writer.writeHeader("Loading Tracks by Index")
	for _, track := range enumeratedTracks {
		_ = writer.writeAlias(strconv.Itoa(track.Number), makeLoadLogic(relayKey, track))
	}

	// Loading Tracks by Tag
	singles, _ := generateTagGroups(enumeratedTracks)
	_ = writer.writeHeader("Loading Tracks by Tag")

	for tag, track := range singles {
		_ = writer.writeAlias(tag, strconv.Itoa(track.Number))
	}

	return nil
}

func DeleteConfigFiles(rootDir string) {
	_ = os.Remove(rootCFGPath(rootDir))
	_ = os.RemoveAll(cfgDirPath(rootDir))
}
