package it

// Must ensures the the expression does not error, it panics in case of an error.
// It is tipically used on the main func, where the error would be fatal, such as
// connecting to a db server etc. Returns the first arg
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

// MustNotErr ensures the the expression does not error, it panics in case of an error.
// It is tipically used on the main func, where the error would be fatal, such as
// connecting to a db server etc.
func MustNotErr(err error) {
	if err != nil {
		panic(err)
	}
}
