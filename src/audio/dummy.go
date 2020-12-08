// +build dummy

package audio

import "errors"

type dummyManipulator struct {

}

func InitializeManipulator() (Manipulator, error) {
	return dummyManipulator{}, nil
}

func (d dummyManipulator) ConvertToWav(inputPath string, outputPath string) error {
	return errors.New("this is a dummy")
}

func (d dummyManipulator) ApplyTransformations(inputPath string, outputPath string, start *float64, end *float64) error {
	return errors.New("this is a dummy")
}



