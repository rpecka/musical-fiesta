package audio

// #cgo LDFLAGS: -L${SRCDIR}/../ffwrapper -lffwrapper -lavformat -lavcodec -lavutil -lswresample
// #include <stdlib.h>
// #include "../ffwrapper/ffwrapper.h"
import "C"
import (
	"fmt"
	"unsafe"
)

type wrapperFFMPEG struct {

}

func InitializeWrapperFFMPEGManipulator() (Manipulator, error) {
	audio := wrapperFFMPEG{

	}
	return &audio, nil
}

func (w *wrapperFFMPEG) ConvertToWav(inputPath string, outputPath string) error {
	cInputPath := C.CString(inputPath)
	defer C.free(unsafe.Pointer(cInputPath))
	cOutputPath := C.CString(outputPath)
	defer C.free(unsafe.Pointer(cOutputPath))
	ret := C.transcode(cInputPath, cOutputPath)
	if ret < 0 {
		return fmt.Errorf("transcoding procedure exited with code: '%v'", ret)
	}
	return nil
}

func (w *wrapperFFMPEG) ApplyTransformations(inputPath string, outputPath string, start *float64, end *float64) error {
	panic("implement me")
}
