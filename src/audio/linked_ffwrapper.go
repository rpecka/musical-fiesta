package audio

// Path to liblzma must be provided by environment variable i.e. CGO_LDFLAGS=/usr/local/Cellar/xz/5.2.5/lib/liblzma.a

// #cgo LDFLAGS: -L${SRCDIR}/../ffwrapper -L${SRCDIR}/../../../ffmpeg/lib -lffwrapper
// #cgo LDFLAGS: -lavformat -lavcodec -lavutil -lswresample
// #cgo LDFLAGS: -lbz2 -lz -liconv
// #cgo LDFLAGS: -framework CoreVideo -framework CoreMedia -framework CoreFoundation
// #cgo LDFLAGS: -framework VideoToolbox -framework AudioToolbox
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
