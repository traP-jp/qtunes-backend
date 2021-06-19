package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hackathon-21-spring-02/back-end/domain"
	"io/ioutil"
	"net/http"
	"time"
)

type ComposersInfo struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Mime            string    `json:"mime"`
	Size            int       `json:"size"`
	Md5             string    `json:"md5"`
	IsAnimatedImage bool      `json:"isAnimatedImage"`
	CreatedAt       time.Time `json:"createdAt"`
	Thumbnails      []struct {
		Type   string `json:"type"`
		Mime   string `json:"mime"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"thumbnails"`
	ChannelID  string `json:"channelId"`
	UploaderID string `json:"uploaderId"`
}

func GetComposers(ctx context.Context, accessToken string) ([]*domain.Composer, error) {
	path := *BaseUrl
	path.Path += "/files"
	req, err := http.NewRequest("GET", path.String(), nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("channelId", SoundChannelId)
	params.Add("limit", "200")
	req.URL.RawQuery = params.Encode()

	req.Header.Set("content-type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Failed In Getting Information:(Status:%d %s)", res.StatusCode, res.Status)
	}
	data := make([]*ComposersInfo, 0)

	body, err := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	composers := make([]*domain.Composer, 0, len(data))
	for _, v := range data {
		composers = append(composers, &domain.Composer{
			ID:        v.UploaderID,
			Name:      v.Name,
			PostCount: 0, //TODO: データベースからとってくるようにする
		})
	}

	return composers, err
}
