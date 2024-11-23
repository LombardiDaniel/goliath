package common

import (
	"strings"
)

const (
	errUniqueConstraint string = "duplicate key value violates unique constraint"
	errNoRows           string = "no rows in result"
)

var (
	conflictErrStrings []string = []string{errUniqueConstraint, errNoRows}
)

func FilterSqlPgError(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	for _, v := range conflictErrStrings {
		if strings.Contains(errStr, v) {
			return ErrDbConflict
		}
	}

	return err
}
