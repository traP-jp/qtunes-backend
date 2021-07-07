package model

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	traq "github.com/sapphi-red/go-traq"
)

type File struct {
	ID             string    `json:"id"  db:"id"`
	Title          string    `json:"title"  db:"title"`
	ComposerID     string    `json:"composer_id"  db:"composer_id"`
	ComposerName   string    `json:"composer_name"  db:"composer_name"`
	MessageID      string    `json:"message_id"  db:"message_id"`
	FavoriteCount  int       `json:"favorite_count"  db:"-"`
	IsFavoriteByMe bool      `json:"is_favorite_by_me"  db:"-"`
	CreatedAt      time.Time `json:"created_at"  db:"created_at"`
}

func GetFiles(ctx context.Context, userID string) ([]*File, error) {
	var files []*File
	err := db.SelectContext(ctx, &files, "SELECT * FROM files ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("Failed to get files: %w", err)
	}

	// DBからお気に入りを取得
	favCountMap, err := getFavoriteCounts(ctx)
	if err != nil {
		return nil, err
	}
	// DBから自分がお気に入りに追加しているかを取得
	myFavMap, err := getMyFavoritesMap(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, v := range files {
		if count, ok := favCountMap[v.ID]; ok {
			v.FavoriteCount = count
		} else {
			v.FavoriteCount = 0
		}

		if _, ok := myFavMap[v.ID]; ok {
			v.IsFavoriteByMe = true
		} else {
			v.IsFavoriteByMe = false
		}
	}

	return files, nil
}

func GetRandomFile(ctx context.Context, userID string) (*File, error) {
	var count int
	err := db.GetContext(ctx, &count, "SELECT COUNT(id) AS count FROM files LIMIT 1")
	if err != nil {
		return nil, err
	}

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(count)

	var file File
	err = db.GetContext(ctx, &file, "SELECT * FROM files LIMIT 1 OFFSET ?", r)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func GetFile(ctx context.Context, userID, fileID string) (*File, error) {
	var file File
	err := db.GetContext(ctx, &file, "SELECT * FROM files WHERE id = ? LIMIT 1", fileID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get files: %w", err)
	}

	// DBからお気に入りを取得
	favCount, err := getFavoriteCount(ctx, fileID)
	if err != nil {
		return nil, err
	}
	// DBから自分がお気に入りに追加しているかを取得
	isFavoriteByMe, err := getMyFavorite(ctx, userID, fileID)
	if err != nil {
		return nil, err
	}

	file.FavoriteCount = favCount
	file.IsFavoriteByMe = isFavoriteByMe

	return &file, nil
}

func GetFileDownload(ctx context.Context, fileID string, accessToken string) (*os.File, *http.Response, error) {
	client, auth := newClient(accessToken)
	file, res, err := client.FileApi.GetFile(auth, fileID, &traq.FileApiGetFileOpts{})
	if err != nil {
		return nil, nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, res, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	return file, res, nil
}

func ToggleFileFavorite(ctx context.Context, userID string, fileID string, favorite bool) error {
	if favorite {
		var file File
		err := db.GetContext(ctx, &file, "SELECT * FROM files WHERE id = ? LIMIT 1", fileID)
		if err != nil {
			return fmt.Errorf("Failed to get files: %w", err)
		}

		if err := insertFileFavorite(ctx, userID, file.ComposerID, fileID); err != nil {
			return err
		}
	} else {
		if err := deleteFileFavorite(ctx, userID, fileID); err != nil {
			return err
		}
	}

	return nil
}
