package common

import "errors"

var (
	ErrAuth                = errors.New("auth error")
	ErrDbConflict          = errors.New("db conflict error")
	ErrDbTransactionCreate = errors.New("could not create DB transaction")
)
