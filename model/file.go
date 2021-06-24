package model

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hackathon-21-spring-02/back-end/domain"
	traq "github.com/sapphi-red/go-traq"
)

type File struct {
	ID           string    `db:"id"`
	Title        string    `db:"title"`
	ComposerID   string    `db:"composer_id"`
	ComposerName string    `db:"composer_name"`
	MessageID    string    `db:"message_id"`
	CreatedAt    time.Time `db:"created_at"`
}

func GetFiles(ctx context.Context, userID string) ([]*domain.File, error) {
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

	res := make([]*domain.File, 0, len(files))
	for _, v := range files {
		f := convertFile(*v, favoriteCounts[v.ID], myFavorites[v.ID])
		res = append(res, &f)
	}

	return res, nil //TODO:domain.Fileを返すべきではない
}

func GetFile(ctx context.Context, userID, fileID string) (*domain.File, error) {
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

	res := convertFile(file, favoriteCount.Count, isFavoriteByMe)

	return &res, nil
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
		err := db.GetContext(ctx, &file, "SELECT * FROM files LIMIT 1")
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

func convertFile(file File, count uint32, isFavorite bool) domain.File {
	return domain.File{
		ID:             file.ID,
		Title:          file.Title,
		ComposerID:     file.ComposerID,
		ComposerName:   file.ComposerName,
		MessageID:      file.MessageID,
		FavoriteCount:  count,
		IsFavoriteByMe: isFavorite,
		CreatedAt:      file.CreatedAt,
	}
}
