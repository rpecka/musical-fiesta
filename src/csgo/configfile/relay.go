package configfile

import (
	"errors"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

const (
	RelayFileName = "fiesta-relay.cfg"
)

const (
	LoadNumberResult = iota
	TagResult
)

type ParseResult struct {
	ResultType  int
	TrackNumber int
	Tag         string
	Offset      *int
}

func makeLoadNumberResult(trackNumber int) *ParseResult {
	return &ParseResult{
		ResultType:  LoadNumberResult,
		TrackNumber: trackNumber,
	}
}

func makeTagResult(tag string) *ParseResult {
	return &ParseResult{
		ResultType: TagResult,
		Tag:        tag,
	}
}

func Parse(relayFilePath string, trackRelayKey, offsetRelayKey string) (*ParseResult, error) {
	f, err := os.Open(relayFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open relay file at: %v because: %v", relayFilePath, err)
	}
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read relay file at: %v because: %v", relayFilePath, err)
	}
	utfBytes, err := charmap.ISO8859_1.NewDecoder().Bytes(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode iso8859_1 bytes to utf-8 because: %v", err)
	}
	utfString := string(utfBytes)
	regex := bindCommandRegex(trackRelayKey)
	matches := regex.FindStringSubmatch(utfString)
	if matches == nil || len(matches) < 2 {
		return nil, errors.New("could not find matching bind command in relay file")
	}
	command := matches[1]
	integer, err := strconv.Atoi(command)
	var result *ParseResult
	if err != nil {
		result = makeTagResult(command)
	} else {
		result = makeLoadNumberResult(integer)
	}

	offsetRegex := bindCommandRegex(offsetRelayKey)
	offsetMatches := offsetRegex.FindStringSubmatch(utfString)
	if offsetMatches == nil || len(offsetMatches) < 2 {
		return nil, errors.New("could not find matching bind command in relay file")
	}
	offsetString := offsetMatches[1]
	integerOffset, err := strconv.Atoi(offsetString)
	if err != nil {
		return result, nil
	} else {
		result.Offset = &integerOffset
		return result, nil
	}
}

func bindCommandRegex(relayKey string) *regexp.Regexp {
	return regexp.MustCompile("bind \"" + relayKey + "\" \"(.+)\"")
}

func offsetBindCommandRegex(relayKey string) *regexp.Regexp {
	return regexp.MustCompile("bind \"" + relayKey + "\" \"(.+)%\"")
}
