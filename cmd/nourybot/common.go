package main

import "maunium.net/go/mautrix/event"

func (app *Application) SendText(evt *event.Event, message string) {
	room := evt.RoomID

	resp, err := app.Mc.SendText(room, message)
	if err != nil {
		app.Log.Error().Err(err).Msg("Failed to send event")
	} else {
		app.Log.Info().Str("event_id", resp.EventID.String()).Msg("Event sent")
	}
}
