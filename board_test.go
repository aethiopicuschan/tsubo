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

func TestFetchBoard(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		boardURL      string
		statusCode    int
		status        string
		body          string
		expectedName  string
		expectedURL   string
		expectedPaths []string
	}{
		{
			name:         "uses board name from html title",
			boardURL:     "https://alpha.example.invalid/board-a/",
			statusCode:   http.StatusOK,
			status:       "200 OK",
			body:         `<html><head><title>Board A - 5ch</title></head><body></body></html>`,
			expectedName: "Board A",
			expectedURL:  "https://alpha.example.invalid/board-a/",
			expectedPaths: []string{
				"/board-a/",
			},
		},
		{
			name:         "falls back to url path when html has no title",
			boardURL:     "https://alpha.example.invalid/board-a/",
			statusCode:   http.StatusOK,
			status:       "200 OK",
			body:         `<html><head></head><body></body></html>`,
			expectedName: "board-a",
			expectedURL:  "https://alpha.example.invalid/board-a/",
			expectedPaths: []string{
				"/board-a/",
			},
		},
		{
			name:         "falls back to last url path element",
			boardURL:     "https://alpha.example.invalid/category/board-a/",
			statusCode:   http.StatusOK,
			status:       "200 OK",
			body:         `<html><head></head><body></body></html>`,
			expectedName: "board-a",
			expectedURL:  "https://alpha.example.invalid/category/board-a/",
			expectedPaths: []string{
				"/category/board-a/",
			},
		},
		{
			name:         "falls back when fetching board name fails",
			boardURL:     "https://alpha.example.invalid/board-a/",
			statusCode:   http.StatusInternalServerError,
			status:       "500 Internal Server Error",
			body:         ``,
			expectedName: "board-a",
			expectedURL:  "https://alpha.example.invalid/board-a/",
			expectedPaths: []string{
				"/board-a/",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			do := func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, tt.boardURL, req.URL.String())
				assert.Contains(t, tt.expectedPaths, req.URL.Path)

				return &http.Response{
					StatusCode: tt.statusCode,
					Status:     tt.status,
					Body:       io.NopCloser(strings.NewReader(tt.body)),
				}, nil
			}

			board, err := tsubo.FetchBoard(
				context.Background(),
				do,
				tt.boardURL,
			)

			require.NoError(t, err)
			require.NotNil(t, board)

			assert.Equal(t, tt.expectedName, board.Name())
			assert.Equal(t, tt.expectedURL, board.URL())
		})
	}
}

func TestFetchBoardError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		boardURL  string
		do        func(req *http.Request) (*http.Response, error)
		assertion func(t *testing.T, err error)
	}{
		{
			name:     "invalid url and fallback fails",
			boardURL: "://invalid-url",
			do: func(req *http.Request) (*http.Response, error) {
				t.Fatal("do should not be called")
				return nil, nil
			},
			assertion: func(t *testing.T, err error) {
				t.Helper()

				assert.ErrorIs(t, err, tsubo.ErrCreateBoard)
				assert.ErrorIs(t, err, tsubo.ErrCreateBoardURL)
			},
		},
		{
			name:     "empty path fallback fails",
			boardURL: "https://alpha.example.invalid/",
			do: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			assertion: func(t *testing.T, err error) {
				t.Helper()

				assert.ErrorIs(t, err, tsubo.ErrCreateBoard)
				assert.ErrorIs(t, err, tsubo.ErrCreateBoardURL)
				assert.ErrorContains(t, err, "board path is empty")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			board, err := tsubo.FetchBoard(
				context.Background(),
				tt.do,
				tt.boardURL,
			)

			require.Error(t, err)
			assert.Nil(t, board)

			tt.assertion(t, err)
		})
	}
}

