package main

import (
	"encoding/json"
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
	case "gofile":
		go app.GofileUpload(evt, fileName)
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

func (app *application) GetGofileServer() string {
	type gofileData struct {
		Server string `json:"server"`
	}

	type gofileResponse struct {
		Status string `json:"status"`
		Data   gofileData
	}

	response, err := http.Get("https://api.gofile.io/getServer")
	if err != nil {
		return ""
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return ""
	}
	var responseObject gofileResponse
	if err = json.Unmarshal(responseData, &responseObject); err != nil {
		return ""
	}

	uploadServer := fmt.Sprintf("https://%s.gofile.io/uploadFile", responseObject.Data.Server)
	return uploadServer
}

func (app *application) GofileUpload(evt *event.Event, fileName string) {
	defer os.Remove(fileName)
	app.SendText(evt, "Uploading to gofile.io...")
	pr, pw := io.Pipe()
	form := multipart.NewWriter(pw)

	type gofileData struct {
		DownloadPage string `json:"downloadPage"`
		Code         string `json:"code"`
		ParentFolder string `json:"parentFolder"`
		FileId       string `json:"fileId"`
		FileName     string `json:"fileName"`
		Md5          string `json:"md5"`
	}

	type gofileResponse struct {
		Status string `json:"status"`
		Data   gofileData
	}

	go func() {
		defer pw.Close()

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

	gofileServer := app.GetGofileServer()
	app.SendText(evt, gofileServer)
	req, err := http.NewRequest(http.MethodPost, gofileServer, pr)
	if err != nil {
		app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
		os.Remove(fileName)
		return
	}
	req.Header.Set("Content-Type", form.FormDataContentType())

	httpClient := http.Client{
		Timeout: 300 * time.Second,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
		os.Remove(fileName)
		app.Log.Error().Err(err).Msgf("Error while sending HTTP request: %s", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		app.Log.Error().Err(err).Msg("Failed to upload")
		os.Remove(fileName)
		app.SendText(evt, fmt.Sprintf("Something went wrong FeelsBadMan: %q", err))
		return
	}

	jsonResponse := new(gofileResponse)
	if err := json.Unmarshal(body, jsonResponse); err != nil {
		app.Log.Error().Err(err).Msg("Failed to upload")
		os.Remove(fileName)
		app.SendText(evt, fmt.Sprintf("Error unmarshalling json response: %q", err))
		return
	}

	var reply = jsonResponse.Data.DownloadPage

	app.SendText(evt, reply)
}
