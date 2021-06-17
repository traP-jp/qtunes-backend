package model

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hackathon-21-spring-02/back-end/domain"
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

func GetFiles(ctx context.Context, accessToken string, userID string) ([]*domain.File, error) {
	path := *baseURL
	path.Path += "/files"
	req, err := http.NewRequest("GET", path.String(), nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("channelId", "8bd9e07a-2c6a-49e6-9961-4f88e83b4918")
	params.Add("limit", "200")
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

	// DBからお気に入りを取得
	favoriteCounts, err := getFavoriteCounts(ctx)
	if err != nil {
		return nil, err
	}
	// DBから自分がお気に入りに追加しているかを取得
	myFavorites, err := getMyFavorites(ctx, userID)
	if err != nil {
		return nil, err
	}

	audioFiles := make([]*domain.File, 0, len(files))
	for _, v := range files {
		if strings.HasPrefix(v.Mime, "audio") {
			audioFiles = append(audioFiles, &domain.File{
				ID:             v.ID,
				ComposerID:     v.UpLoaderId,
				FavoriteCount:  favoriteCounts[v.UpLoaderId],
				IsFavoriteByMe: myFavorites[v.ID],
				CreatedAt:      v.CreatedAt,
			})
		}
	}

	return audioFiles, nil
}

func GetFileDownload(ctx context.Context, fileID string, accessToken string) (*http.Response, error) {
	path := *baseURL
	path.Path += "/files/" + fileID
	req, err := http.NewRequest("GET", path.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, err
	}

	return res, nil
}
