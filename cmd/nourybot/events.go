package main

import "maunium.net/go/mautrix/event"

func (app *Application) ParseEvent(evt *event.Event) {
	// TODO:
	// Log the events or whatever, I don't even know what events there all are rn.
	app.Log.Info().Msgf("Event: %s", evt.Content.AsMessage().Body)
}
