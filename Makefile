
CC = gcc
AR = ar
CD = cd

SRC_DIR = src
FFWRAPPER_DIR = ${SRC_DIR}/ffwrapper

GO_FILES = $(shell find src/ -type f -name '*.go')

BUILD_DIR = build
LIB_DIR=lib

FIESTA_BASE_NAME = fiesta

ifeq (${OS},Windows_NT)
	FIESTA_EXECUTABLE_NAME = ${FIESTA_BASE_NAME}.exe
	CGO_LDFLAGS+=-static
	CGO_LDFLAGS+=-L/mingw64/usr/local/lib
else
	FIESTA_EXECUTABLE_NAME = ${FIESTA_BASE_NAME}

	ifndef LZMA_PATH
		$(error LZMA_PATH is not set. The full path to liblzma.a must be \
			located at LZMA_PATH. Did you install xz?)
	endif
	CGO_LDFLAGS+=${LZMA_PATH}
endif

FFWRAPPER_CFLAGS=$(shell pkg-config --static --cflags libavformat)

.DEFAULT_GOAL = ${FIESTA_EXECUTABLE_NAME}

clean-libs:
	rm -rf ${LIB_DIR}

clean-builds:
	rm -rf ${BUILD_DIR}

clean-fiesta:
	rm -f ${FIESTA_EXECUTABLE_NAME}

clean: clean-libs clean-builds clean-fiesta
	echo done

${BUILD_DIR}/ffwrapper.o: ${FFWRAPPER_DIR}/ffwrapper.c ${FFWRAPPER_DIR}/ffwrapper.h | ${BUILD_DIR}
	${CC} -c ${FFWRAPPER_DIR}/ffwrapper.c ${FFWRAPPER_CFLAGS} -o ${BUILD_DIR}/ffwrapper.o

${LIB_DIR}/libffwrapper.a: ${BUILD_DIR}/ffwrapper.o | ${LIB_DIR}
	${AR} rs ${LIB_DIR}/libffwrapper.a ${BUILD_DIR}/ffwrapper.o

${BUILD_DIR}:
	mkdir -p ${BUILD_DIR}

${LIB_DIR}:
	mkdir -p ${LIB_DIR}

${FIESTA_EXECUTABLE_NAME}: ${LIB_DIR}/libffwrapper.a ${GO_FILES}
	CGO_LDFLAGS="${CGO_LDFLAGS}" go build -o ${FIESTA_EXECUTABLE_NAME} src/main.go
