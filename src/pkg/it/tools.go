package it

// Must ensures the the expression does not error, it panics in case of an error.
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
