package main

import (
	"archive/zip"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.bug.st/serial"
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
		panic(err)
	}
	req.Header.Set("Idempotency-Key", indem)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var res map[string]any
	err = json.Unmarshal(body, &res)
	if err != nil {
		panic(err)
	}

	return res
}

var accessToken string = ""

func (a *App) FinishRegistration(serialNumber string) map[string]any {
	resp, err := http.Get("https://play.date/api/v2/device/register/" + serialNumber + "/complete/get")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var res map[string]any
	err = json.Unmarshal(body, &res)
	if err != nil {
		panic(err)
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
		panic(err)
	}
	req.Header.Set("Authorization", "Token "+accessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var res map[string]string
	err = json.Unmarshal(body, &res)
	if err != nil {
		panic(err)
	}

	decryptKey, err := base64.StdEncoding.DecodeString(res["decryption_key"])
	if err != nil {
		panic(err)
	}

	f, err := os.CreateTemp("", "PlaydateOS.*.pdos")
	pdosFilename = f.Name()
	if err != nil {
		panic(err)
	}
	defer f.Close()
	resp, err = http.Get(res["url"])
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		panic(err)
	}

	f, err = os.CreateTemp("", "PlaydateOS.*.pdkey")
	pdkeyFilename = f.Name()
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.Write(decryptKey)
	if err != nil {
		panic(err)
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
		panic(err)
	}

	defer f.Close()
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		panic(err)
	}

	operatingSystems[selectedOS] = OSInfo{f.Name(), targetFilename, url}
}

var pdosExtractPath string = ""

func (a *App) ExtractPlaydateOS(funnyloader bool) {
	extractPath, err := os.MkdirTemp("", "PlaydateOS.*")
	if err != nil {
		panic(err)
	}

	zipReader, err := zip.OpenReader(pdosFilename)
	if err != nil {
		panic(err)
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		filePath := filepath.Join(extractPath, f.Name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				panic(err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}
		srcFile, err := f.Open()
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			panic(err)
		}

		dstFile.Close()
		srcFile.Close()
	}

	if funnyloader {
		err = os.Mkdir(filepath.Join(extractPath, "System", "Launchers"), os.ModePerm)
		if err != nil {
			panic(err)
		}
		err = os.Rename(filepath.Join(extractPath, "System", "Launcher.pdx"), filepath.Join(extractPath, "System", "Launchers", "StockLauncher.pdx"))
		if err != nil {
			panic(err)
		}
	} else {
		err = os.Rename(filepath.Join(extractPath, "System", "Launcher.pdx"), filepath.Join(extractPath, "System", "StockLauncher.pdx"))
		if err != nil {
			panic(err)
		}
	}

	pdosExtractPath = extractPath
}

func (a *App) ExtractOS(selectedOS string, pdxPath string) {
	extractPath, err := os.MkdirTemp("", selectedOS)
	if err != nil {
		panic(err)
	}

	zipReader, err := zip.OpenReader(operatingSystems[selectedOS].filename)
	if err != nil {
		panic(err)
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		filePath := filepath.Join(extractPath, f.Name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				panic(err)
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}
		srcFile, err := f.Open()
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(dstFile, srcFile); err != nil {
			panic(err)
		}

		dstFile.Close()
		srcFile.Close()
	}

	if operatingSystems[selectedOS].targetFilename == "Launcher.pdx" {
		err = os.Rename(filepath.Join(extractPath, pdxPath), filepath.Join(pdosExtractPath, "System", "Launcher.pdx"))
		if err != nil {
			panic(err)
		}
		return
	}
	err = os.Rename(filepath.Join(extractPath, pdxPath), filepath.Join(pdosExtractPath, "System", "Launchers", operatingSystems[selectedOS].targetFilename))
	if err != nil {
		panic(err)
	}
}

var pdosPatchedPath string = ""

func (a *App) CompressPlaydateOS() {
	zipFile, err := os.CreateTemp("", "PlaydateOS-Patched.*.pdos")
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	absSourceDir, err := filepath.Abs(pdosExtractPath)
	if err != nil {
		panic(err)
	}

	err = filepath.Walk(absSourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		relPath, err := filepath.Rel(absSourceDir, filePath)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		zipName := strings.ReplaceAll(relPath, "\\", "/")
		if info.IsDir() {
			zipName += "/"
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = zipName
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
	if err != nil {
		panic(err)
	}

	pdosPatchedPath = zipFile.Name()
}

func (a *App) GetSerialPorts() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		panic(err)
	}
	return ports
}

func (a *App) UploadPatchedPlaydateOS(selectedPort string) {
	port, err := serial.Open(selectedPort, &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
	})
	if err != nil {
		panic(err)
	}
	defer port.Close()

	_, err = port.Write([]byte("\ndatadisk\n"))
	if err != nil {
		panic(err)
	}

	time.Sleep(8 * time.Second)

	info, err := FindMount("PLAYDATE")
	if err != nil {
		panic(err)
	}

	i, err := os.Open(pdosPatchedPath)
	if err != nil {
		panic(err)
	}
	defer i.Close()

	o, err := os.Create(filepath.Join(info.MountPoint, "PlaydateOS-Patched.pdos"))
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(o, i)
	if err != nil {
		panic(err)
	}

	if err = o.Sync(); err != nil {
		panic(err)
	}
	if err = o.Close(); err != nil {
		panic(err)
	}

	i, err = os.Open(pdkeyFilename)
	if err != nil {
		panic(err)
	}
	defer i.Close()

	o, err = os.Create(filepath.Join(info.MountPoint, "PlaydateOS.pdkey"))
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(o, i)
	if err != nil {
		panic(err)
	}

	if err = o.Sync(); err != nil {
		panic(err)
	}
	if err = o.Close(); err != nil {
		panic(err)
	}

	err = UnmountAndEject(info)
	if err != nil {
		panic(err)
	}
}

func (a *App) InstallPatchedPlaydateOS(selectedPort string) {
	port, err := serial.Open(selectedPort, &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	})
	if err != nil {
		panic(err)
	}
	defer port.Close()

	_, err = port.Write([]byte("\nfwup /PlaydateOS-Patched.pdos /PlaydateOS.pdkey\n"))
	if err != nil {
		panic(err)
	}
}

func (a *App) CleanUp(selectedPort string) {
	port, err := serial.Open(selectedPort, &serial.Mode{
		BaudRate: 115200,
		DataBits: 8,
		Parity:   serial.NoParity,
	})
	if err != nil {
		panic(err)
	}
	defer port.Close()

	_, err = port.Write([]byte("\ndatadisk\n"))
	if err != nil {
		panic(err)
	}

	time.Sleep(8 * time.Second)

	info, err := FindMount("PLAYDATE")
	if err != nil {
		panic(err)
	}

	os.Remove(filepath.Join(info.MountPoint, "PlaydateOS-Patched.pdos"))
	os.Remove(filepath.Join(info.MountPoint, "PlaydateOS.pdkey"))

	err = UnmountAndEject(info)
	if err != nil {
		panic(err)
	}
}
