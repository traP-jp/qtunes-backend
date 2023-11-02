package model

import (
	"context"

	traq "github.com/traPtitech/go-traq"
)

func NewTraqClient(accessToken string) (*traq.APIClient, context.Context) {
	client := traq.NewAPIClient(traq.NewConfiguration())
	auth := context.WithValue(context.Background(), traq.ContextAccessToken, accessToken)

	return client, auth
}
