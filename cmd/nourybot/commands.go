package main

import (
	"strings"

	"maunium.net/go/mautrix/event"
)

func (app *Application) ParseCommand(evt *event.Event) {
	// commandName is the actual name of the command without the prefix.
	// e.g. `!ping` would be `ping`.
	//commandName := strings.ToLower(strings.SplitN(evt.Content.AsMessage().Body, " ", 2)[0][1:])

	// cmdParams are additional command parameters.
	// e.g. `!weather san antonio`
	// cmdParam[0] is `san` and cmdParam[1] = `antonio`.
	cmdParams := strings.SplitN(evt.Content.AsMessage().Body, " ", 500)
	app.Log.Info().Msgf("cmdParams: %s", cmdParams)

	// msgLen is the amount of words in a message without the prefix.
	// Useful to check if enough cmdParams are provided.
	//msgLen := len(strings.SplitN(evt.Content.AsMessage().Body, " ", -2))

	app.Log.Info().Msgf("Command: %s", evt.Content.AsMessage().Body)

	switch evt.Content.AsMessage().Body {
	case "!xd":
		app.SendText(evt, "xd !")
		return
	}
}
