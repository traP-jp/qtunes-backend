package model

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/antihax/optional"
	"github.com/hackathon-21-spring-02/back-end/domain"
	"github.com/sapphi-red/go-traq"
)

func GetComposers(ctx context.Context, accessToken string) ([]*domain.Composer, error) {
	client, auth := newClient(accessToken)
	users, res, err := client.UserApi.GetUsers(auth, &traq.UserApiGetUsersOpts{IncludeSuspended: optional.NewBool(true)})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	postCountByUser, err := caluculateCount(accessToken)
	if err != nil {
		return nil, err
	}

	composers := make([]*domain.Composer, 0, len(users))
	for _, user := range users {
		if val, ok := postCountByUser[user.Id]; ok && val > 0 {
			composers = append(composers, &domain.Composer{
				ID:        user.Id,
				Name:      user.Name,
				PostCount: val,
			})
		}
	}
	return composers, err
}

func GetComposer(ctx context.Context, accessToken string, composerID string) (*domain.Composer, error) {
	client, auth := newClient(accessToken)
	user, res, err := client.UserApi.GetUser(auth, composerID)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	postCountByUser, err := caluculateCount(accessToken)
	if err != nil {
		return nil, err
	}

	composer := &domain.Composer{
		ID:        composerID,
		Name:      user.Name,
		PostCount: postCountByUser[user.Id],
	}
	return composer, err
}

func GetComposerByName(ctx context.Context, accessToken string, name string) (*domain.Composer, error) {
	client, auth := newClient(accessToken)
	users, res, err := client.UserApi.GetUsers(auth, &traq.UserApiGetUsersOpts{Name: optional.NewString(name)})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}
	if len(users) != 1 {
		return nil, fmt.Errorf("Invalid name")
	}

	postCountByUser, err := caluculateCount(accessToken)
	if err != nil {
		return nil, err
	}

	u := users[0]
	composer := domain.Composer{
		ID:        u.Id,
		Name:      u.Name,
		PostCount: postCountByUser[u.Id],
	}

	return &composer, nil
}

func GetComposerFiles(ctx context.Context, accessToken string, composerID string, userID string) ([]*domain.File, error) {
	client, auth := newClient(accessToken)
	user, res, err := client.UserApi.GetUser(auth, composerID)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	files, err := GetFiles(ctx, accessToken, composerID)
	if err != nil {
		return nil, err
	}

	getMyFavorites, err := getMyFavorites(ctx, userID)
	if err != nil {
		return nil, err
	}

	composerFiles := make([]*domain.File, 0, len(files))
	for _, file := range files {
		if file.ComposerID == composerID {
			composerFiles = append(composerFiles, &domain.File{
				ID:             file.ID,
				Title:          format(file.Title),
				ComposerID:     composerID,
				ComposerName:   user.Name,
				FavoriteCount:  file.FavoriteCount,
				IsFavoriteByMe: getMyFavorites[file.ID],
				CreatedAt:      file.CreatedAt,
			})
		}
	}

	return composerFiles, nil
}

func caluculateCount(accessToken string) (map[string]int, error) {
	cnt := make(map[string]int)
	files, res, err := getAllFiles(accessToken)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	for _, v := range files {
		if strings.HasPrefix(v.Mime, "audio") {
			cnt[*v.UploaderId]++
		}
	}

	return cnt, nil
}

func getAllFiles(accessToken string) ([]traq.FileInfo, *http.Response, error) {
	var (
		files []traq.FileInfo
		res   *http.Response
	)

	client, auth := newClient(accessToken)
	for i := 0; ; i += 200 {
		f, res, err := client.FileApi.GetFiles(auth, &traq.FileApiGetFilesOpts{
			ChannelId: optional.NewInterface(SoundChannelId),
			Limit:     optional.NewInt32(200),
			Offset:    optional.NewInt32(int32(i)),
		})
		if err != nil {
			return nil, res, err
		}
		if res.StatusCode != http.StatusOK {
			return nil, res, nil
		}
		if len(files) == 0 {
			break
		}

		files = append(files, f...)
	}

	return files, res, nil
}
