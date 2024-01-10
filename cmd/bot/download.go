package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/wader/goutubedl"
	"maunium.net/go/mautrix/event"
)

const (
	//CATBOX_ENDPOINT = "https://litterbox.catbox.moe/resources/internals/api.php"
	//KAPPA_ENDPOINT  = "https://kappa.lol/api/upload"
	YAF_ENDPOINT = "https://i.yaf.li/upload"
)

func (app *application) NewDownload(destination string, evt *event.Event, link string) {
	uuid := uuid.NewString()

	app.Log.Info().Msgf("Link: %s", link)
	switch destination {
	case "yaf":
		app.YafDownload(evt, link, uuid)
	// case "catbox":
	// 	app.CatboxDownload(target, link, identifier, msg)
	//case "kappa":
	//	app.KappaDownload(target, link, identifier, msg)
	case "gofile":
		app.GofileDownload(evt, link)
	}
}

func (app *application) YafDownload(evt *event.Event, link, uuid string) {
	goutubedl.Path = "yt-dlp"
	var rExt string

	app.SendText(evt, "Downloading...")
	result, err := goutubedl.New(context.Background(), link, goutubedl.Options{})
	if err != nil {
		app.Log.Error().Err(err).Msg("Failed to download")
		app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
		return
	}

	// For some reason youtube links return webm as result.Info.Ext but
	// are in reality mp4.
	if strings.HasPrefix(link, "https://www.youtube.com/") || strings.HasPrefix(link, "https://youtu.be/") {
		rExt = "mp4"
	} else {
		rExt = result.Info.Ext
	}

	downloadResult, err := result.Download(context.TODO(), "best")
	if err != nil {
		app.Log.Error().Err(err).Msg("Failed to download")
		app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
		return
	}
	defer downloadResult.Close()

	fileName := fmt.Sprintf("%s.%s", uuid, rExt)
	f, err := os.Create(fileName)
	if err != nil {
		app.Log.Error().Err(err).Msg("Failed to download")
		app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
		return
	}
	defer f.Close()
	io.Copy(f, downloadResult)

	app.NewUpload("yaf", evt, fileName)
}

func (app *application) GofileDownload(evt *event.Event, link string) {
	goutubedl.Path = "yt-dlp"
	var rExt string

	app.SendText(evt, "Downloading...")
	result, err := goutubedl.New(context.Background(), link, goutubedl.Options{})
	if err != nil {
		app.Log.Error().Err(err).Msg("Failed to download")
		app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
		return
	}
	// For some reason youtube links return webm as result.Info.Ext but
	// are in reality mp4.
	if strings.HasPrefix(link, "https://www.youtube.com/") || strings.HasPrefix(link, "https://youtu.be/") {
		rExt = "mp4"
	} else {
		rExt = result.Info.Ext
	}
	safeFilename := fmt.Sprintf("download_%s", result.Info.Title)
	downloadResult, err := result.Download(context.Background(), "best")
	if err != nil {
		app.Log.Error().Err(err).Msg("Failed to download")
		app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
		return
	}
	fileName := fmt.Sprintf("%s.%s", safeFilename, rExt)
	f, err := os.Create(fileName)
	//app.Send(target, fmt.Sprintf("Filename: %s", fileName), msg)

	if err != nil {
		app.Log.Error().Err(err).Msg("Failed to download")
		app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
		return
	}
	defer f.Close()
	if _, err = io.Copy(f, downloadResult); err != nil {
		app.Log.Error().Err(err).Msg("Failed to download")
		app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
		return
	}

	downloadResult.Close()
	f.Close()

	app.SendText(evt, fileName)
	go app.NewUpload("gofile", evt, fileName)

}
