package main

import (
	"maunium.net/go/mautrix/event"
)

func (app *Application) ParseEvent(evt *event.Event) {
	if evt.Content.AsMessage().Body == "xD" {
		resp, err := app.Mc.SendText(evt.RoomID, "hehe")
		if err != nil {
			app.Log.Error().Err(err).Msg("Failed to send event")
		} else {
			app.Log.Info().Str("event_id", resp.EventID.String()).Msg("Event sent")
		}
	}
}
