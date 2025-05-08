package constants

import "errors"

var (
	ErrNoRows              = errors.New("db no rows")
	ErrAuth                = errors.New("auth error")
	ErrDbConflict          = errors.New("db conflict error")
	ErrDbTransactionCreate = errors.New("could not create DB transaction")
)
