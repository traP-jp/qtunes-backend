package model

import (
	"context"
	"fmt"
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
	err := db.SelectContext(ctx, &files, "SELECT * FROM files")
	if err != nil {
		return nil, fmt.Errorf("Failed to get files: %w", err)
	}

	// DBからお気に入りを取得
	favoriteCounts, err := getFavoriteCounts(ctx)
	if err != nil {
		return nil, err
	}
	// DBから自分がお気に入りに追加しているかを取得
	myFavorites, err := getMyFavorites(ctx, userID)
	if err != nil {
		return nil, err
	}

	for _, v := range files {
		if cnt, ok := favoriteCounts[v.ID]; ok {
			v.FavoriteCount = cnt
		} else {
			v.FavoriteCount = 0
		}

		if flag, ok := myFavorites[v.ID]; ok {
			v.IsFavoriteByMe = flag
		} else {
			v.IsFavoriteByMe = false
		}
	}

	return files, nil
}

func GetRandomFile(ctx context.Context, userID string) (*File, error) {
	var rand int
	err := db.GetContext(ctx, &rand, "SELECT FLOOR(COUNT(id)*RAND()) AS rand FROM files LIMIT 1")
	if err != nil {
		return nil, err
	}

	var file File
	err = db.GetContext(ctx, &file, "SELECT * FROM files LIMIT 1 OFFSET ?", rand)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func GetFile(ctx context.Context, userID, fileID string) (*File, error) {
	var file File
	err := db.GetContext(ctx, &file, "SELECT * FROM files LIMIT 1")
	if err != nil {
		return nil, fmt.Errorf("Failed to get files: %w", err)
	}

	// DBからお気に入りを取得
	favoriteCount, err := getFavoriteCount(ctx, fileID)
	if err != nil {
		return nil, err
	}
	// DBから自分がお気に入りに追加しているかを取得
	isFavoriteByMe, err := getMyFavorite(ctx, userID, fileID)
	if err != nil {
		return nil, err
	}

	file.FavoriteCount = favoriteCount.Count
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
