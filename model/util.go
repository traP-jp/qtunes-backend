package model

import (
	"database/sql"
	"errors"
	"net/url"
)

const SoundChannelId = "8bd9e07a-2c6a-49e6-9961-4f88e83b4918"

var (
	BaseUrl, _ = url.Parse("https://q.trap.jp/api/v3")
	DBErrs     = map[string]error{
		"ErrNoRows": sql.ErrNoRows,
		"NoChange":  errors.New("No Change"),
	}
)
