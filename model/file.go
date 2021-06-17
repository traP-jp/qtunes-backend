package model

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type FileInfo struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Mime            string    `json:"mime"`
	Size            int       `json:"size"`
	Md5             string    `json:"md5"`
	IsAnimatedImage bool      `json:"isAnimatedImage"`
	CreatedAt       time.Time `json:"createAt"`
	Thumbnails      []struct {
		Type   string `json:"type"`
		Mime   string `json:"mime"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	}
	ChannelId  string `json:"channelId"`
	UpLoaderId string `json:"upLoaderId"`
}

var baseURL, _ = url.Parse("https://q.trap.jp/api/v3")

func GetFiles(ctx context.Context, accessToken string) ([]*FileInfo, error) {
	path := *baseURL
	path.Path += "/files"
	req, err := http.NewRequest("GET", path.String(), nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("channelId", "8bd9e07a-2c6a-49e6-9961-4f88e83b4918")
	req.URL.RawQuery = params.Encode()

	req.Header.Set("Authorization", "Bearer "+accessToken)
	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	var files []*FileInfo
	err = json.NewDecoder(res.Body).Decode(&files)
	if err != nil {
		return nil, err
	}
	
	var audioFiles []*FileInfo
	var fileElement FileInfo
	for i := 0; i < len(files); i++ {
		fileElement = *files[i]
		if strings.HasPrefix(fileElement.Mime, "audio") {
			audioFiles = append(audioFiles, files[i])
		}

	}

	return audioFiles, nil
}
