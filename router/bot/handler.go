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
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	event := c.Request().Header.Get(botEventHeader)
	if len(event) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	ctx := c.Request().Context()

	switch event {
	case "MESSAGE_CREATED":
		payload := &traqbot.MessageCreatedPayload{}
		err := c.Bind(payload)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err).Error())
		}

		err = MessageCreatedHandler(ctx, accessToken, payload)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to handle event: %w", err).Error())
		}
	case "MESSAGE_UPDATED":
		payload := &traqbot.MessageUpdatedPayload{}
		err := c.Bind(payload)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err).Error())
		}

		err = MessageUpdatedHandler(ctx, accessToken, payload)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to handle event: %w", err).Error())
		}
	case "MESSAGE_DELETED":
		payload := &traqbot.MessageDeletedPayload{}
		err := c.Bind(payload)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err).Error())
		}

		err = MessageDeletedHandler(ctx, payload)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to handle event: %w", err).Error())
		}
	case "PING":
		payload := &traqbot.PingPayload{}
		err := c.Bind(payload)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("failed to bind request: %w", err).Error())
		}

		err = PingHandler(payload)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("failed to handle event: %w", err).Error())
		}
	}

	return echo.NewHTTPError(http.StatusNoContent)
}
