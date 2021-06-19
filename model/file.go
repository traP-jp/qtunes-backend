package model

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/antihax/optional"
	"github.com/hackathon-21-spring-02/back-end/domain"
	traq "github.com/sapphi-red/go-traq"
)

func GetFiles(ctx context.Context, accessToken string, userID string) ([]*domain.File, error) {
	client, auth := newClient(accessToken)
	files, res, err := client.FileApi.GetFiles(auth, &traq.FileApiGetFilesOpts{
		ChannelId: optional.NewInterface(SoundChannelId),
		Limit:     optional.NewInt32(200),
	})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

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
				Title:          v.Name,
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
	path := *baseURL
	path.Path += "/files"
	req, err := http.NewRequest("GET", path.String(), nil)
	if err != nil {
		return nil, err
	}
	params := req.URL.Query()
	params.Add("channelId", "8bd9e07a-2c6a-49e6-9961-4f88e83b4918") // TODO:あとでSoundChannelIDに変える
	params.Add("limit", "200")
	req.URL.RawQuery = params.Encode()

	req.Header.Set("Authorization", "Bearer "+accessToken)
	httpClient := http.DefaultClient
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	var files []*FileInfo
	err = json.NewDecoder(res.Body).Decode(&files)
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
		if strings.HasPrefix(f.Mime, "audio") {
			audioFile.ID = f.ID
			audioFile.ComposerID = f.UpLoaderId
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
	client, auth := newClient(accessToken)
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
		Title:          file.Name,
		ComposerID:     *file.UploaderId,
		ComposerName:   user.Name,
		FavoriteCount:  favoriteCount.Count,
		IsFavoriteByMe: isFavoriteByMe,
		CreatedAt:      file.CreatedAt,
	}

	return audioFile, nil
}

func GetFileDownload(ctx context.Context, fileID string, accessToken string) (*http.Response, error) {
	client, auth := newClient(accessToken)
	_, res, err := client.FileApi.GetFile(auth, fileID, &traq.FileApiGetFileOpts{})
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
	}

	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, err
	}

	return res, nil
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