func TestFetchBoardName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		boardURL     string
		body         string
		expectedName string
	}{
		{
			name:         "plain title",
			boardURL:     "https://alpha.example.invalid/board-a/",
			body:         `<html><head><title>Board A</title></head><body></body></html>`,
			expectedName: "Board A",
		},
		{
			name:         "normalized title",
			boardURL:     "https://alpha.example.invalid/board-a/",
			body:         `<html><head><title>Board A - 5ch</title></head><body></body></html>`,
			expectedName: "Board A",
		},
		{
			name:         "trimmed title",
			boardURL:     "https://alpha.example.invalid/board-a/",
			body:         `<html><head><title>  Board A  </title></head><body></body></html>`,
			expectedName: "Board A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			do := func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, tt.boardURL, req.URL.String())

				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
					Body:       io.NopCloser(strings.NewReader(tt.body)),
				}, nil
			}

			name, err := tsubo.FetchBoardName(
				context.Background(),
				do,
				tt.boardURL,
			)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedName, name)
		})
	}
}

func TestFetchBoardNameError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		boardURL  string
		do        func(req *http.Request) (*http.Response, error)
		assertion func(t *testing.T, err error)
	}{
		{
			name:     "create request error",
			boardURL: "://invalid-url",
			do: func(req *http.Request) (*http.Response, error) {
				t.Fatal("do should not be called")
				return nil, nil
			},
			assertion: func(t *testing.T, err error) {
				t.Helper()

				assert.ErrorIs(t, err, tsubo.ErrCreateRequest)
			},
		},
		{
			name:     "do error",
			boardURL: "https://alpha.example.invalid/board-a/",
			do: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			assertion: func(t *testing.T, err error) {
				t.Helper()

				assert.ErrorIs(t, err, tsubo.ErrFetchBoardName)
				assert.ErrorContains(t, err, "network error")
			},
		},
		{
			name:     "unexpected status",
			boardURL: "https://alpha.example.invalid/board-a/",
			do: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     "500 Internal Server Error",
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			},
			assertion: func(t *testing.T, err error) {
				t.Helper()

				assert.ErrorIs(t, err, tsubo.ErrFetchBoardName)
				assert.ErrorIs(t, err, tsubo.ErrUnexpectedBoardStatus)
				assert.ErrorContains(t, err, "500 Internal Server Error")
			},
		},
		{
			name:     "read body error",
			boardURL: "https://alpha.example.invalid/board-a/",
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
		{
			name:     "parse board name error",
			boardURL: "https://alpha.example.invalid/board-a/",
			do: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
					Body:       io.NopCloser(strings.NewReader(`<html><body></body></html>`)),
				}, nil
			},
			assertion: func(t *testing.T, err error) {
				t.Helper()

				assert.ErrorIs(t, err, tsubo.ErrParseBoardName)
				assert.ErrorIs(t, err, tsubo.ErrBoardNameNotFound)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			name, err := tsubo.FetchBoardName(
				context.Background(),
				tt.do,
				tt.boardURL,
			)

			require.Error(t, err)
			assert.Empty(t, name)

			tt.assertion(t, err)
		})
	}
}

func TestBoardNameFromURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		boardURL     string
		expectedName string
		expectedErr  error
	}{
		{
			name:         "simple board",
			boardURL:     "https://alpha.example.invalid/board-a/",
			expectedName: "board-a",
		},
		{
			name:         "nested board",
			boardURL:     "https://alpha.example.invalid/category/board-a/",
			expectedName: "board-a",
		},
		{
			name:         "without trailing slash",
			boardURL:     "https://alpha.example.invalid/board-a",
			expectedName: "board-a",
		},
		{
			name:        "root path",
			boardURL:    "https://alpha.example.invalid/",
			expectedErr: tsubo.ErrCreateBoardURL,
		},
		{
			name:        "invalid url",
			boardURL:    "://invalid-url",
			expectedErr: tsubo.ErrCreateBoardURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			name, err := tsubo.BoardNameFromURL(tt.boardURL)

			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expectedName, name)
		})
	}
}

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
