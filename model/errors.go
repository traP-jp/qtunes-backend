package model

import (
	"database/sql"
	"errors"
)

var (
	ErrNoRows   = sql.ErrNoRows
	ErrNoChange = errors.New("No Change")
)
