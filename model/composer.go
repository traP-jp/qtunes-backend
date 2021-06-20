package model

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/antihax/optional"
	"github.com/hackathon-21-spring-02/back-end/domain"
	"github.com/sapphi-red/go-traq"
)

// TODO: 変数名ちゃんと考える
type composerInfo struct {
	PostCount int
	UpdatedAt time.Time
}

func GetComposers(ctx context.Context, accessToken string) ([]*domain.Composer, error) {
	client, auth := newClient(accessToken)
	users, res, err := client.UserApi.GetUsers(auth, &traq.UserApiGetUsersOpts{IncludeSuspended: optional.NewBool(true)})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	info, err := getComposersInfo(accessToken)
	if err != nil {
		return nil, err
	}

	composers := make([]*domain.Composer, 0, len(users))
	for _, user := range users {
		if val, ok := info[user.Id]; ok && val.PostCount > 0 {
			composers = append(composers, &domain.Composer{
				ID:        user.Id,
				Name:      user.Name,
				PostCount: val.PostCount,
				UpdatedAt: val.UpdatedAt,
			})
		}
	}

	sort.Slice(composers, func(i, j int) bool {
		return composers[i].UpdatedAt.After(composers[j].UpdatedAt)
	})

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

	postCountByUser, err := getComposersInfo(accessToken)
	if err != nil {
		return nil, err
	}

	composer := &domain.Composer{
		ID:        composerID,
		Name:      user.Name,
		PostCount: postCountByUser[user.Id].PostCount,
		UpdatedAt: postCountByUser[user.Id].UpdatedAt,
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

	postCountByUser, err := getComposersInfo(accessToken)
	if err != nil {
		return nil, err
	}

	u := users[0]
	composer := domain.Composer{
		ID:        u.Id,
		Name:      u.Name,
		PostCount: postCountByUser[u.Id].PostCount,
		UpdatedAt: postCountByUser[u.Id].UpdatedAt,
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

// TODO: 関数名考える
func getComposersInfo(accessToken string) (map[string]*composerInfo, error) {
	info := make(map[string]*composerInfo)
	files, err := getAllFiles(accessToken)
	if err != nil {
		return nil, err
	}

	for _, v := range files {
		if strings.HasPrefix(v.Mime, "audio") {
			if _, ok := info[*v.UploaderId]; !ok {
				info[*v.UploaderId] = &composerInfo{
					PostCount: 1,
					UpdatedAt: v.CreatedAt,
				}
				continue
			}

			info[*v.UploaderId].PostCount++
			if v.CreatedAt.After(info[*v.UploaderId].UpdatedAt) {
				info[*v.UploaderId].UpdatedAt = v.CreatedAt
			}
		}
	}

	return info, nil
}
