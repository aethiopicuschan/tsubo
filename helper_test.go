package tsubo_test

import "errors"

type errReadCloser struct{}

func (errReadCloser) Read(_ []byte) (int, error) {
	return 0, errors.New("read error")
}

func (errReadCloser) Close() error {
	return nil
}
