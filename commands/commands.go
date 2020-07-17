package commands

import (
	"fiesta/library"
	"github.com/desertbit/grumble"
)

func SetUpCommands(app *grumble.App, library library.Library) {
	addDeleteTrack(app, library)
}
