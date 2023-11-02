package bot

import (
	"log"

	traqbot "github.com/traPtitech/traq-bot"
)

// PingHandler pingを処理する
func PingHandler(_ *traqbot.PingPayload) error {
	log.Println("pong!")

	return nil
}
