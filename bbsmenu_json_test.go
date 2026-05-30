package tsubo_test

import (
	"encoding/json"
	"testing"

	"github.com/aethiopicuschan/tsubo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBBSMenuJSON(t *testing.T) {
	t.Parallel()

	data := []byte(`{
		"last_modify_string": "2026/05/31 12:34:56",
		"last_modify": 1780000000,
		"description": "test menu",
		"menu_list": [
			{
				"category_number": "1",
				"category_total": 2,
				"category_name": "Category A",
				"category_content": [
					{
						"board_name": "Board A",
						"url": "https://alpha.example.invalid/board-a/"
					},
					{
						"board_name": "Board B",
						"url": "https://beta.example.invalid/board-b/"
					}
				]
			},
			{
				"category_number": "2",
				"category_total": 1,
				"category_name": "Category B",
				"category_content": [
					{
						"board_name": "Board C",
						"url": "https://gamma.example.invalid/board-c/"
					}
				]
			}
		]
	}`)

	menu, err := tsubo.ParseBBSMenuJSON(data)
	require.NoError(t, err)
	require.NotNil(t, menu)

	assert.Equal(t, "2026/05/31 12:34:56", menu.LastModifiedString())
	assert.Equal(t, 1780000000, menu.LastModified())
	assert.Equal(t, "test menu", menu.Description())

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

func TestParseBBSMenuJSONInvalid(t *testing.T) {
	t.Parallel()

	menu, err := tsubo.ParseBBSMenuJSON([]byte(`{`))

	require.Error(t, err)
	assert.Nil(t, menu)
	assert.ErrorIs(t, err, tsubo.ErrParseBBSMenuJSON)

	var syntaxError *json.SyntaxError
	assert.ErrorAs(t, err, &syntaxError)
}

func TestParseBBSMenuJSONEmptyMenuList(t *testing.T) {
	t.Parallel()

	menu, err := tsubo.ParseBBSMenuJSON([]byte(`{
		"last_modify_string": "2026/05/31 12:34:56",
		"last_modify": 1780000000,
		"description": "empty menu",
		"menu_list": []
	}`))

	require.NoError(t, err)
	require.NotNil(t, menu)

	assert.Equal(t, "2026/05/31 12:34:56", menu.LastModifiedString())
	assert.Equal(t, 1780000000, menu.LastModified())
	assert.Equal(t, "empty menu", menu.Description())
	assert.Empty(t, menu.Categories())
}
