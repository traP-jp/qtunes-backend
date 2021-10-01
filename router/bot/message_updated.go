package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/hackathon-21-spring-02/back-end/model"
	traqbot "github.com/traPtitech/traq-bot"
)

//MessageUpdatedHandler MessageUpdatedイベントを処理する
func MessageUpdatedHandler(ctx context.Context, accessToken string, payload *traqbot.MessageUpdatedPayload) error {
	fileIDs := extractFileIDs(payload.Message.Text)

	insertReq := make([]*model.File, 0, len(fileIDs))
	traqapi := model.NewTraqAPI(accessToken)
	for _, v := range fileIDs {
		file, err := traqapi.GetFileMeta(v)
		if err != nil {
			return err
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

	if err := model.InsertFiles(ctx, insertReq); err != nil {
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
