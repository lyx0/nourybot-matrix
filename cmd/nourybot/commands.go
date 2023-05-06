package main

import (
	"strings"

	"github.com/lyx0/nourybot-matrix/internal/commands"
	"maunium.net/go/mautrix/event"
)

func (app *Application) ParseCommand(evt *event.Event) {
	commandName := strings.ToLower(strings.SplitN(evt.Content.AsMessage().Body, " ", 2)[0][1:])
	app.Log.Info().Msgf("Command: %s", commandName)

	switch commandName {
	case "xd":
		app.SendText(evt, "XD !")
		return

	case "ping":
		if resp, err := commands.Ping(); err != nil {
			app.Log.Error().Err(err).Msg("Failed to send Ping")
			app.SendText(evt, "Something went wrong.")
			return
		} else {
			app.SendText(evt, resp)
			return
		}

	case "xkcd":
		if resp, err := commands.Xkcd(); err != nil {
			app.Log.Error().Err(err).Msg("Failed to send Ping")
			app.SendText(evt, "Something went wrong.")
			return
		} else {
			app.SendText(evt, resp)
			return
		}
	}
}
