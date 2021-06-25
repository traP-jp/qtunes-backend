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
	matches := embURLRegex.FindAllStringSubmatch(payload.Message.Text, -1)

	client, auth := model.NewTraqClient(accessToken)
	for _, v := range matches {
		file, res, err := client.FileApi.GetFileMeta(auth, v[0])
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("failed in HTTP request:(status:%d %s)", res.StatusCode, res.Status)
		}

		if strings.HasPrefix(file.Mime, "audio") {
			req := model.File{
				ID:           file.Id,
				Title:        file.Name,
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

	var newMap map[string]struct{}
	for _, v := range matches {
		newMap[v[1]] = struct{}{}
	}
	oldArr, err := model.GetFileIDsInMessage(ctx, payload.Message.ID)
	if err != nil {
		return err
	}
	for _, v := range oldArr {
		if _, ok := newMap[v]; !ok {
			model.DeleteFile(ctx, v)
		}
	}

	return nil
}
