package model

import (
	"context"
	"net/http"
	"net/url"
)

var baseURL, _ = url.Parse("https://q.trap.jp/api/v3")

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
