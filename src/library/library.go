package library

import (
	"encoding/json"
	"errors"
	"fiesta/src/audio"
	"fiesta/src/csgo/commands"
	"fiesta/src/util"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	libraryFileName      = "library.json"
	libraryTracksDirName = "tracks"
)

type realLibrary struct {
	path        string
	manipulator audio.Manipulator
}

type libraryFile struct {
	Tracks []track `json:"tracks"`
}

type Library interface {
	Tracks() ([]track, error)
	GetTrack(trackNumber int) (*track, error)
	Import(trackPath string) error
	ImportDir(trackDirPath string) (failures []string, err error)
	DeleteTrack(trackNumber int) error
}

func InitializeLibrary(libraryDir string, manipulator audio.Manipulator) (Library, error) {
	if !util.Exists(libraryDir) {
		return nil, errors.New("Library directory: " + libraryDir + " does not exist")
	}

	lib := realLibrary{
		path:        libraryDir,
		manipulator: manipulator,
	}

	libraryFilePath := lib.libraryFilePath()
	if !util.Exists(libraryFilePath) {
		f, err := os.Create(libraryFilePath)
		if err != nil {
			return nil, err
		}
		byteValue, err := json.Marshal(libraryFile{Tracks: []track{}})
		if err != nil {
			return nil, err
		}
		_, err = f.Write(byteValue)
		if err != nil {
			return nil, err
		}
		_ = f.Close()

	}

	tracksDirPath := lib.tracksDirPath()
	if !util.Exists(tracksDirPath) {
		err := os.MkdirAll(tracksDirPath, 0755)
		if err != nil {
			return nil, err
		}
	}

	return &lib, nil
}

func (l realLibrary) libraryFilePath() string {
	return filepath.Join(l.path, libraryFileName)
}

func (l realLibrary) tracksDirPath() string {
	return filepath.Join(l.path, libraryTracksDirName)
}

func (l realLibrary) readLibraryFile() (*libraryFile, error) {
	f, err := os.Open(l.libraryFilePath())
	if err != nil {
		return nil, err
	}
	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var libFile libraryFile
	err = json.Unmarshal(byteValue, &libFile)
	if err != nil {
		return nil, err
	}

	_ = f.Close()
	return &libFile, nil
}

func (l realLibrary) writeLibraryFile(libFile libraryFile) error {
	data, err := json.MarshalIndent(libFile, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(l.libraryFilePath(), data, 0644)
}

func (l realLibrary) insertTrack(t track) error {
	libFile, err := l.readLibraryFile()
	if err != nil {
		return err
	}

	// There is probably a better way to do this
	libFile.Tracks = append(libFile.Tracks, t)
	sort.Slice(libFile.Tracks, func(i, j int) bool {
		return libFile.Tracks[i].Name < libFile.Tracks[j].Name
	})

	err = l.writeLibraryFile(*libFile)
	return err
}

func (l realLibrary) generateTagsFromFilename(trackFilename string) []string {
	trimmed := strings.TrimSpace(trackFilename)
	words := strings.Split(trimmed, " ")
	for index, word := range words {
		words[index] = strings.ToLower(word)
	}
	uniqueWords := util.Unique(words)
	return uniqueWords
}

func trackNumberToIndex(trackNumber int) int {
	return trackNumber - 1
}

func validateTrackNumber(trackNumber int, libFile libraryFile) error {
	if trackNumber <= 0 {
		return errors.New("track number must be greater than zero")
	}
	if trackNumber > len(libFile.Tracks) {
		return errors.New("track number is out of bounds: " + string(len(libFile.Tracks)))
	}
	return nil
}

func (l *realLibrary) Tracks() ([]track, error) {
	libFile, err := l.readLibraryFile()
	if err != nil {
		return nil, err
	}
	return libFile.Tracks, nil
}

func (l *realLibrary) GetTrack(trackNumber int) (*track, error) {
	libFile, err := l.readLibraryFile()
	if err != nil {
		return nil, err
	}
	err = validateTrackNumber(trackNumber, *libFile)
	if err != nil {
		return nil, err
	}
	trackIndex := trackNumberToIndex(trackNumber)
	return &libFile.Tracks[trackIndex], nil
}

func (l *realLibrary) Import(trackPath string) error {
	if !util.Exists(trackPath) {
		return errors.New("the path: `" + trackPath + "` does not exist")
	}
	inputFileName := util.RemoveFileExtension(filepath.Base(trackPath))
	outputFileName := inputFileName + ".wav"
	outputFilePath := filepath.Join(l.tracksDirPath(), outputFileName)
	if util.Exists(outputFilePath) {
		return errors.New("Could not import track to destination: " + outputFilePath + " because a file already exists at that location")
	}

	err := l.manipulator.ConvertToWav(trackPath, outputFilePath)
	if err != nil {
		return err
	}

	tags := l.generateTagsFromFilename(inputFileName)
	legalTags := util.Filter(tags, func(s string) bool {
		return !commands.IsIllegal(s)
	})

	track := track{
		Name: inputFileName,
		Path: outputFilePath,
		Tags: legalTags,
	}

	err = l.insertTrack(track)
	if err != nil {
		return err
	}
	return nil
}

func (l *realLibrary) ImportDir(trackDirPath string) (failures []string, err error) {
	if !util.Exists(l.tracksDirPath()) {
		return nil, errors.New("the path: `" + trackDirPath + "` does not exist")
	}
	contents, err := ioutil.ReadDir(trackDirPath)
	if err != nil {
		return nil, fmt.Errorf("could not import directory of tracks because: %v", err)
	}
	failures = make([]string, 0)
	for _, file := range contents {
		if file.IsDir() {
			continue
		}
		trackPath := filepath.Join(trackDirPath, file.Name())
		err = l.Import(trackPath)
		if err != nil {
			failures = append(failures, trackPath)
		}
	}
	return failures, nil
}

func (l *realLibrary) DeleteTrack(trackNumber int) error {
	libFile, err := l.readLibraryFile()
	if err != nil {
		return err
	}
	err = validateTrackNumber(trackNumber, *libFile)
	if err != nil {
		return err
	}
	trackIndex := trackNumberToIndex(trackNumber)
	trackPath := libFile.Tracks[trackIndex].Path

	libFile.Tracks = append(libFile.Tracks[:trackIndex], libFile.Tracks[trackIndex+1:]...)
	err = l.writeLibraryFile(*libFile)

	err = os.Remove(trackPath)
	if err != nil {
		// TODO: Need some kind of logging for this but it seems non-fatal to me
	}
	return err
}