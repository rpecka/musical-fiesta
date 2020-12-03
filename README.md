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
	2. Configure FFmpeg: `./configure --enable-static --disable-shared --disable-debug --disable-doc --disable-asm --disable-securetransport --prefix=$(pwd)`
	3. Build and install FFmpeg: `make install`
	4. ``export FFMPEG_INCLUDE_PATH=`pwd`/include"``
	5. ``export FFMPEG_LIB_PATH=`pwd`/lib``
3. Install `xz`
	1. `brew install xz`
	2. Fiesta requires that the path to the static library liblzma.a be defined in the variable `LZMA_PATH` i.e. `export LZMA_PATH=/usr/local/Cellar/xz/5.2.5/lib/liblzma.a`
4. Build fiesta using `make`

## Windows
1. Install [`golang`](https://golang.org/dl/)
2. Install [`ffmpeg`](https://ffmpeg.org/download.html)
3. Ensure `ffmpeg` and go are both accessible from `cmd`. You can test this by running the following commands in `cmd`:
```
$: where go
```
and getting an output like:
```
C:\Go\bin\go.exe
```
and running
```
$: where ffmpeg
```
and getting an output like:
```
C:\Program Files (x86)\ffmepg\ffmpeg.exe
```
Note that the output does not need to be exactly the same, it only needs to not show an error.
If either of these commands fail, you will need to add the directories where you installed those tools to your PATH environment variable. You can do this in Control Panel > System > Advanced System Settings > Environment Variables

4. Clone the git repository to your machine and navigate to it in cmd
5. Run `go run src\main.go`
6. Musical Fiesta will start automatically. To run again in the future, you can run the executable called `main` that you just created by double-clicking on it.
