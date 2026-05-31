package tsubo_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aethiopicuschan/tsubo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBoard(t *testing.T) {
	t.Parallel()

	board := tsubo.NewBoard(
		"Board A",
		"https://alpha.example.invalid/board-a/",
	)

	require.NotNil(t, board)
	assert.Equal(t, "Board A", board.Name())
	assert.Equal(t, "https://alpha.example.invalid/board-a/", board.URL())
}

func TestFetchSubject(t *testing.T) {
	t.Parallel()

	board := *tsubo.NewBoard(
		"Board A",
		"https://alpha.example.invalid/board-a/",
	)

	body := "1234567890.dat<>Thread A (10)\n"

	do := func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(
			t,
			"https://alpha.example.invalid/board-a/subject.txt",
			req.URL.String(),
		)

		return &http.Response{
			StatusCode: http.StatusOK,
			Status:     "200 OK",
			Body:       io.NopCloser(strings.NewReader(body)),
		}, nil
	}

	subject, err := tsubo.FetchSubject(
		context.Background(),
		do,
		board,
	)

	require.NoError(t, err)
	require.NotNil(t, subject)

	threads := subject.Threads()
	require.Len(t, threads, 1)

	assert.Equal(t, "1234567890", threads[0].Key())
	assert.Equal(t, "Thread A", threads[0].Title())
	assert.Equal(t, 10, threads[0].ResCount())
}

func TestFetchSubjectError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		board     tsubo.Board
		do        func(req *http.Request) (*http.Response, error)
		assertion func(t *testing.T, err error)
	}{
		{
			name: "invalid subject url",
			board: *tsubo.NewBoard(
				"Board A",
				"://invalid-url",
			),
			do: func(req *http.Request) (*http.Response, error) {
				t.Fatal("do should not be called")
				return nil, nil
			},
			assertion: func(t *testing.T, err error) {
				t.Helper()

				assert.ErrorIs(t, err, tsubo.ErrCreateSubjectURL)
				assert.ErrorContains(t, err, "board URL")
			},
		},
		{
			name: "invalid board url",
			board: *tsubo.NewBoard(
				"Board A",
				"http://[::1",
			),
			do: func(req *http.Request) (*http.Response, error) {
				t.Fatal("do should not be called")
				return nil, nil
			},
			assertion: func(t *testing.T, err error) {
				t.Helper()

				assert.ErrorIs(t, err, tsubo.ErrCreateSubjectURL)
				assert.ErrorContains(t, err, "missing ']' in host")
			},
		},
		{
			name: "do error",
			board: *tsubo.NewBoard(
				"Board A",
				"https://alpha.example.invalid/board-a/",
			),
			do: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			assertion: func(t *testing.T, err error) {
				t.Helper()

				assert.ErrorIs(t, err, tsubo.ErrFetchSubject)
				assert.ErrorContains(t, err, "network error")
			},
		},
		{
			name: "unexpected status",
			board: *tsubo.NewBoard(
				"Board A",
				"https://alpha.example.invalid/board-a/",
			),
			do: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     "500 Internal Server Error",
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			},
			assertion: func(t *testing.T, err error) {
				t.Helper()

				assert.ErrorIs(t, err, tsubo.ErrFetchSubject)
				assert.ErrorIs(t, err, tsubo.ErrUnexpectedSubjectStatus)
				assert.ErrorContains(t, err, "500 Internal Server Error")
			},
		},
		{
			name: "read body error",
			board: *tsubo.NewBoard(
				"Board A",
				"https://alpha.example.invalid/board-a/",
			),
			do: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
					Body:       errReadCloser{},
				}, nil
			},
			assertion: func(t *testing.T, err error) {
				t.Helper()

				assert.ErrorIs(t, err, tsubo.ErrReadBody)
				assert.ErrorContains(t, err, "read error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			subject, err := tsubo.FetchSubject(
				context.Background(),
				tt.do,
				tt.board,
			)

			require.Error(t, err)
			assert.Nil(t, subject)

			tt.assertion(t, err)
		})
	}
}
