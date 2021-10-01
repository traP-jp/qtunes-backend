package bot

import (
	"context"
	"fmt"
	"strings"

	"github.com/hackathon-21-spring-02/back-end/model"
	traqbot "github.com/traPtitech/traq-bot"
)

//MessageCreatedHandler MessageCreatedイベントを処理する
func MessageCreatedHandler(ctx context.Context, accessToken string, payload *traqbot.MessageCreatedPayload) error {
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

	return nil
}
