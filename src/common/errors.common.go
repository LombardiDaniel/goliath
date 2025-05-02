package common

import "errors"

var (
	ErrAuth                = errors.New("authError")
	ErrDbConflict          = errors.New("dbConflictError")
	ErrDbTransactionCreate = errors.New("could not create DB transaction")
)
