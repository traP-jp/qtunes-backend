package bot

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hackathon-21-spring-02/back-end/model"
	traqbot "github.com/traPtitech/traq-bot"
)

//MessageUpdatedHandler MessageUpdatedイベントを処理する
func MessageUpdatedHandler(ctx context.Context, accessToken string, payload *traqbot.MessageUpdatedPayload) error {
	fileIDs := extractFileIDs(payload.Message.Text)

	client, auth := model.NewTraqClient(accessToken)
	for _, v := range fileIDs {
		file, res, err := client.FileApi.GetFileMeta(auth, v)
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
		}

		if strings.HasPrefix(file.Mime, "audio") {
			req := model.File{
				ID:           file.Id,
				Title:        removeExtensions(file.Name),
				ComposerID:   payload.Message.User.ID,
				ComposerName: payload.Message.User.Name,
				MessageID:    payload.Message.ID,
				CreatedAt:    payload.Message.CreatedAt,
			}

			err = model.InsertFile(ctx, &req)
			if err != nil {
				return fmt.Errorf("failed to insert file: %w", err)
			}
		}
	}

	newMap := map[string]struct{}{}
	for _, v := range fileIDs {
		newMap[v] = struct{}{}
	}
	oldArr, err := model.GetFileIDsInMessage(ctx, payload.Message.ID)
	if err != nil {
		return err
	}

	req := make([]string, 0, len(oldArr)) // 削除するメッセージの配列
	for _, v := range oldArr {
		if _, ok := newMap[v]; !ok {
			req = append(req, v)
		}
	}
	if err := model.DeleteFiles(ctx, req); err != nil {
		return err
	}

	return nil
}
