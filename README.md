# Musical Fiesta
A command line tool for running a CSGO soundboard

# Installation
A prebuilt version of the tool is not yet available so you will need to build the tool yourself.
## macOS
1. Install [`golang`](https://golang.org/dl/)
2. Install [`ffmpeg`](https://ffmpeg.org/download.html)
3. Ensure `ffmpeg` and go are both accessible from your terminal. You can test this by running the following commands in your terminal:
```
$: which go
```
and getting an output like:
```
/usr/local/bin/go
```
and running
```
$: which ffmpeg
```
and getting an output like:
```
/usr/local/bin/ffmpeg
```
Note that the output does not need to be exactly the same, it only needs to not show an error.
If either of these commands fail, you will need to add the directories where you installed those tools to your PATH environment variable.
4. Clone the git repository to your machine and navigate to it in your terminal
5. run `go run src/main.go`
6. Musical Fiesta will start automatically. To run again in the future, you can run the executable called `main` that you just created by double-clicking on it.

## Windows
1. Install [`golang`](https://golang.org/dl/)
2. Install [`ffmpeg`](https://ffmpeg.org/download.html)
3. Ensure `ffmpeg` and go are both accessible from `cmd`. You can test this by running the following commands in `cmd`:
```
$: where go
```
and getting an output like:
```
C:\Program Files (x86)/go/go
```
and running
```
$: which ffmpeg
```
and getting an output like:
```
C:\Program Files (x86)/ffmepg/ffmpeg
```
Note that the output does not need to be exactly the same, it only needs to not show an error.
If either of these commands fail, you will need to add the directories where you installed those tools to your PATH environment variable. You can do this in Control Panel > System > Advanced System Settings > Environment Variables
4. Clone the git repository to your machine and navigate to it in cmd
5. Run `go run src/main.go`
6. Musical Fiesta will start automatically. To run again in the future, you can run the executable called `main` that you just created by double-clicking on it.
