
CC = clang
AR = ar
CD = cd

SRC_DIR = src
FFWRAPPER_DIR = ${SRC_DIR}/ffwrapper

GO_FILES = $(shell find src/ -type f -name '*.go')

BUILD_DIR = build
LIB_DIR=$(shell pwd)/lib

.DEFAULT_GOAL = fiesta

clean-libs:
	rm -rf ${LIB_DIR}

clean-builds:
	rm -rf ${BUILD_DIR}

clean-fiesta:
	rm fiesta

clean: clean-fiesta clean-libs clean-builds
	echo done

${BUILD_DIR}/ffwrapper.o: ${FFWRAPPER_DIR}/ffwrapper.c ${FFWRAPPER_DIR}/ffwrapper.h | ${BUILD_DIR}
ifndef FFMPEG_INCLUDE_PATH
	$(error FFMPEG_INCLUDE_PATH is not set. The full path to FFmpeg's includes must be located at FFMPEG_INCLUDE_PATH. \
		Did you build ffmpeg?)
endif
	${CC} -c ${FFWRAPPER_DIR}/ffwrapper.c -I${FFMPEG_INCLUDE_PATH} -o ${BUILD_DIR}/ffwrapper.o

${LIB_DIR}/libffwrapper.a: ${BUILD_DIR}/ffwrapper.o | ${LIB_DIR}
	${AR} rs ${LIB_DIR}/libffwrapper.a ${BUILD_DIR}/ffwrapper.o

${BUILD_DIR}:
	mkdir ${BUILD_DIR}

${LIB_DIR}:
	mkdir ${LIB_DIR}

fiesta: ${LIB_DIR}/libffwrapper.a ${GO_FILES}
ifndef FFMPEG_LIB_PATH
	$(error FFMPEG_LIB_PATH is not set. The full path to FFmpeg's libraries must be located at FFMPEG_LIB_PATH. \
		Did you build ffmpeg?)
endif
ifndef LZMA_PATH
	$(error LZMA_PATH is not set. The full path to liblzma.a must be located at LZMA_PATH. Did you install xz?)
endif
	CGO_LDFLAGS="${LZMA_PATH} -L${FFMPEG_LIB_PATH} -L${LIB_DIR} -lffwrapper" go build -o fiesta src/main.go
