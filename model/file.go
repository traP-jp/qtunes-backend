package model

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/antihax/optional"
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

	return res, nil
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

func ToggleFileFavorite(ctx context.Context, accessToken string, userID string, fileID string, favorite bool) error {
	if favorite {
		client, auth := newClient(accessToken)
		file, res, err := client.FileApi.GetFileMeta(auth, fileID)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
		}

		if err := insertFileFavorite(ctx, userID, *file.UploaderId, fileID); err != nil {
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

//TODO後で消す
// 拡張子を除く
func format(str string) string {
	rep := regexp.MustCompile(`\.[A-Za-z0-9]{3,5}`)
	return rep.ReplaceAllString(str, "")
}

//TODO後で消す
// offsetを変えて全ファイルを取得
func getAllFiles(accessToken string) ([]traq.FileInfo, error) {
	var files []traq.FileInfo

	client, auth := newClient(accessToken)
	for i := 0; ; i += 200 {
		f, res, err := client.FileApi.GetFiles(auth, &traq.FileApiGetFilesOpts{
			ChannelId: optional.NewInterface(SoundChannelId),
			Limit:     optional.NewInt32(200),
			Offset:    optional.NewInt32(int32(i)),
		})
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			return nil, err
		}
		if len(f) == 0 {
			break
		}

		files = append(files, f...)
	}

	return files, nil
}
