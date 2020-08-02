package audio



type Manipulator interface {
	ConvertToWav(inputPath string, outputPath string) error
	ApplyTransformations(inputPath string, outputPath string, start *float64, end *float64) error
}
