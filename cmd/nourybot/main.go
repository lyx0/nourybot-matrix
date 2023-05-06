package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/chzyer/readline"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto/cryptohelper"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type config struct {
	matrixHomeserver string
	matrixUser       string
	matrixPass       string
	database         string
}

type Application struct {
	Mc  *mautrix.Client
	Log zerolog.Logger
}

var debug = flag.Bool("debug", false, "Enable debug logs")

func main() {
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	var cfg config
	cfg.matrixHomeserver = os.Getenv("MATRIX_HOMESERVER")
	cfg.matrixUser = os.Getenv("MATRIX_USERNAME")
	cfg.matrixPass = os.Getenv("MATRIX_PASSWORD")
	cfg.database = os.Getenv("SQLITE_DATABASE")

	client, err := mautrix.NewClient(cfg.matrixHomeserver, "", "")
	if err != nil {
		panic(err)
	}

	rl, err := readline.New("[no room]> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()
	log := zerolog.New(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.Out = rl.Stdout()
		w.TimeFormat = time.Stamp
	})).With().Timestamp().Logger()
	if !*debug {
		log = log.Level(zerolog.InfoLevel)
	}
	client.Log = log

	var lastRoomID id.RoomID

	app := &Application{
		Mc:  client,
		Log: log,
	}

	syncer := app.Mc.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) {
		lastRoomID = evt.RoomID
		rl.SetPrompt(fmt.Sprintf("%s> ", lastRoomID))
		app.Log.Info().
			Str("sender", evt.Sender.String()).
			Str("type", evt.Type.String()).
			Str("id", evt.ID.String()).
			Str("body", evt.Content.AsMessage().Body).
			Msg("Received message")

		app.ParseEvent(evt)
	})
	syncer.OnEventType(event.StateMember, func(source mautrix.EventSource, evt *event.Event) {
		if evt.GetStateKey() == app.Mc.UserID.String() && evt.Content.AsMember().Membership == event.MembershipInvite {
			_, err := app.Mc.JoinRoomByID(evt.RoomID)
			if err == nil {
				lastRoomID = evt.RoomID
				rl.SetPrompt(fmt.Sprintf("%s> ", lastRoomID))
				app.Log.Info().
					Str("room_id", evt.RoomID.String()).
					Str("inviter", evt.Sender.String()).
					Msg("Joined room after invite")
			} else {
				app.Log.Error().Err(err).
					Str("room_id", evt.RoomID.String()).
					Str("inviter", evt.Sender.String()).
					Msg("Failed to join room after invite")
			}
		}
	})

	cryptoHelper, err := cryptohelper.NewCryptoHelper(app.Mc, []byte("meow"), cfg.database)
	if err != nil {
		panic(err)
	}

	cryptoHelper.LoginAs = &mautrix.ReqLogin{
		Type:       mautrix.AuthTypePassword,
		Identifier: mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: cfg.matrixUser},
		Password:   cfg.matrixPass,
	}

	err = cryptoHelper.Init()
	if err != nil {
		panic(err)
	}

	app.Mc.Crypto = cryptoHelper

	app.Log.Info().Msg("Now running")
	syncCtx, cancelSync := context.WithCancel(context.Background())
	var syncStopWait sync.WaitGroup
	syncStopWait.Add(1)

	go func() {
		err = app.Mc.SyncWithContext(syncCtx)
		defer syncStopWait.Done()
		if err != nil && !errors.Is(err, context.Canceled) {
			panic(err)
		}
	}()

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		if lastRoomID == "" {
			app.Log.Error().Msg("Wait for an incoming message before sending messages")
			continue
		}
		resp, err := app.Mc.SendText(lastRoomID, line)
		if err != nil {
			app.Log.Error().Err(err).Msg("Failed to send event")
		} else {
			app.Log.Info().Str("event_id", resp.EventID.String()).Msg("Event sent")
		}
	}
	cancelSync()
	syncStopWait.Wait()
	err = cryptoHelper.Close()
	if err != nil {
		app.Log.Error().Err(err).Msg("Error closing database")
	}
}
