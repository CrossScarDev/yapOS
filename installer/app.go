package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetPin(serialNumber string) map[string]any {
	const possibleChars string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	indem := ""
	for range 16 {
		indem += string(possibleChars[rand.IntN(len(possibleChars))])
	}

	req, err := http.NewRequest("GET", "https://play.date/api/v2/device/register/"+serialNumber+"/get", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Idempotency-Key", indem)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var res map[string]any
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

var accessToken string = ""

func (a *App) FinishRegistration(serialNumber string) map[string]any {
	resp, err := http.Get("https://play.date/api/v2/device/register/" + serialNumber + "/complete/get")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var res map[string]any
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Fatal(err)
	}

	return res
}

func (a *App) DownloadOS(accessToken string) {
	req, err := http.NewRequest("GET", "https://play.date/api/v2/firmware/?current_version=2.6.1", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Token "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var res map[string]string
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Fatal(err)
	}

	decryptKey, err := base64.StdEncoding.DecodeString(res["decryption_key"])
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.CreateTemp("", "PlaydateOS.*.pdos")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	resp, err = http.Get(res["url"])
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	f, err = os.CreateTemp("", "PlaydateOS.*.pdkey")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.Write(decryptKey)
	if err != nil {
		log.Fatal(err)
	}
}
