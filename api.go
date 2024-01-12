package main

import (
	"encoding/json"
	"fmt"
	"github.com/wii-tools/libtpl"
	"image/png"
	"io"
	"net/http"
)

const (
	catalogUrl = "https://hbb1.oscwii.org/api/v3/contents"
	iconUrl    = "https://hbb1.oscwii.org/api/v3/contents/%s/icon.png"
)

type App struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Slug        string `json:"slug"`
	Version     string `json:"version"`
	ReleaseDate int64  `json:"release_date"`
	Shop        struct {
		TitleID string `json:"title_id"`
		Version int    `json:"title_version"`
	} `json:"shop"`
	Description struct {
		Long  string `json:"long"`
		Short string `json:"short"`
	} `json:"description"`
}

func GetCatalog() ([]App, error) {
	response, err := http.Get(catalogUrl)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)

	var apps []App
	err = json.Unmarshal(data, &apps)
	if err != nil {
		return nil, err
	}

	return apps, nil
}

func (a *App) GetIcon() ([]byte, error) {
	response, err := http.Get(fmt.Sprintf(iconUrl, a.Slug))
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	img, err := png.Decode(response.Body)
	if err != nil {
		return nil, err
	}

	return libtpl.ToRGB5A3(img)
}
