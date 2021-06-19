package model

import (
	"context"

	traq "github.com/sapphi-red/go-traq"
)

func newClient(accessToken string) (*traq.APIClient, context.Context) {
	client := traq.NewAPIClient(traq.NewConfiguration())
	auth := context.WithValue(context.Background(), traq.ContextAccessToken, accessToken)

	return client, auth
}
