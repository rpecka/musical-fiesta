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

func Parse(relayFilePath string, relayKey string) (*ParseResult, error) {
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
	regex := bindCommandRegex(relayKey)
	matches := regex.FindStringSubmatch(utfString)
	if matches == nil || len(matches) < 2 {
		return nil, errors.New("could not find matching bind command in relay file")
	}
	command := matches[1]
	integer, err := strconv.Atoi(command)
	if err != nil {
		return makeTagResult(command), nil
	} else {
		return makeLoadNumberResult(integer), nil
	}
}

func bindCommandRegex(relayKey string) *regexp.Regexp {
	return regexp.MustCompile("bind \"" + relayKey + "\" \"(.+)\"")
}
