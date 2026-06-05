package tsubo_test

import (
	"testing"

	"github.com/aethiopicuschan/tsubo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/encoding/japanese"
)

func TestParseBoardNameHTML(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		data         []byte
		expectedName string
	}{
		{
			name: "plain title",
			data: []byte(`
				<html>
					<head>
						<title>Board A</title>
					</head>
					<body></body>
				</html>
			`),
			expectedName: "Board A",
		},
		{
			name: "normalizes 5ch suffix",
			data: []byte(`
				<html>
					<head>
						<title>Board A - 5ch</title>
					</head>
					<body></body>
				</html>
			`),
			expectedName: "Board A",
		},
		{
			name: "trims spaces",
			data: []byte(`
				<html>
					<head>
						<title>  Board A  </title>
					</head>
					<body></body>
				</html>
			`),
			expectedName: "Board A",
		},
		{
			name: "uses first non empty title",
			data: []byte(`
				<html>
					<head>
						<title>Board A</title>
						<title>Board B</title>
					</head>
					<body></body>
				</html>
			`),
			expectedName: "Board A",
		},
		{
			name: "nested title text",
			data: []byte(`
				<html>
					<head>
						<title>Board <span>A</span></title>
					</head>
					<body></body>
				</html>
			`),
			expectedName: "Board <span>A</span>",
		},
		{
			name: "html entity",
			data: []byte(`
				<html>
					<head>
						<title>Board &amp; Test</title>
					</head>
					<body></body>
				</html>
			`),
			expectedName: "Board & Test",
		},
		{
			name: "html entity with suffix",
			data: []byte(`
				<html>
					<head>
						<title>Board &amp; Test - 5ch</title>
					</head>
					<body></body>
				</html>
			`),
			expectedName: "Board & Test",
		},
		{
			name: "emoji",
			data: []byte(`
				<html>
					<head>
						<title>🍣 Board A 🚀</title>
					</head>
					<body></body>
				</html>
			`),
			expectedName: "🍣 Board A 🚀",
		},
		{
			name: "emoji and html entity",
			data: []byte(`
				<html>
					<head>
						<title>🍣 Board &amp; Test 🚀</title>
					</head>
					<body></body>
				</html>
			`),
			expectedName: "🍣 Board & Test 🚀",
		},
		{
			name: "emoji html entity and suffix",
			data: []byte(`
				<html>
					<head>
						<title>🍣 Board &amp; Test 🚀 - 5ch</title>
					</head>
					<body></body>
				</html>
			`),
			expectedName: "🍣 Board & Test 🚀",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			name, err := tsubo.ParseBoardNameHTML(tt.data)

			require.NoError(t, err)
			assert.Equal(t, tt.expectedName, name)
		})
	}
}

func TestParseBoardNameHTMLShiftJIS(t *testing.T) {
	t.Parallel()

	encoded, err := japanese.ShiftJIS.NewEncoder().Bytes(
		[]byte(`
			<html>
				<head>
					<title>カテゴリA - 5ch</title>
				</head>
				<body></body>
			</html>
		`),
	)
	require.NoError(t, err)

	name, err := tsubo.ParseBoardNameHTML(encoded)

	require.NoError(t, err)
	assert.Equal(t, "カテゴリA", name)
}

func TestParseBoardNameHTMLError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		data        []byte
		expectedErr error
	}{
		{
			name:        "title not found",
			data:        []byte(`<html><head></head><body></body></html>`),
			expectedErr: tsubo.ErrBoardNameNotFound,
		},
		{
			name:        "empty title",
			data:        []byte(`<html><head><title>   </title></head><body></body></html>`),
			expectedErr: tsubo.ErrBoardNameNotFound,
		},
		{
			name:        "empty document",
			data:        []byte(``),
			expectedErr: tsubo.ErrBoardNameNotFound,
		},
		{
			name:        "only whitespace",
			data:        []byte(`   `),
			expectedErr: tsubo.ErrBoardNameNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			name, err := tsubo.ParseBoardNameHTML(tt.data)

			require.Error(t, err)
			assert.Empty(t, name)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestParseBoardNameHTMLShiftJISAndEntity(t *testing.T) {
	t.Parallel()

	encoded, err := japanese.ShiftJIS.NewEncoder().Bytes(
		[]byte(`
			<html>
				<head>
					<title>カテゴリ&amp;A - 5ch</title>
				</head>
				<body></body>
			</html>
		`),
	)
	require.NoError(t, err)

	name, err := tsubo.ParseBoardNameHTML(encoded)

	require.NoError(t, err)
	assert.Equal(t, "カテゴリ&A", name)
}

func TestNormalizeBoardName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain",
			input:    "Board A",
			expected: "Board A",
		},
		{
			name:     "trims spaces",
			input:    "  Board A  ",
			expected: "Board A",
		},
		{
			name:     "removes 5 channel japanese suffix",
			input:    "Board A＠5ちゃんねる",
			expected: "Board A",
		},
		{
			name:     "removes 5ch board suffix",
			input:    "Board A＠5ch掲示板",
			expected: "Board A",
		},
		{
			name:     "removes long japanese suffix",
			input:    "Board A - 5ちゃんねる掲示板",
			expected: "Board A",
		},
		{
			name:     "removes short 5ch suffix",
			input:    "Board A - 5ch",
			expected: "Board A",
		},
		{
			name:     "removes generic board suffix",
			input:    "Board A＠掲示板",
			expected: "Board A",
		},
		{
			name:     "removes suffix and trims again",
			input:    "  Board A   - 5ch  ",
			expected: "Board A",
		},
		{
			name:     "does not remove suffix-like middle text",
			input:    "Board A - 5ch Mirror",
			expected: "Board A - 5ch Mirror",
		},
		{
			name:     "empty",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, tsubo.NormalizeBoardName(tt.input))
		})
	}
}
