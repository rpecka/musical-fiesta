.DEFAULT_GOAL = default

clean:
	rm fiesta

lzma:
ifndef LZMA_PATH
	$(error LZMA_PATH is not set. The full path to liblzma.a must be located at LZMA_PATH. Did you install xz?)
endif


default: lzma
	CGO_LDFLAGS=${LZMA_PATH} go build -o fiesta src/main.go
