package model

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
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
	params.Add("channelId", "8bd9e07a-2c6a-49e6-9961-4f88e83b4918") // TODO:あとでSoundChannelIDに変える
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
				FavoriteCount:  favoriteCounts[v.ID],
				IsFavoriteByMe: myFavorites[v.ID],
				CreatedAt:      v.CreatedAt,
			})
		}
	}

	return audioFiles, nil
}

func GetRandomFile(ctx context.Context, accessToken string, userID string) (*domain.File, error) {
	path := *baseURL
	path.Path += "/files"
	req, err := http.NewRequest("GET", path.String(), nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("channelId", "8bd9e07a-2c6a-49e6-9961-4f88e83b4918") // TODO:あとでSoundChannelIDに変える
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

	audioFile := &domain.File{}
	for {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(len(files))
		if err != nil {
			return nil, fmt.Errorf("failed to generate random integer: %w", err)
		}
		f := files[r]
		if strings.HasPrefix(f.Mime, "audio") {
			audioFile.ID = f.ID
			audioFile.ComposerID = f.UpLoaderId
			audioFile.CreatedAt = f.CreatedAt
			break
		}
	}

	// DBからお気に入りを取得
	favoriteCount, err := getFavoriteCount(ctx, audioFile.ID)
	if err != nil {
		return nil, err
	}
	// DBから自分がお気に入りに追加しているかを取得
	isFavoriteByMe, err := getMyFavorite(ctx, userID, audioFile.ID)
	if err != nil {
		return nil, err
	}

	audioFile.FavoriteCount = favoriteCount.Count
	audioFile.IsFavoriteByMe = isFavoriteByMe

	return audioFile, nil
}

func GetFile(ctx context.Context, accessToken string, userID, fileID string) (*domain.File, error) {
	path := *baseURL
	path.Path += fmt.Sprintf("/files/%s/meta", fileID)
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

	file := FileInfo{}
	err = json.NewDecoder(res.Body).Decode(&file)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(file.Mime, "audio") {
		return nil, fmt.Errorf("")
	}

	// DBからお気に入りを取得
	favoriteCount, err := getFavoriteCount(ctx, fileID)
	if err != nil {
		return nil, err
	}
	// DBから自分がお気に入りに追加しているかを取得
	isFavoriteByMe, err := getMyFavorite(ctx, userID, fileID)
	if err != nil {
		return nil, err
	}

	audioFile := &domain.File{
		ID:             file.ID,
		ComposerID:     file.UpLoaderId,
		FavoriteCount:  favoriteCount.Count,
		IsFavoriteByMe: isFavoriteByMe,
		CreatedAt:      file.CreatedAt,
	}

	return audioFile, nil
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
