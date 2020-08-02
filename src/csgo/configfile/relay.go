package configfile

import (
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

type ParseResult struct {
	TrackNumber int
	Offset      *int
}

func bindCommandRegex(relayKey string) *regexp.Regexp {
	return regexp.MustCompile("bind \"" + relayKey + "\" \"(.+)\"")
}

func offsetBindCommandRegex(relayKey string) *regexp.Regexp {
	return regexp.MustCompile("bind \"" + relayKey + "\" \"(.+)%\"")
}

func parse(relayFileString string, trackRelayKey, offsetRelayKey string) (*ParseResult, error) {
	regex := bindCommandRegex(trackRelayKey)
	matches := regex.FindStringSubmatch(relayFileString)
	if matches == nil || len(matches) < 2 {
		return nil, nil
	}
	command := matches[1]
	integer, err := strconv.Atoi(command)
	var result *ParseResult
	if err != nil {
		return nil, fmt.Errorf("could not parse track to load because: %v", err)
	} else {
		result = &ParseResult{
			TrackNumber: integer,
			Offset:      nil,
		}
	}

	offsetRegex := offsetBindCommandRegex(offsetRelayKey)
	offsetMatches := offsetRegex.FindStringSubmatch(relayFileString)
	if offsetMatches == nil || len(offsetMatches) < 2 {
		return result, nil
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
	return parse(utfString, trackRelayKey, offsetRelayKey)
}
