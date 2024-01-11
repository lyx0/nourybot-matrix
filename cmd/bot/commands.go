package main

import (
	"strings"

	"github.com/lyx0/nourybot-matrix/pkg/owm"
	"maunium.net/go/mautrix/event"
)

func (app *application) ParseCommand(evt *event.Event) {
	var reply string
	// commandName is the actual name of the command without the prefix.
	// e.g. `!ping` would be `ping`.
	commandName := strings.ToLower(strings.SplitN(evt.Content.AsMessage().Body, " ", 2)[0][1:])

	// cmdParams are additional command parameters.
	// e.g. `!weather san antonio`
	// cmdParam[0] is `san` and cmdParam[1] = `antonio`.
	cmdParams := strings.SplitN(evt.Content.AsMessage().Body, " ", 500)
	app.Log.Info().Msgf("cmdParams: %s", cmdParams)

	// msgLen is the amount of words in a message without the prefix.
	// Useful to check if enough cmdParams are provided.
	//msgLen := len(strings.SplitN(evt.Content.AsMessage().Body, " ", -2))

	app.Log.Info().Msgf("Command: %s", evt.Content.AsMessage().Body)

	//message := evt.Content.AsMessage().Body
	switch commandName {
	case "xd":
		app.SendText(evt, "xd !")
		return
	case "gofile":
		app.NewDownload("gofile", evt, cmdParams[1])
		return

	case "weather":
		reply, _ = owm.Weather(evt.Content.AsMessage().Body[9:len(evt.Content.AsMessage().Body)])

	case "yaf":
		app.NewDownload("yaf", evt, cmdParams[1])
		return

	}
	if reply != "" {
		app.SendText(evt, reply)
	}
}
