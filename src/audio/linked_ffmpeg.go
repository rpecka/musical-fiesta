package audio

/*
	#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"github.com/giorgisio/goav/avcodec"
	"github.com/giorgisio/goav/avformat"
	"github.com/giorgisio/goav/avutil"
	"github.com/giorgisio/goav/swresample"
	"unsafe"
)

type linkedFFMPEG struct {

}

type transcoder struct {
	inFmtCtx *avformat.Context
	inCodecCtx *avcodec.Context
	outFmtCtx *avformat.Context
	outCodecCtx *avcodec.Context
	resampleCtx *swresample.Context
	fifo *avutil.AVAudioFifo
	pts int
}

func InitializeLinkedFFMPEG() Manipulator {
	return &linkedFFMPEG{}
}

func (f *linkedFFMPEG) ConvertToWav(inputPath string, outputPath string) error {
	inFmtCtx, inCodecCtx, err := openInput(inputPath)
	if err != nil {
		return err
	}
	defer inFmtCtx.AvformatCloseInput()
	defer inCodecCtx.AvcodecFreeContext()

	outFmtCtx, outCodecCtx, err := openOutput(outputPath)
	if err != nil {
		return err
	}
	defer outFmtCtx.Pb().Close()
	defer outFmtCtx.AvformatFreeContext()
	defer outCodecCtx.AvcodecFreeContext()

	resampleCtx, err := initResampler(inCodecCtx, outCodecCtx)
	if err != nil {
		return err
	}
	defer resampleCtx.SwrFree()

	fifo, err := initFifo(outCodecCtx)
	if err != nil {
		return err
	}
	defer fifo.AVAudioFifoFree()

	t := transcoder{
		inFmtCtx:     inFmtCtx,
		inCodecCtx:   inCodecCtx,
		outFmtCtx:    outFmtCtx,
		outCodecCtx:  outCodecCtx,
		resampleCtx:  resampleCtx,
		fifo:         fifo,
		pts:          0,
	}

	return t.transcode()
}

func (f *linkedFFMPEG) ApplyTransformations(inputPath string, outputPath string, start *float64, end *float64) error {
	//avformat_seek_file seek to timestamp
	return nil
}

func (t transcoder) transcode() error {
	err := writeOutputFileHeader(t.outFmtCtx)
	if err != nil {
		return err
	}

	for {
		outFrameSize := t.outCodecCtx.FrameSize()
		finished := false

		for t.fifo.Size() < outFrameSize {
			finished, err = t.readDecodeConvertAndStore()
			if err != nil {
				return fmt.Errorf("could not convert data: %v", err)
			}

			if finished {
				break
			}
		}

		for t.fifo.Size() >= outFrameSize || (finished && t.fifo.Size() > 0) {
			err = t.loadEncodeAndWrite()
			if err != nil {
				return err
			}
		}

		if finished {
			dataWritten := true
			for dataWritten {
				dataWritten, err = t.encodeAudioFrame(nil)
				if err != nil {
					return err
				}
			}
			break
		}
	}
	return nil
}

func openInput(path string) (inFmtCtx *avformat.Context, inCodecCtx *avcodec.Context, err error) {
	ret := avformat.AvformatOpenInput(&inFmtCtx, path, nil, nil)
	if ret < 0 {
		err = fmt.Errorf("could not open input file '%v' (error '%v')", path, avutil.ErrorFromCode(ret))
		return
	}

	*inFmtCtx.Flags() |= avformat.AVFMT_FLAG_BITEXACT

	ret = avformat.AvformatFindStreamInfo(inFmtCtx, nil)
	if ret < 0 {
		inFmtCtx.AvformatCloseInput()
		err = fmt.Errorf("could not find stream info (error '%v')", avutil.ErrorFromCode(ret))
		return
	}

	if inFmtCtx.NbStreams() != 1 {
		inFmtCtx.AvformatCloseInput()
		err = fmt.Errorf("expected one input stream but found %v", inFmtCtx.NbStreams())
		return
	}

	inCodec := avcodec.AvcodecFindDecoder(*inFmtCtx.Streams()[0].CodecParameters().AvCodecGetId())
	if inCodec == nil {
		inFmtCtx.AvformatCloseInput()
		err = fmt.Errorf("could not find input codec")
		return
	}

	inCodecCtx = inCodec.AvcodecAllocContext3()
	if inCodecCtx == nil {
		inFmtCtx.AvformatCloseInput()
		err = fmt.Errorf("could not allocate a decoding context")
		return
	}

	*inCodecCtx.Flags() |= avformat.AVFMT_FLAG_BITEXACT

	ret = avcodec.AvcodecParametersToContext(inCodecCtx, inFmtCtx.Streams()[0].CodecParameters())
	if ret < 0 {
		inFmtCtx.AvformatCloseInput()
		inCodecCtx.AvcodecFreeContext()
		err = fmt.Errorf("could not initialize stream parameters with demuxer information (error '%v')", avutil.ErrorFromCode(ret))
		return
	}

	ret = inCodecCtx.AvcodecOpen2(inCodec, nil)
	if ret < 0 {
		inFmtCtx.AvformatCloseInput()
		inCodecCtx.AvcodecFreeContext()
		err = fmt.Errorf("failed to open the decoder (error '%v')", avutil.ErrorFromCode(ret))
		return
	}

	return
}

func openOutput(path string) (outFmtCtx *avformat.Context, outCodecCtx *avcodec.Context, err error) {
	outIOCtx, err := avformat.AvIOOpen(path, avformat.AVIO_FLAG_WRITE)
	if err != nil {
		return
	}

	outFmtCtx = avformat.AvformatAllocContext()
	if outFmtCtx == nil {
		err = fmt.Errorf("could not allocate output format context")
		outIOCtx.Close()
		return
	}

	*outFmtCtx.Flags() |= avformat.AVFMT_FLAG_BITEXACT
	outFmtCtx.SetPb(outIOCtx)

	outFmt := avformat.AvGuessFormat(nil, &path, nil)
	if outFmt == nil {
		err = errors.New("could not find output file format")
		_ = outIOCtx.Close()
		outFmtCtx.AvformatFreeContext()
		return
	}
	outFmtCtx.SetOformat(outFmt)

	outFmtCtx.SetUrl(path)

	outCodec := avcodec.AvcodecFindEncoder(avcodec.CodecId(avcodec.AV_CODEC_ID_PCM_S16LE))
	if outCodec == nil {
		err = errors.New("could not find PCMS16LE encoder")
		_ = outIOCtx.Close()
		outFmtCtx.AvformatFreeContext()
		return
	}

	stream := avformat.AvformatNewStream(outFmtCtx, nil)
	if stream == nil {
		err = errors.New("could not create new stream")
		_ = outIOCtx.Close()
		outFmtCtx.AvformatFreeContext()
		return
	}

	avCtx := outCodec.AvcodecAllocContext3()
	if avCtx == nil {
		err = errors.New("could not allocate an encoding context")
		_ = outIOCtx.Close()
		outFmtCtx.AvformatFreeContext()
		return
	}

	*outCodecCtx.Flags() |= avformat.AVFMT_FLAG_BITEXACT

	avCtx.SetChannels(1)
	avCtx.SetSampleRate(22050)

	if outFmtCtx.Oformat().Flags() & avformat.AVFMT_GLOBALHEADER != 0 {
		*avCtx.Flags() |= avformat.AVFMT_GLOBALHEADER
	}

	ret := avCtx.AvcodecOpen2(outCodec, nil)
	if ret < 0 {
		err = fmt.Errorf("could not open output codec (error '%v')", avutil.ErrorFromCode(ret))
		_ = outIOCtx.Close()
		outFmtCtx.AvformatFreeContext()
		avCtx.AvcodecFreeContext()
		return
	}

	ret = avcodec.AvcodecParametersFromContext(stream.CodecParameters(), avCtx)
	if ret < 0 {
		err = fmt.Errorf("could not initialize stream parameters (error '%v')", avutil.ErrorFromCode(ret))
		_ = outIOCtx.Close()
		outFmtCtx.AvformatFreeContext()
		avCtx.AvcodecFreeContext()
		return
	}

	outCodecCtx = avCtx
	return
}

func initResampler(inCodecCtx *avcodec.Context, outCodecCtx *avcodec.Context) (resampleCtx *swresample.Context, err error) {
	resampleCtx = swresample.SwrAllocSetOpts(nil, avutil.GetDefaultChannelLayout(outCodecCtx.Channels()),
		(avutil.AVSampleFormat)(outCodecCtx.SampleFmt()), outCodecCtx.SampleRate(),
		avutil.GetDefaultChannelLayout(inCodecCtx.Channels()), (avutil.AVSampleFormat)(inCodecCtx.SampleFmt()),
		inCodecCtx.SampleRate(), 0, nil)
	if resampleCtx == nil {
		err = errors.New("could not allocate resample context")
		return
	}

	ret := resampleCtx.SwrInit()
	if ret < 0 {
		err = fmt.Errorf("could not open resample context (error '%v')", avutil.ErrorFromCode(ret))
		resampleCtx.SwrFree()
		return
	}
	return
}

func initFifo(outCodecCtx *avcodec.Context) (fifo *avutil.AVAudioFifo, err error) {
	fifo = avutil.AVAudioFifoAlloc(outCodecCtx.SampleFmt(), outCodecCtx.Channels(), 1)
	if fifo == nil {
		err = errors.New("could not allocate FIFO")
		return
	}
	return
}

func writeOutputFileHeader(outFmtCtx *avformat.Context) error {
	ret := outFmtCtx.AvformatWriteHeader(nil)
	if ret < 0 {
		return fmt.Errorf("could not write output file header (error '%v')", avutil.ErrorFromCode(ret))
	}
	return nil
}

func (t *transcoder) readDecodeConvertAndStore() (finished bool, err error) {
	inputFrame := avutil.AvFrameAlloc()
	if inputFrame == nil {
		err = errors.New("failed to allocate new frame")
		return
	}
	defer avutil.AvFrameFree(inputFrame)

	var convertedInputSamples [][]uint8 = nil
	defer func() {
		if convertedInputSamples != nil {
			avutil.AvFreep(unsafe.Pointer(&convertedInputSamples))
		}
	}()

	dataPresent, finished, err := decodeAudioFrame(inputFrame, t.inFmtCtx, t.inCodecCtx)
	if err != nil {
		return
	}

	if finished {
		return
	}

	if dataPresent {
		err = initConvertedSamples(&convertedInputSamples, t.outCodecCtx, inputFrame.Samples())
		if err != nil {
			err = fmt.Errorf("could not initialize converted samples array: %v", err)
			return
		}

		err = convertSamples(inputFrame.ExtendedData(), convertedInputSamples, inputFrame.Samples(), t.resampleCtx)
		if err != nil {
			err = fmt.Errorf("could not convert samples: %v", err)
			return
		}

		err = addSamplesToFifo(t.fifo, convertedInputSamples, inputFrame.Samples())
		if err != nil {
			err = fmt.Errorf("could not write samples to FIFO: %v", err)
			return
		}
	}
	return
}

func decodeAudioFrame(frame *avutil.Frame, inFmtCtx *avformat.Context, inCodecCtx *avcodec.Context) (dataPresent bool, finished bool, err error) {
	inputPacket := &avcodec.Packet{}
	inputPacket.AvInitPacket()
	inputPacket.SetData(nil)
	inputPacket.SetSize(0)
	defer inputPacket.AvPacketUnref()

	ret := inFmtCtx.AvReadFrame(inputPacket)
	if ret < 0 {
		if ret == avutil.AvErrorEOF {
			finished = true
		} else {
			err = fmt.Errorf("could not read frame (error '%v')", avutil.ErrorFromCode(ret))
			return
		}
	}

	ret = inCodecCtx.AvcodecSendPacket(inputPacket)
	if ret < 0 {
		err = fmt.Errorf("could not send packet for decoding (error '%v')", avutil.ErrorFromCode(ret))
		return
	}

	ret = inCodecCtx.AvcodecReceiveFrame(frame)
	if ret < avutil.AvErrorEAGAIN {
		return
	} else if ret == avutil.AvErrorEOF {
		finished = true
		return
	} else if ret < 0 {
		err = fmt.Errorf("could not decode frame (error '%v')", avutil.ErrorFromCode(ret))
		return
	} else {
		dataPresent = true
		return
	}
}

func initConvertedSamples(convertedInputSamples *[][]uint8, outCodecCtx *avcodec.Context, frameSize int) error {
	newMemory := C.calloc(outCodecCtx.Channels(), 8)  // TODO: use sizeof(**uint8_t) here
	if newMemory == nil {
		return errors.New("could not allocate converted input sample pointers")
	}
	*convertedInputSamples = newMemory

	ret := avutil.AVSamplesAlloc(*convertedInputSamples, outCodecCtx.Channels(), frameSize, outCodecCtx.SampleFmt(), 0)
	if ret < 0 {
		avutil.AvFreep(unsafe.Pointer(convertedInputSamples))
		C.free(*convertedInputSamples)
		return fmt.Errorf("could not allocate converted input samples (error '%v')", avutil.ErrorFromCode(ret))
	}
	return nil
}

func convertSamples(inputData [][]uint8, convertedData [][]uint8, frameSize int, swrCtx *swresample.Context) error {
	ret := swrCtx.SwrConvert(convertedData, frameSize, inputData, frameSize)
	if ret < 0 {
		return fmt.Errorf("could not convert input samples (error '%v')", avutil.ErrorFromCode(ret))
	}

	return nil
}

func addSamplesToFifo(fifo *avutil.AVAudioFifo, convertedInputSamples [][]uint8, frameSize int) error {
	ret := fifo.AVAudioFifoRealloc(fifo.Size() + frameSize)
	if ret < 0 {
		return fmt.Errorf("could not reallocate FIFO (error '%v')", avutil.ErrorFromCode(ret))
	}

	ret = fifo.AVAudioWrite(convertedInputSamples, frameSize)
	if ret < 0 {
		return fmt.Errorf("could not write data to FIFO (error '%v')", avutil.ErrorFromCode(ret))
	}
	return nil
}

func (t *transcoder) loadEncodeAndWrite() error {
	frameSize := t.outCodecCtx.FrameSize()
	if t.fifo.Size() < frameSize {
		frameSize = t.fifo.Size()
	}

	outputFrame, err := initOutputFrame(t.outCodecCtx, frameSize)
	if err != nil {
		return fmt.Errorf("could not initialize output frame: %v", err)
	}
	defer avutil.AvFrameFree(outputFrame)

	framesRead := t.fifo.AVAudioFifoRead(outputFrame.Data(), frameSize)
	if framesRead < frameSize {
		return errors.New("could not read data from FIFO")
	}

	_, err = t.encodeAudioFrame(outputFrame)
	return err
}

func initOutputFrame(outCodecCtx avcodec.Context, frameSize int) (frame *avutil.Frame, err error) {
	frame = avutil.AvFrameAlloc()
	if frame == nil {
		err = errors.New("could not allocate output frame")
		return
	}

	frame.SetSamples(frameSize)
	frame.SetChannelLayout(outCodecCtx.ChannelLayout())
	frame.SetFormat(outCodecCtx.SampleFmt())
	frame.SetSampleRate(outCodecCtx.SampleRate())

	ret := avutil.AvFrameGetBuffer(frame, 0)
	if ret < 0 {
		err = fmt.Errorf("could not allocate output frame samples (error '%v')", avutil.ErrorFromCode(ret))
		avutil.AvFrameFree(frame)
		return
	}

	return
}

func (t *transcoder) encodeAudioFrame(frame *avutil.Frame) (dataPresent bool, err error) {
	outputPacket := avcodec.Packet{}
	outputPacket.AvInitPacket()
	outputPacket.SetData(nil)
	outputPacket.SetSize(0)
	defer outputPacket.AvPacketUnref()

	if frame != nil {
		frame.SetPts(t.pts)
		t.pts += frame.Samples()
	}

	ret := t.outCodecCtx.AVcodecSendFrame(frame)
	if ret == avutil.AvErrorEOF {
		return
	} else if ret < 0 {
		err = fmt.Errorf("could not send packet for encoding (error '%v')", avutil.ErrorFromCode(ret))
		return
	}

	ret = t.outCodecCtx.AVcodecReceivePacket(outputPacket)
	if ret == avutil.AvErrorEAGAIN {
		return
	} else if ret == avutil.AvErrorEOF {
		return
	} else if ret < 0 {
		err = fmt.Errorf("could not encode frame (error '%v')", avutil.ErrorFromCode(ret))
		return
	} else {
		dataPresent = true
	}

	ret = t.outFmtCtx.AvWriteFrame(outputPacket)
	if ret < 0 {
		err = fmt.Errorf("could not write frame (error '%v')", avutil.ErrorFromCode(ret))
		return
	}
	return
}
