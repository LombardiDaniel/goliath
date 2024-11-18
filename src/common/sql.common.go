package common

import (
	"strings"
)

const (
	ErrUniqueConstraint string = "duplicate key value violates unique constraint"
)

func FilterSqlPgError(err error) error {
	if err == nil {
		return nil
	}

	if strings.Contains(err.Error(), ErrUniqueConstraint) {
		return ErrDbConflict
	}

	return err
}
