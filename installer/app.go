package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/rand/v2"
	"net/http"
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

func (a *App) GetPin(serialNumber string) (map[string]any, error) {
	const possibleChars string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	indem := ""
	for range 16 {
		indem += string(possibleChars[rand.IntN(len(possibleChars))])
	}

	req, err := http.NewRequest("GET", "https://play.date/api/v2/device/register/"+serialNumber+"/get", nil)
	if err != nil {
		return nil, errors.New("Failed to construct request for: " + "https://play.date/api/v2/device/register/" + serialNumber + "/get")
	}
	req.Header.Set("Idempotency-Key", indem)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Failed to fetch: " + "https://play.date/api/v2/device/register/" + serialNumber + "/get")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Failed to read body of: " + "https://play.date/api/v2/device/register/" + serialNumber + "/get")
	}
	var res map[string]any
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, errors.New("Failed to parse json response of: " + "https://play.date/api/v2/device/register/" + serialNumber + "/get")
	}

	return res, nil
}
