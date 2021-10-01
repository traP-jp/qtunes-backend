package model

import (
	"context"
	"fmt"
	"net/http"
	"os"

	traq "github.com/sapphi-red/go-traq"
)

type GetUsersOpts traq.UserApiGetUsersOpts
type GetFileOpts traq.FileApiGetFileOpts

type TraqAPI interface {
	GetUsers(opts *GetUsersOpts) ([]traq.User, error)
	GetUser(id string) (*traq.UserDetail, error)
	GetFile(id string, opts *GetFileOpts) (*os.File, *http.Response, error)
	GetFileMeta(id string) (*traq.FileInfo, error)
}

type traqAPI struct {
	client *traq.APIClient
	auth   context.Context
}

func NewTraqAPI(accessToken string) TraqAPI {
	client := traq.NewAPIClient(traq.NewConfiguration())
	auth := context.WithValue(context.Background(), traq.ContextAccessToken, accessToken)

	return &traqAPI{client, auth}
}

func (a *traqAPI) GetUsers(opts *GetUsersOpts) ([]traq.User, error) {
	users, res, err := a.client.UserApi.GetUsers(a.auth, &traq.UserApiGetUsersOpts{
		IncludeSuspended: opts.IncludeSuspended,
		Name:             opts.Name,
	})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	return users, nil
}

func (a *traqAPI) GetUser(id string) (*traq.UserDetail, error) {
	user, res, err := a.client.UserApi.GetUser(a.auth, id)
	if res.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (a *traqAPI) GetFile(id string, opts *GetFileOpts) (*os.File, *http.Response, error) {
	file, res, err := a.client.FileApi.GetFile(a.auth, id, &traq.FileApiGetFileOpts{
		Dl: opts.Dl,
	})
	if err != nil {
		return nil, res, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, res, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	return file, res, nil
}

func (a *traqAPI) GetFileMeta(id string) (*traq.FileInfo, error) {
	file, res, err := a.client.FileApi.GetFileMeta(a.auth, id)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	return &file, nil
}
