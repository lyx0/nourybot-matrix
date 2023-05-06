package main

import (
	"strings"

	"github.com/lyx0/nourybot-matrix/internal/commands"
	"maunium.net/go/mautrix/event"
)

func (app *Application) ParseCommand(evt *event.Event) {
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
	msgLen := len(strings.SplitN(evt.Content.AsMessage().Body, " ", -2))

	app.Log.Info().Msgf("Command: %s", commandName)

	switch commandName {
	case "xd":
		app.SendText(evt, "XD !")
		return

	case "currency":
		if msgLen < 4 {
			app.SendText(evt, "Not enough arguments provided")
			return
		} else {
			if resp, err := commands.Currency(cmdParams[1], cmdParams[2], cmdParams[4]); err != nil {
				app.Log.Error().Err(err).Msg("failed to handle currency command")
				app.SendText(evt, "Something went wrong.")
				return
			} else {
				app.SendText(evt, resp)
				return
			}
		}

	case "phonetic":
		msg := evt.Content.AsMessage().Body[9:len(evt.Content.AsMessage().Body)]
		if resp, err := commands.Phonetic(msg); err != nil {
			app.Log.Error().Err(err).Msg("failed to handle phonetic command")
			app.SendText(evt, "Something went wrong.")
			return
		} else {
			app.SendText(evt, resp)
			return
		}

	case "ping":
		if resp, err := commands.Ping(); err != nil {
			app.Log.Error().Err(err).Msg("failed to handle ping command")
			app.SendText(evt, "Something went wrong.")
			return
		} else {
			app.SendText(evt, resp)
			return
		}

	case "random":
		if msgLen == 2 && cmdParams[1] == "xkcd" {
			if resp, err := commands.RandomXkcd(); err != nil {
				app.Log.Error().Err(err).Msg("failed to handle random->xkcd command")
				app.SendText(evt, "Something went wrong.")
				return
			} else {
				app.SendText(evt, resp)
				return
			}
		}

	case "xkcd":
		if msgLen == 2 && cmdParams[1] == "random" {
			if resp, err := commands.RandomXkcd(); err != nil {
				app.Log.Error().Err(err).Msg("failed to handle xkcd->randomXkcd command")
				app.SendText(evt, "Something went wrong.")
				return
			} else {
				app.SendText(evt, resp)
				return
			}
		} else if msgLen == 1 {
			if resp, err := commands.Xkcd(); err != nil {
				app.Log.Error().Err(err).Msg("failed to handle xkcd command")
				app.SendText(evt, "Something went wrong.")
				return
			} else {
				app.SendText(evt, resp)
				return
			}
		}

	case "weather":
		if msgLen < 2 {
			app.SendText(evt, "Not enough arguments provided")
			return
		} else {
			location := evt.Content.AsMessage().Body[9:len(evt.Content.AsMessage().Body)]
			if resp, err := commands.Weather(location); err != nil {
				app.Log.Error().Err(err).Msg("Failed to handle Weather command")
				app.SendText(evt, "Something went wrong.")
				return
			} else {
				app.SendText(evt, resp)
				return
			}
		}
	}
}
