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

func TestDetectBBSMenuFormat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		menuURL     string
		contentType string
		body        string
		expected    tsubo.BBSMenuFormat
	}{
		{
			name:        "content type json",
			contentType: "application/json; charset=utf-8",
			expected:    tsubo.BBSMenuFormatJSON,
		},
		{
			name:        "content type html",
			contentType: "text/html; charset=utf-8",
			expected:    tsubo.BBSMenuFormatHTML,
		},
		{
			name:     "url json",
			menuURL:  "https://menu.example.invalid/bbsmenu.json",
			expected: tsubo.BBSMenuFormatJSON,
		},
		{
			name:     "url html",
			menuURL:  "https://menu.example.invalid/bbsmenu.html",
			expected: tsubo.BBSMenuFormatHTML,
		},
		{
			name:     "body object json",
			body:     `{"menu_list":[]}`,
			expected: tsubo.BBSMenuFormatJSON,
		},
		{
			name:     "body array json",
			body:     `[{"name":"test"}]`,
			expected: tsubo.BBSMenuFormatJSON,
		},
		{
			name:     "fallback html",
			body:     `<html><body></body></html>`,
			expected: tsubo.BBSMenuFormatHTML,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := tsubo.DetectBBSMenuFormat(
				tt.menuURL,
				tt.contentType,
				[]byte(tt.body),
			)

			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestParseBBSMenu(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		data   []byte
		format tsubo.BBSMenuFormat
	}{
		{
			name:   "json",
			format: tsubo.BBSMenuFormatJSON,
			data: []byte(`{
				"menu_list": [
					{
						"category_number": "1",
						"category_total": 1,
						"category_name": "Category A",
						"category_content": [
							{
								"board_name": "Board A",
								"url": "https://alpha.example.invalid/board-a/"
							}
						]
					}
				]
			}`),
		},
		{
			name:   "html",
			format: tsubo.BBSMenuFormatHTML,
			data: []byte(`
				<html>
					<body>
						<b>Category A</b>
						<a href="https://alpha.example.invalid/board-a/">Board A</a>
					</body>
				</html>
			`),
		},
		{
			name:   "auto json",
			format: tsubo.BBSMenuFormatAuto,
			data: []byte(`{
				"menu_list": [
					{
						"category_number": "1",
						"category_total": 1,
						"category_name": "Category A",
						"category_content": [
							{
								"board_name": "Board A",
								"url": "https://alpha.example.invalid/board-a/"
							}
						]
					}
				]
			}`),
		},
		{
			name:   "auto html",
			format: tsubo.BBSMenuFormatAuto,
			data: []byte(`
				<html>
					<body>
						<b>Category A</b>
						<a href="https://alpha.example.invalid/board-a/">Board A</a>
					</body>
				</html>
			`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			menu, err := tsubo.ParseBBSMenu(tt.data, tt.format)
			require.NoError(t, err)
			require.NotNil(t, menu)

			categories := menu.Categories()
			require.Len(t, categories, 1)

			assert.Equal(t, "Category A", categories[0].Name())

			boards := categories[0].Boards()
			require.Len(t, boards, 1)

			assert.Equal(t, "Board A", boards[0].Name())
			assert.Equal(t, "https://alpha.example.invalid/board-a/", boards[0].URL())
		})
	}
}

func TestParseBBSMenuUnknownFormat(t *testing.T) {
	t.Parallel()

	menu, err := tsubo.ParseBBSMenu(
		[]byte{},
		tsubo.BBSMenuFormat(999),
	)

	require.Error(t, err)
	assert.Nil(t, menu)
	assert.ErrorIs(t, err, tsubo.ErrUnknownBBSMenuFormat)
}

func TestFetchBBSMenu(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		menuURL     string
		contentType string
		body        string
		wantName    string
	}{
		{
			name:        "fetch json",
			menuURL:     "https://menu.example.invalid/bbsmenu.json",
			contentType: "application/json",
			body: `{
				"menu_list": [
					{
						"category_number": "1",
						"category_total": 1,
						"category_name": "Category A",
						"category_content": [
							{
								"board_name": "Board A",
								"url": "https://alpha.example.invalid/board-a/"
							}
						]
					}
				]
			}`,
			wantName: "Category A",
		},
		{
			name:        "fetch html",
			menuURL:     "https://menu.example.invalid/bbsmenu.html",
			contentType: "text/html; charset=utf-8",
			body: `
				<html>
					<body>
						<b>Category B</b>
						<a href="https://beta.example.invalid/board-b/">Board B</a>
					</body>
				</html>
			`,
			wantName: "Category B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			do := func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, tt.menuURL, req.URL.String())

				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
					Header: http.Header{
						"Content-Type": []string{tt.contentType},
					},
					Body: io.NopCloser(
						strings.NewReader(tt.body),
					),
				}, nil
			}

			menu, err := tsubo.FetchBBSMenu(
				context.Background(),
				do,
				tt.menuURL,
			)

			require.NoError(t, err)
			require.NotNil(t, menu)
			assert.Equal(t, tt.menuURL, menu.BBSMenuURL())

			categories := menu.Categories()
			require.Len(t, categories, 1)
			assert.Equal(t, tt.wantName, categories[0].Name())
		})
	}
}

