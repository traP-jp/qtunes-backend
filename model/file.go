package model

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/antihax/optional"
	"github.com/hackathon-21-spring-02/back-end/domain"
	"github.com/jmoiron/sqlx"
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

func GetFiles(ctx context.Context, accessToken string, userID string) ([]*domain.File, error) {
	files, err := getAllFiles(accessToken)
	if err != nil {
		return nil, err
	}

	client, auth := NewTraqClient(accessToken)
	users, res, err := client.UserApi.GetUsers(auth, &traq.UserApiGetUsersOpts{IncludeSuspended: optional.NewBool(true)})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	userIdMap := map[string]string{}
	for _, v := range users {
		userIdMap[v.Id] = v.Name
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

	audioFiles := make([]*domain.File, 0, len(files))
	for _, v := range files {
		if strings.HasPrefix(v.Mime, "audio") {
			audioFiles = append(audioFiles, &domain.File{
				ID:             v.Id,
				Title:          format(v.Name),
				ComposerID:     *v.UploaderId,
				ComposerName:   userIdMap[*v.UploaderId],
				FavoriteCount:  favoriteCounts[v.Id],
				IsFavoriteByMe: myFavorites[v.Id],
				CreatedAt:      v.CreatedAt,
			})
		}
	}

	return audioFiles, nil
}

func GetRandomFile(ctx context.Context, accessToken string, userID string) (*domain.File, error) {
	files, err := getAllFiles(accessToken)
	if err != nil {
		return nil, err
	}

	audioFile := &domain.File{}
	for {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(len(files))
		if err != nil {
			return nil, fmt.Errorf("failed to generate random integer: %w", err)
		}
		f := files[r]

		client, auth := NewTraqClient(accessToken)
		user, res, err := client.UserApi.GetUser(auth, *f.UploaderId)
		if err != nil {
			return nil, err
		}
		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
		}
		if strings.HasPrefix(f.Mime, "audio") {
			audioFile.ID = f.Id
			audioFile.Title = format(f.Name)
			audioFile.ComposerName = user.Name
			audioFile.ComposerID = *f.UploaderId
			audioFile.CreatedAt = f.CreatedAt
			break
		}
	}

	// DBからお気に入りを取得
	favoriteCount, err := getFavoriteCount(ctx, audioFile.ID)
	if err != nil {
		return nil, err
	}
	// DBから自分がお気に入りに追加しているかを取得
	isFavoriteByMe, err := getMyFavorite(ctx, userID, audioFile.ID)
	if err != nil {
		return nil, err
	}

	audioFile.FavoriteCount = favoriteCount.Count
	audioFile.IsFavoriteByMe = isFavoriteByMe

	return audioFile, nil
}

func GetFile(ctx context.Context, accessToken string, userID, fileID string) (*domain.File, error) {
	client, auth := NewTraqClient(accessToken)
	file, res, err := client.FileApi.GetFileMeta(auth, fileID)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	if !strings.HasPrefix(file.Mime, "audio") {
		return nil, fmt.Errorf("")
	}

	user, res, err := client.UserApi.GetUser(auth, *file.UploaderId)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
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

	audioFile := &domain.File{
		ID:             file.Id,
		Title:          format(file.Name),
		ComposerID:     *file.UploaderId,
		ComposerName:   user.Name,
		FavoriteCount:  favoriteCount.Count,
		IsFavoriteByMe: isFavoriteByMe,
		CreatedAt:      file.CreatedAt,
	}

	return audioFile, nil
}

func GetFileDownload(ctx context.Context, fileID string, accessToken string) (*os.File, *http.Response, error) {
	client, auth := NewTraqClient(accessToken)
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
		client, auth := NewTraqClient(accessToken)
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

func GetFileIDsInMessage(ctx context.Context, messageID string) ([]string, error) {
	var files []string
	err := db.SelectContext(ctx, &files, "SELECT id FROM files WHERE message_id = ?", messageID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get files: %w", err)
	}
	return files, nil
}

func InsertFile(ctx context.Context, file *File) error {
	_, err := db.ExecContext(
		ctx,
		"INSERT IGNORE INTO files (id, title, composer_id, composer_name, message_id, created_at)VALUES (?, ?, ?, ?, ?, ?)",
		file.ID, file.Title, file.ComposerID, file.ComposerName, file.MessageID, file.CreatedAt,
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

func DeleteFilesFromMessageId(ctx context.Context, messageID string) error {
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

// 拡張子を除く
func format(str string) string {
	rep := regexp.MustCompile(`\.[A-Za-z0-9]{3,5}`)
	return rep.ReplaceAllString(str, "")
}

// offsetを変えて全ファイルを取得
func getAllFiles(accessToken string) ([]traq.FileInfo, error) {
	var files []traq.FileInfo

	client, auth := NewTraqClient(accessToken)
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
