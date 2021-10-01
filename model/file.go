package model

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
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
	if err == sql.ErrNoRows {
		return []*File{}, nil
	}
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
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func GetFile(ctx context.Context, userID, fileID string) (*File, error) {
	var file File
	err := db.GetContext(ctx, &file, "SELECT * FROM files WHERE id = ? LIMIT 1", fileID)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
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

func GetFileDownload(ctx context.Context, fileID, accessToken string) (*os.File, *http.Response, error) {
	traqapi := NewTraqAPI(accessToken)
	file, res, err := traqapi.GetFile(fileID, &GetFileOpts{})
	if err != nil {
		return nil, nil, err
	}

	return file, res, nil
}

func ToggleFileFavorite(ctx context.Context, userID, fileID string, isFavorite bool) error {
	opts := FavoriteOptions{
		UserID:  userID,
		SoundID: fileID,
	}

	if isFavorite {
		var composerID string
		err := db.GetContext(ctx, &composerID, "SELECT composer_id FROM files WHERE id = ? LIMIT 1", fileID)
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("Failed to get files: %w", err)
		}

		opts.ComposerID = composerID
		if err := insertFileFavorite(ctx, opts); err != nil {
			return err
		}
	} else {
		if err := deleteFileFavorite(ctx, opts); err != nil {
			return err
		}
	}

	return nil
}

func GetFileIDsInMessage(ctx context.Context, messageID string) ([]string, error) {
	var files []string
	err := db.SelectContext(ctx, &files, "SELECT id FROM files WHERE message_id = ?", messageID)
	if err == sql.ErrNoRows {
		return []string{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to get files: %w", err)
	}
	return files, nil
}

func InsertFiles(ctx context.Context, files []*File) error {
	_, err := db.NamedExecContext(
		ctx,
		`INSERT IGNORE INTO files (id, title, composer_id, composer_name, message_id, created_at)
		VALUES (:id, :title, :composer_id, :composer_name, :message_id, :created_at)`,
		files,
	)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFiles(ctx context.Context, fileIDs []string) error {
	query := "DELETE favorites, files FROM files LEFT JOIN favorites ON favorites.sound_id = files.id WHERE files.id IN (?)"
	query, params, err := sqlx.In(query, fileIDs)
	if err != nil {
		return fmt.Errorf("Failed to delete files: %w", err)
	}

	_, err = db.ExecContext(ctx, query, params...)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFilesFromMessageID(ctx context.Context, messageID string) error {
	_, err := db.ExecContext(
		ctx,
		"DELETE favorites, files FROM files LEFT JOIN favorites ON favorites.sound_id = files.id WHERE files.message_id = ?",
		messageID,
	)
	if err != nil {
		return err
	}

	return nil
}

func FindFileFromComposerName(ctx context.Context, composerName string) ([]*File, error) {
	var file []*File
	likeComposer := "%" + composerName + "%"
	err := db.SelectContext(ctx, &file, "SELECT * FROM files WHERE composer_name LIKE ? ORDER BY CASE WHEN composer_name = ? THEN 0 WHEN composer_name LIKE ? THEN 1 ELSE 2 END", likeComposer, composerName, likeComposer[1:])
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to get composerName: %w", err)
	}
	return file, nil
}

func FindFileFromTitle(ctx context.Context, songTitle string) ([]*File, error) {
	var file []*File
	likeTitle := "%" + songTitle + "%"
	err := db.SelectContext(ctx, &file, "SELECT * FROM files WHERE title LIKE ? ORDER BY CASE WHEN title = ? THEN 0 WHEN title LIKE ? THEN 1 ELSE 2 END", likeTitle, songTitle, likeTitle[1:])
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("Failed to get title: %w", err)
	}
	return file, nil
}