func TestFetchBBSMenuError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		menuURL   string
		do        func(req *http.Request) (*http.Response, error)
		assertion func(t *testing.T, err error)
	}{
		{
			name:    "invalid url",
			menuURL: "://invalid-url",
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
			name:    "do error",
			menuURL: "https://menu.example.invalid/bbsmenu.json",
			do: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			assertion: func(t *testing.T, err error) {
				t.Helper()
				assert.ErrorIs(t, err, tsubo.ErrFetchBBSMenu)
				assert.ErrorContains(t, err, "network error")
			},
		},
		{
			name:    "unexpected status",
			menuURL: "https://menu.example.invalid/bbsmenu.json",
			do: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     "500 Internal Server Error",
					Body: io.NopCloser(
						strings.NewReader(""),
					),
				}, nil
			},
			assertion: func(t *testing.T, err error) {
				t.Helper()
				assert.ErrorIs(t, err, tsubo.ErrFetchBBSMenu)
				assert.ErrorIs(t, err, tsubo.ErrUnexpectedBBSMenuStatus)
				assert.ErrorContains(t, err, "500 Internal Server Error")
			},
		},
		{
			name:    "read body error",
			menuURL: "https://menu.example.invalid/bbsmenu.json",
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
			name:    "parse error",
			menuURL: "https://menu.example.invalid/bbsmenu.json",
			do: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: io.NopCloser(
						strings.NewReader("{"),
					),
				}, nil
			},
			assertion: func(t *testing.T, err error) {
				t.Helper()
				assert.ErrorIs(t, err, tsubo.ErrParseBBSMenu)
				assert.ErrorIs(t, err, tsubo.ErrParseBBSMenuJSON)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			menu, err := tsubo.FetchBBSMenu(
				context.Background(),
				tt.do,
				tt.menuURL,
			)

			require.Error(t, err)
			assert.Nil(t, menu)

			tt.assertion(t, err)
		})
	}
}

func TestBBSMenuCategoriesReturnsCopy(t *testing.T) {
	t.Parallel()

	menu, err := tsubo.ParseBBSMenuJSON([]byte(`{
		"menu_list": [
			{
				"category_number": "1",
				"category_total": 1,
				"category_name": "Category A",
				"category_content": [
					{
						"board_name": "Board A",
						"url": "https://alpha.example.invalid/board-a/"
					}
				]
			}
		]
	}`))
	require.NoError(t, err)

	categories := menu.Categories()
	require.Len(t, categories, 1)

	categories[0] = tsubo.BBSMenuCategory{}

	assert.Equal(t, "Category A", menu.Categories()[0].Name())
}

func TestBBSMenuCategoryBoardsReturnsCopy(t *testing.T) {
	t.Parallel()

	menu, err := tsubo.ParseBBSMenuJSON([]byte(`{
		"menu_list": [
			{
				"category_number": "1",
				"category_total": 1,
				"category_name": "Category A",
				"category_content": [
					{
						"board_name": "Board A",
						"url": "https://alpha.example.invalid/board-a/"
					}
				]
			}
		]
	}`))
	require.NoError(t, err)

	category := menu.Categories()[0]
	boards := category.Boards()
	require.Len(t, boards, 1)

	boards[0] = tsubo.Board{}

	assert.Equal(t, "Board A", category.Boards()[0].Name())
}
