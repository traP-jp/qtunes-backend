package bot

import (
	"context"
	"fmt"

	"github.com/hackathon-21-spring-02/back-end/model"
	traqbot "github.com/traPtitech/traq-bot"
)

//MessageDeletedHandler MessageDeletedイベントを処理する
func MessageDeletedHandler(ctx context.Context, payload *traqbot.MessageDeletedPayload) error {
	err := model.DeleteFilesFromMessageId(ctx, payload.Message.ID)
	if err != nil {
		return fmt.Errorf("failed to delete files: %w", err)
	}

	return nil
}
