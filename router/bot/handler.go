package bot

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	traqbot "github.com/traPtitech/traq-bot"
)

const (
	botEventHeader = "X-TRAQ-BOT-EVENT"
	botTokenHeader = "X-TRAQ-BOT-TOKEN"
)

var (
	verificationToken = os.Getenv("BOT_VERIFICATION_TOKEN")
)

func Handler(c echo.Context) error {
	token := c.Request().Header.Get(botTokenHeader)
	if token == verificationToken {
		return c.NoContent(http.StatusForbidden)
	}

	event := c.Request().Header.Get(botEventHeader)
	if len(event) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	switch event {
	case "MESSAGE_CREATED":
		ctx := c.Request().Context()
		sess, err := session.Get("sessions", c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("Failed In Getting Session:%w", err))
		}
		accessToken := sess.Values["accessToken"].(string)
		payload := &traqbot.MessageCreatedPayload{}
		err = c.Bind(payload)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err))
		}

		err = MessageCreatedHandler(ctx, accessToken, payload)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("failed to handle event: %w", err))
		}
	case "PING":
		payload := &traqbot.PingPayload{}
		err := c.Bind(payload)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err))
		}

		err = PingHandler(payload)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("failed to handle event: %w", err))
		}
	}

	return c.NoContent(http.StatusNoContent)
}
