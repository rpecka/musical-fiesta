package audio

import (
	"fiesta/src/crossplatform"
	"fmt"
	"os/exec"
	"strings"
)

const (
	ffmpegExecutableName = "ffmpeg"

	mapMetadata   = "-map_metadata"
	bitexact      = "-bitexact"
	audioQuality  = "-aq"
	sampleRate    = "-ar"
	audioChannels = "-ac"
	audioCodec    = "-acodec"
	pcmS16le      = "pcm_s16le"
)

type Manipulator interface {
	ConvertToWav(inputPath string, outputPath string) error
}

type ffmpegAudioManipulator struct {
	ffmpegPath string
}

func InitializeAudioManipulator() (Manipulator, error) {
	findCommand := crossplatform.WhichCommand
	output, err := exec.Command(findCommand, ffmpegExecutableName).Output()
	if err != nil {
		return nil, err
	}
	ffmpegPath := strings.TrimSpace(string(output))
	audio := ffmpegAudioManipulator{
		ffmpegPath: ffmpegPath,
	}
	return &audio, nil
}

func (f *ffmpegAudioManipulator) ConvertToWav(inputPath string, outputPath string) error {
	output, err := exec.Command(f.ffmpegPath,
		"-i", inputPath,
		mapMetadata, "-1",
		bitexact,
		audioQuality, "100",
		sampleRate, "22050",
		audioChannels, "1",
		audioCodec, pcmS16le,
		outputPath).CombinedOutput()
	outputString := string(output)
	if err != nil {
		return fmt.Errorf("ffmpeg execution error: %v - %v", err, outputString)
	}
	return nil
}
