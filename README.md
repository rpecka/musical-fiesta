# Musical Fiesta
A command line tool for running a CSGO soundboard

![Go](https://github.com/rpecka/musical-fiesta/workflows/Go/badge.svg?branch=master)

# Installation
A prebuilt version of the tool is not yet available so you will need to build the tool yourself.

# Building
Please note that I've only ever built this on my own machines and it's not unlikely that you will run into issues. Please don't hesitate to reach out by filing an issue.
## macOS
1. Install [`golang`](https://golang.org/dl/)
2. Build FFmpeg
	1. Clone FFmpeg https://github.com/FFmpeg/FFmpeg
	2. Check out version `n4.3.1`
	3. Configure FFmpeg: `./configure --pkg-config-flags=--static --enable-static --disable-shared --disable-debug --disable-doc --disable-asm --disable-network --disable-securetransport --disable-programs --disable-avdevice --disable-swscale --disable-avfilter --prefix=$(pwd)`
	4. Build and install FFmpeg: `make install`
	5. ``export FFMPEG_INCLUDE_PATH=`pwd`/include"``
	6. ``export FFMPEG_LIB_PATH="`pwd`/lib"``
3. Install `xz`
	1. `brew install xz`
	2. Fiesta requires that the path to xz's pkgconfig is in `PKG_CONFIG_PATH`
4. Build fiesta using `make`

## Windows
Building on Windows has only been tested with MSYS2 amd MinGW-w64, which you can get at http://msys2.github.io/

1. Install `go`: `pacman -S mingw64/mingw-w64-x86_64-go`
2. Build FFmpeg using the MinGW-w64 console
	1. Install pkg-config: `pacman -S mingw64/mingw-w64-x86_64-pkg-config`
	2. Clone FFmpeg https://github.com/FFmpeg/FFmpeg
	3. Check out version `n4.3.1`
	4. Configure FFmpeg: `./configure --pkg-config-flags=--static --enable-static --disable-shared --disable-debug --disable-doc --disable-asm --disable-network --disable-securetransport --disable-x86asm --disable-avdevice --disable-swscale --disable-avfilter --disable-programs --prefix=/mingw64/usr/local`
	5. Build and install FFmpeg: `make install`
	6. Ensure your `PKG_CONFIG_PATH` environment variable contains the directory with the FFmpeg .pc files. If you followed the configure instruction above, this is `/mingw64/usr/local/lib/pkgconfig`
3. Build fiesta using `make`
