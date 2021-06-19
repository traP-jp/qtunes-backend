package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hackathon-21-spring-02/back-end/domain"
	"io/ioutil"
	"net/http"
	"strings"
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
	/*for _, v := range data {
		composers = append(composers, &domain.Composer{
			ID:        v.UploaderID,
			Name:      v.Name,
			PostCount: 0,
		})
	}*/
	users,err:=GetUsers(ctx)
	if err != nil {
		return nil, err
	}

	uuidToName:=make(map[string]string,len(data))
	for _,v:=range users{
		uuidToName[v.ID]=v.Name
	}
	for _,user:=range users{
		/*postByUser,err:=GetPostCounter(ctx,accessToken,user.ID)
		if err != nil {
			return nil, fmt.Errorf("Faied in Getting PostCount %v",err)
		}*/
		for _,composer:=range data{
			count,_:=GetPostCounter2(ctx,accessToken,user.ID)
			composers=append(composers,&domain.Composer{
				ID:        composer.ID,
				Name:      uuidToName[user.ID],
				PostCount: *count,
			})
		}
	}
	return composers, err
}

func GetPostCounter(ctx context.Context, accessToken string, userID string)(*map[string]int,error) {
	path := *baseURL
	path.Path += "/files"
	req, err := http.NewRequest("GET", path.String(), nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("channelId", SoundChannelId)
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
		return nil,err
	}
	postCountByUser:=make(map[string]int,0)
	for _,v:=range files{
		if strings.HasPrefix(v.Mime, "audio") {
			if _, ok := postCountByUser[v.UpLoaderId]; ok && v.UpLoaderId==userID {
				postCountByUser[v.UpLoaderId]++
			}else if _, ok := postCountByUser[v.UpLoaderId]; ok && v.UpLoaderId != userID {
				postCountByUser[v.UpLoaderId]=1
			}
		}
	}

	return &postCountByUser,err
}
func GetPostCounter2(ctx context.Context, accessToken string, userID string)(*int,error) {
	path := *baseURL
	path.Path += "/files"
	req, err := http.NewRequest("GET", path.String(), nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("channelId", SoundChannelId)
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
		return nil,err
	}
	count:=0
	for _,v:=range files{
		if strings.HasPrefix(v.Mime, "audio") {
			if v.UpLoaderId==userID {
				count++
			}
		}
	}

	return &count,err
}

