package tsubo_test

import (
	"testing"

	"github.com/aethiopicuschan/tsubo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBBSMenuHTML(t *testing.T) {
	t.Parallel()

	data := []byte(`
		<html>
			<body>
				<a href="https://portal.example.invalid/">Portal</a>
				<br><br><b>Category A</b><br>
				<a href="https://alpha.example.invalid/board-a/">Board A</a><br>
				<a href="https://beta.example.invalid/board-b/">Board B</a><br>
				<br><br><b>Category B</b><br>
				<a href="https://gamma.example.invalid/board-c/">Board C</a><br>
			</body>
		</html>
	`)

	menu, err := tsubo.ParseBBSMenuHTML(data)
	require.NoError(t, err)
	require.NotNil(t, menu)

	categories := menu.Categories()
	require.Len(t, categories, 2)

	assert.Equal(t, "1", categories[0].Number())
	assert.Equal(t, "Category A", categories[0].Name())
	assert.Equal(t, 2, categories[0].Total())

	boardsA := categories[0].Boards()
	require.Len(t, boardsA, 2)
	assert.Equal(t, "Board A", boardsA[0].Name())
	assert.Equal(t, "https://alpha.example.invalid/board-a/", boardsA[0].URL())
	assert.Equal(t, "Board B", boardsA[1].Name())
	assert.Equal(t, "https://beta.example.invalid/board-b/", boardsA[1].URL())

	assert.Equal(t, "2", categories[1].Number())
	assert.Equal(t, "Category B", categories[1].Name())
	assert.Equal(t, 1, categories[1].Total())

	boardsB := categories[1].Boards()
	require.Len(t, boardsB, 1)
	assert.Equal(t, "Board C", boardsB[0].Name())
	assert.Equal(t, "https://gamma.example.invalid/board-c/", boardsB[0].URL())
}

func TestParseBBSMenuHTMLSkipsAnchorsBeforeFirstCategory(t *testing.T) {
	t.Parallel()

	menu, err := tsubo.ParseBBSMenuHTML([]byte(`
		<html>
			<body>
				<a href="https://portal.example.invalid/">Portal</a>
				<b>Category A</b>
				<a href="https://alpha.example.invalid/board-a/">Board A</a>
			</body>
		</html>
	`))

	require.NoError(t, err)
	require.NotNil(t, menu)

	categories := menu.Categories()
	require.Len(t, categories, 1)

	boards := categories[0].Boards()
	require.Len(t, boards, 1)

	assert.Equal(t, "Board A", boards[0].Name())
}

func TestParseBBSMenuHTMLSkipsInvalidAnchors(t *testing.T) {
	t.Parallel()

	menu, err := tsubo.ParseBBSMenuHTML([]byte(`
		<html>
			<body>
				<b>Category A</b>
				<a>Missing Href</a>
				<a href="https://alpha.example.invalid/board-a/"></a>
				<a href="https://beta.example.invalid/board-b/">Board B</a>
			</body>
		</html>
	`))

	require.NoError(t, err)
	require.NotNil(t, menu)

	categories := menu.Categories()
	require.Len(t, categories, 1)

	boards := categories[0].Boards()
	require.Len(t, boards, 1)

	assert.Equal(t, "Board B", boards[0].Name())
	assert.Equal(t, "https://beta.example.invalid/board-b/", boards[0].URL())
}

func TestParseBBSMenuHTMLEmpty(t *testing.T) {
	t.Parallel()

	menu, err := tsubo.ParseBBSMenuHTML([]byte(``))

	require.NoError(t, err)
	require.NotNil(t, menu)
	assert.Empty(t, menu.Categories())
}
