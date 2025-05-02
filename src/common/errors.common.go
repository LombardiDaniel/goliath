package common

import "errors"

var (
	ErrAuth       = errors.New("auth error")
	ErrDbConflict = errors.New("db conflict error")
	ErrNoRows     = errors.New("db no rows")
)
