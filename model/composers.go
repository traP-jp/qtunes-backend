package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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
func GetComposers() (map[string]string,error){
	v:=url.Values{}

	v.Set("channelId",soundChannelId)
	reqBody:=strings.NewReader(v.Encode())
	path:=*BaseUrl
	path.Path+="/files"
	req,err:=http.NewRequest("GET",path.String(),reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type","application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil,err
	}

	if res.StatusCode != 200 {
		return  nil,fmt.Errorf("Failed In Getting Information:(Status:%d %s)", res.StatusCode, res.Status)
	}
	data:=make([]*ComposersInfo,0)

	body,err:=ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	var result =make(map[string]string,0)
	for _,r:=range data{
		result[r.UploaderID]=r.Name
	}
	return result,err


}

