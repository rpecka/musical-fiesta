package commands

import (
	"fiesta/src/library"
	"fiesta/src/settings"
	"github.com/desertbit/grumble"
)

func SetUpCommands(app *grumble.App, settings *settings.Settings, library *library.Library) {
	addDeleteTrack(app, library)
	addTrack(app, library)
	addStart(app, settings, library)
}
