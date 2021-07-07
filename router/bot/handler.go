package bot

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	traqbot "github.com/traPtitech/traq-bot"
)

const (
	botEventHeader = "X-TRAQ-BOT-EVENT"
	botTokenHeader = "X-TRAQ-BOT-TOKEN"
)

func Handler(c echo.Context) error {
	token := c.Request().Header.Get(botTokenHeader)
	if token != verificationToken {
		return c.NoContent(http.StatusUnauthorized)
	}

	event := c.Request().Header.Get(botEventHeader)
	if len(event) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	ctx := c.Request().Context()

	switch event {
	case "MESSAGE_CREATED":
		payload := &traqbot.MessageCreatedPayload{}
		err := c.Bind(payload)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err))
		}

		err = MessageCreatedHandler(ctx, accessToken, payload)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("failed to handle event: %w", err))
		}
	case "MESSAGE_UPDATED":
		payload := &traqbot.MessageUpdatedPayload{}
		err := c.Bind(payload)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err))
		}

		err = MessageUpdatedHandler(ctx, accessToken, payload)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, fmt.Errorf("failed to handle event: %w", err))
		}
	case "MESSAGE_DELETED":
		payload := &traqbot.MessageDeletedPayload{}
		err := c.Bind(payload)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err))
		}

		err = MessageDeletedHandler(ctx, payload)
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
