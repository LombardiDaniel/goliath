package common

const (
	errUniqueConstraint string = "duplicate key value violates unique constraint"
	errNoRows           string = "no rows in result"
)

func FilterSqlPgError(err error) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	switch errStr {
	case errUniqueConstraint:
		return ErrDbConflict
	case errNoRows:
		return ErrNoRows
	}

	return err
}
