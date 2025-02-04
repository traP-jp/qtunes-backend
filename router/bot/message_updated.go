package bot

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hackathon-21-spring-02/back-end/model"
	traqbot "github.com/traPtitech/traq-bot"
)

// MessageUpdatedHandler MessageUpdatedイベントを処理する
func MessageUpdatedHandler(ctx context.Context, accessToken string, payload *traqbot.MessageUpdatedPayload) error {
	fileIDs := extractFileIDs(payload.Message.Text)

	insertReq := make([]*model.File, 0, len(fileIDs))
	client, auth := model.NewTraqClient(accessToken)
	for _, v := range fileIDs {
		file, res, err := client.FileApi.GetFileMeta(auth, v).Execute()
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
		}

		if strings.HasPrefix(file.Mime, "audio") {
			insertReq = append(insertReq, &model.File{
				ID:           file.Id,
				Title:        removeExtensions(file.Name),
				ComposerID:   payload.Message.User.ID,
				ComposerName: payload.Message.User.Name,
				MessageID:    payload.Message.ID,
				CreatedAt:    payload.Message.CreatedAt,
			})
		}
	}

	existingFileIDs, err := model.GetFileIDsInMessage(ctx, payload.Message.ID)
	if err != nil {
		return err
	}
	existingFileMap := make(map[string]struct{}, len(existingFileIDs))
	for _, id := range existingFileIDs {
		existingFileMap[id] = struct{}{}
	}

	uniqueInsertReq := make([]*model.File, 0, len(insertReq))
	for _, file := range insertReq {
		if _, exists := existingFileMap[file.ID]; !exists {
			uniqueInsertReq = append(uniqueInsertReq, file)
		}
	}

	if err := model.InsertFiles(ctx, uniqueInsertReq); err != nil {
		return fmt.Errorf("failed to insert file: %w", err)
	}

	newMap := map[string]struct{}{}
	for _, v := range fileIDs {
		newMap[v] = struct{}{}
	}
	oldArr, err := model.GetFileIDsInMessage(ctx, payload.Message.ID)
	if err != nil {
		return err
	}

	deleteFileIds := make([]string, 0, len(oldArr))
	for _, v := range oldArr {
		if _, ok := newMap[v]; !ok {
			deleteFileIds = append(deleteFileIds, v)
		}
	}
	err = model.DeleteFiles(ctx, deleteFileIds)

	return err
}
