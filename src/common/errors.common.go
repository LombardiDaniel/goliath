package common

import "errors"

var (
	ErrAuth       = errors.New("authError")
	ErrDbConflict = errors.New("dbConflictError")
)
