package model

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	traq "github.com/sapphi-red/go-traq"
)

func newClient(accessToken string) (*traq.APIClient, context.Context) {
	client := traq.NewAPIClient(traq.NewConfiguration())
	auth := context.WithValue(context.Background(), traq.ContextAccessToken, accessToken)

	return client, auth
}

func getThumbnailImage(accessToken string, fileID string) (*http.Response, error) {
	path := *BaseUrl
	path.Path += fmt.Sprintf("/files/%s/thumbnail", fileID)
	form := url.Values{}
	form.Set("type", "waveform")
	reqBody := strings.NewReader(form.Encode())
	req, err := http.NewRequest("GET", path.String(), reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	return res, nil
}
