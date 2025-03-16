package main

import (
	"archive/zip"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
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

var (
	pdosFilename  string
	pdkeyFilename string
)

func (a *App) DownloadPlaydateOS(accessToken string) {
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
	pdosFilename = f.Name()
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
	pdkeyFilename = f.Name()
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err = f.Write(decryptKey)
	if err != nil {
		log.Fatal(err)
	}
}

type OSInfo struct {
	filename       string
	targetFilename string
	url            string
}

var operatingSystems map[string]OSInfo = map[string]OSInfo{}

func (a *App) DownloadOS(selectedOS string, url string, filename string, targetFilename string) {
	f, err := os.CreateTemp("", filename)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	operatingSystems[selectedOS] = OSInfo{f.Name(), targetFilename, url}
}

var pdosExtractPath string = ""

func (a *App) ExtractPlaydateOS(funnyloader bool) {
	extractPath, err := os.MkdirTemp("", "PlaydateOS.*")
	if err != nil {
		log.Fatal(err)
	}

	zipReader, err := zip.OpenReader(pdosFilename)
	if err != nil {
		log.Fatal(err)
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		filePath := filepath.Join(extractPath, f.Name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				log.Fatal(err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			log.Fatal(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			log.Fatal(err)
		}
		srcFile, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			log.Fatal(err)
		}

		dstFile.Close()
		srcFile.Close()
	}

	if funnyloader {
		err = os.Mkdir(filepath.Join(extractPath, "System", "Launchers"), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		err = os.Rename(filepath.Join(extractPath, "System", "Launcher.pdx"), filepath.Join(extractPath, "System", "Launchers", "StockLauncher.pdx"))
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = os.Rename(filepath.Join(extractPath, "System", "Launcher.pdx"), filepath.Join(extractPath, "System", "StockLauncher.pdx"))
		if err != nil {
			log.Fatal(err)
		}
	}

	pdosExtractPath = extractPath
}

func (a *App) ExtractOS(selectedOS string, pdxPath string) {
	extractPath, err := os.MkdirTemp("", selectedOS)
	if err != nil {
		log.Fatal(err)
	}

	zipReader, err := zip.OpenReader(operatingSystems[selectedOS].filename)
	if err != nil {
		log.Fatal(err)
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		filePath := filepath.Join(extractPath, f.Name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				log.Fatal(err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			log.Fatal(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			log.Fatal(err)
		}
		srcFile, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			log.Fatal(err)
		}

		dstFile.Close()
		srcFile.Close()
	}

	if operatingSystems[selectedOS].targetFilename == "Launcher.pdx" {
		err = os.Rename(filepath.Join(extractPath, pdxPath), filepath.Join(pdosExtractPath, "System", "Launcher.pdx"))
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	err = os.Rename(filepath.Join(extractPath, pdxPath), filepath.Join(pdosExtractPath, "System", "Launchers", operatingSystems[selectedOS].targetFilename))
	if err != nil {
		log.Fatal(err)
	}
}
