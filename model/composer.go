package model

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/antihax/optional"
	"github.com/sapphi-red/go-traq"
)

type Composer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	PostCount int       `json:"post_count"`
	UpdatedAt time.Time `json:"-"`
}

func GetComposers(ctx context.Context, accessToken string) ([]*Composer, error) {
	client, auth := NewTraqClient(accessToken)
	users, res, err := client.UserApi.GetUsers(auth, &traq.UserApiGetUsersOpts{IncludeSuspended: optional.NewBool(true)})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	composersMap, err := getComposersMap(ctx)
	if err != nil {
		return nil, err
	}

	composers := make([]*Composer, 0, len(users))
	for _, user := range users {
		if val, ok := composersMap[user.Id]; ok && val.PostCount > 0 {
			composers = append(composers, &Composer{
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

func GetComposer(ctx context.Context, accessToken string, composerID string) (*Composer, error) {
	client, auth := NewTraqClient(accessToken)
	user, res, err := client.UserApi.GetUser(auth, composerID)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	postCountByUser, err := getComposersMap(ctx)
	if err != nil {
		return nil, err
	}

	composer := &Composer{
		ID:        composerID,
		Name:      user.Name,
		PostCount: postCountByUser[user.Id].PostCount,
		UpdatedAt: postCountByUser[user.Id].UpdatedAt,
	}
	return composer, err
}

func GetComposerByName(ctx context.Context, accessToken string, name string) (*Composer, error) {
	client, auth := NewTraqClient(accessToken)
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

	postCountByUser, err := getComposersMap(ctx)
	if err != nil {
		return nil, err
	}

	u := users[0]
	composer := Composer{
		ID:        u.Id,
		Name:      u.Name,
		PostCount: postCountByUser[u.Id].PostCount,
		UpdatedAt: postCountByUser[u.Id].UpdatedAt,
	}

	return &composer, nil
}

func GetComposerFiles(ctx context.Context, accessToken string, composerID string, userID string) ([]*File, error) {
	client, auth := NewTraqClient(accessToken)
	user, res, err := client.UserApi.GetUser(auth, composerID)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	files, err := GetFiles(ctx, composerID)
	if err != nil {
		return nil, err
	}

	myFavMap, err := getMyFavoritesMap(ctx, userID)
	if err != nil {
		return nil, err
	}

	composerFiles := make([]*File, 0, len(files))
	for _, file := range files {
		if file.ComposerID == composerID {
			f := &File{
				ID:            file.ID,
				Title:         file.Title,
				ComposerID:    composerID,
				ComposerName:  user.Name,
				FavoriteCount: file.FavoriteCount,
				CreatedAt:     file.CreatedAt,
			}
			if _, ok := myFavMap[file.ID]; ok {
				f.IsFavoriteByMe = true
			} else {
				f.IsFavoriteByMe = false
			}
			composerFiles = append(composerFiles, f)
		}
	}

	return composerFiles, nil
}

func getComposersMap(ctx context.Context) (map[string]*Composer, error) {
	composersMap := make(map[string]*Composer)
	var files []*File
	err := db.SelectContext(ctx, &files, "SELECT composer_id, created_at FROM files")
	if err != nil {
		return nil, fmt.Errorf("Failed to get files: %w", err)
	}

	for _, v := range files {
		if _, ok := composersMap[v.ComposerID]; !ok {
			composersMap[v.ComposerID] = &Composer{
				PostCount: 1,
				UpdatedAt: v.CreatedAt,
			}
			continue
		}

		composersMap[v.ComposerID].PostCount++
		if v.CreatedAt.After(composersMap[v.ComposerID].UpdatedAt) {
			composersMap[v.ComposerID].UpdatedAt = v.CreatedAt
		}
	}

	return composersMap, nil
}
