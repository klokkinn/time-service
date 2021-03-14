package core

func (s StorageError) Unwrap() error {
	return s.err
}

func (s StorageError) StatusCode() int {
	return s.statusCode
}

func (s StorageError) Error() string {
	return s.err.Error()
}
