// Copyright (C) 2017 Tulir Asokan
// Copyright (C) 2018-2020 Luca Weiss
// Copyright (C) 2023 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

var debug = flag.Bool("debug", false, "Enable debug logs")
var env = flag.String("env", "dev", "Environment to run in (dev/prod)")

// var database = flag.String("database", "test.db", "SQLite database path")
type config struct {
	homeserver    string
	username      string
	password      string
	database      string
	db_account_id string
}

type application struct {
	MatrixClient *mautrix.Client
	Log          zerolog.Logger
}

func main() {
	var cfg config
	flag.Parse()
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	if *env == "prod" {
		cfg.homeserver = os.Getenv("PROD_MATRIX_HOMESERVER")
		cfg.username = os.Getenv("PROD_MATRIX_USERNAME")
		cfg.password = os.Getenv("PROD_MATRIX_PASSWORD")
		cfg.database = os.Getenv("PROD_SQLITE_DATABASE")
		cfg.db_account_id = os.Getenv("PROD_DB_ACCOUNT_ID")
	} else {
		cfg.homeserver = os.Getenv("DEV_MATRIX_HOMESERVER")
		cfg.username = os.Getenv("DEV_MATRIX_USERNAME")
		cfg.password = os.Getenv("DEV_MATRIX_PASSWORD")
		cfg.database = os.Getenv("DEV_SQLITE_DATABASE")
		cfg.db_account_id = os.Getenv("DEV_DB_ACCOUNT_ID")
	}

	client, err := mautrix.NewClient(cfg.homeserver, "", "")
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

	app := &application{
		MatrixClient: client,
		Log:          log,
	}
	var lastRoomID id.RoomID

	syncer := app.MatrixClient.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) {
		lastRoomID = evt.RoomID
		rl.SetPrompt(fmt.Sprintf("%s> ", lastRoomID))
		if evt.Content.AsMessage().Body[:1] == "!" {
			app.Log.Info().
				Str("sender", evt.Sender.String()).
				Str("type", evt.Type.String()).
				Str("id", evt.ID.String()).
				Str("body", evt.Content.AsMessage().Body).
				Msg("Received  xdddddddddddddddddddddddd message")

			app.ParseCommand(evt)
			return
		} else {
			app.ParseEvent(evt)
			return
		}
	})
	syncer.OnEventType(event.StateMember, func(source mautrix.EventSource, evt *event.Event) {
		if evt.GetStateKey() == app.MatrixClient.UserID.String() && evt.Content.AsMember().Membership == event.MembershipInvite {
			_, err := app.MatrixClient.JoinRoomByID(context.TODO(), evt.RoomID)
			if err == nil {
				lastRoomID = evt.RoomID
				rl.SetPrompt(fmt.Sprintf("%s> ", lastRoomID))
				log.Info().
					Str("room_id", evt.RoomID.String()).
					Str("inviter", evt.Sender.String()).
					Msg("Joined room after invite")
			} else {
				log.Error().Err(err).
					Str("room_id", evt.RoomID.String()).
					Str("inviter", evt.Sender.String()).
					Msg("Failed to join room after invite")
			}
		}

	})

	cryptoHelper, err := cryptohelper.NewCryptoHelper(app.MatrixClient, []byte("meow"), cfg.database)
	if err != nil {
		panic(err)
	}

	// You can also store the user/device IDs and access token and put them in the client beforehand instead of using LoginAs.
	//client.UserID = "..."
	//client.DeviceID = "..."
	//client.AccessToken = "..."
	// You don't need to set a device ID in LoginAs because the crypto helper will set it for you if necessary.
	cryptoHelper.LoginAs = &mautrix.ReqLogin{
		Type:       mautrix.AuthTypePassword,
		Identifier: mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: cfg.username},
		Password:   cfg.password,
	}
	// If you want to use multiple clients with the same DB, you should set a distinct database account ID for each one.
	cryptoHelper.DBAccountID = cfg.db_account_id
	err = cryptoHelper.Init()
	if err != nil {
		panic(err)
	}
	// Set the client crypto helper in order to automatically encrypt outgoing messages
	app.MatrixClient.Crypto = cryptoHelper

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()

	log.Info().Msg("Now running")
	syncCtx, cancelSync := context.WithCancel(context.Background())
	var syncStopWait sync.WaitGroup
	syncStopWait.Add(1)

	go func() {
		err = app.MatrixClient.SyncWithContext(syncCtx)
		defer syncStopWait.Done()
		if err != nil && !errors.Is(err, context.Canceled) {
			panic(err)
		}
	}()

	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}
		if lastRoomID == "" {
			log.Error().Msg("Wait for an incoming message before sending messages")
			continue
		}
		resp, err := app.MatrixClient.SendText(context.TODO(), lastRoomID, line)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send event")
		} else {
			log.Info().Str("event_id", resp.EventID.String()).Msg("Event sent")
		}
	}
	cancelSync()
	syncStopWait.Wait()
	err = cryptoHelper.Close()
	if err != nil {
		log.Error().Err(err).Msg("Error closing database")
	}
}

func (app *application) SendText(evt *event.Event, message string) {
	room := evt.RoomID

	resp, err := app.MatrixClient.SendText(context.TODO(), room, message)
	if err != nil {
		app.Log.Error().Err(err).Msg("Failed to send event")
	} else {
		app.Log.Info().Str("event_id", resp.EventID.String()).Msg("Event sent")
	}
}
