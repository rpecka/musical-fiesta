package commands

import (
	"fiesta/src/library"
	"github.com/desertbit/grumble"
)

func SetUpCommands(app *grumble.App, library library.Library) {
	addDeleteTrack(app, library)
	addTrack(app, library)
}
