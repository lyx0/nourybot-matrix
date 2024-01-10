package main

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"maunium.net/go/mautrix/event"
)

func (app *application) NewUpload(destination string, evt *event.Event, fileName string) {
	switch destination {
	case "yaf":
		go app.YafUpload(evt, fileName)
		//case "catbox":
		//	go app.CatboxUpload(target, fileName, identifier, msg)
		//case "kappa":
		//	go app.KappaUpload(target, fileName, identifier, msg)
		//case "gofile":
		//	go app.GofileUpload(target, fileName, identifier, msg)
	}
}

func (app *application) YafUpload(evt *event.Event, fileName string) {
	defer os.Remove(fileName)
	app.SendText(evt, "Uploading to yaf.li...")
	pr, pw := io.Pipe()
	form := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()

		err := form.WriteField("name", "xd")
		if err != nil {
			app.Log.Error().Err(err).Msg("Failed to upload")
			app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
			return
		}

		file, err := os.Open(fileName) // path to image file
		if err != nil {
			app.Log.Error().Err(err).Msg("Failed to upload")
			os.Remove(fileName)
			app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
			return
		}

		w, err := form.CreateFormFile("file", fileName)
		if err != nil {
			app.Log.Error().Err(err).Msg("Failed to upload")
			os.Remove(fileName)
			app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
			return
		}

		_, err = io.Copy(w, file)
		if err != nil {
			app.Log.Error().Err(err).Msg("Failed to upload")
			os.Remove(fileName)
			app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
			return
		}

		form.Close()
	}()

	req, err := http.NewRequest(http.MethodPost, YAF_ENDPOINT, pr)
	if err != nil {
		app.Log.Error().Err(err).Msg("Failed to upload")
		os.Remove(fileName)
		app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
		return
	}
	req.Header.Set("Content-Type", form.FormDataContentType())

	httpClient := http.Client{
		Timeout: 300 * time.Second,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		app.Log.Error().Err(err).Msg("Failed to upload")
		os.Remove(fileName)
		app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		app.Log.Error().Err(err).Msg("Failed to upload")
		app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
		os.Remove(fileName)
		return
	}

	var reply = string(body[:])

	app.SendText(evt, reply)
}
