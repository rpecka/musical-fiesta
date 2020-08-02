package audio

import (
	"fiesta/src/defaults"
	"fmt"
	"os/exec"
	"strconv"
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
	ApplyTransformations(inputPath string, outputPath string, start *float64, end *float64) error
}

type ffmpegAudioManipulator struct {
	ffmpegPath string
}

func InitializeAudioManipulator() (Manipulator, error) {
	findCommand := defaults.WhichCommand
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

func makeBaseArgs(inputPath string) []string {
	return []string{
		"-i", inputPath,
		mapMetadata, "-1",
		bitexact,
		"-y",
	}
}

func (f *ffmpegAudioManipulator) ConvertToWav(inputPath string, outputPath string) error {
	args := makeBaseArgs(inputPath)
	args = append(args,
		audioQuality, "100",
		sampleRate, "22050",
		audioChannels, "1",
		audioCodec, pcmS16le,
		outputPath,
	)
	output, err := exec.Command(f.ffmpegPath, args...).CombinedOutput()
	outputString := string(output)
	if err != nil {
		return fmt.Errorf("ffmpeg execution error: %v - %v", err, outputString)
	}
	return nil
}

func (f *ffmpegAudioManipulator) ApplyTransformations(inputPath string, outputPath string, start *float64, end *float64) error {
	args := makeBaseArgs(inputPath)
	if start != nil {
		args = append(args, "-ss", strconv.FormatFloat(*start, 'f', 1, 64))
	}
	if end != nil {
		args = append(args, "-to", strconv.FormatFloat(*end, 'f', 1, 64))
	}
	args = append(args, outputPath)
	output, err := exec.Command(f.ffmpegPath, args...).CombinedOutput()
	outputString := string(output)
	if err != nil {
		return fmt.Errorf("ffmpeg execution error: %v - %v", err, outputString)
	}
	return nil
}
