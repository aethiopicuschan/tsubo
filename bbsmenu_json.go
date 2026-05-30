package tsubo

import (
	"encoding/json"
	"errors"
)

type bbsMenuJSON struct {
	LastModifyString string                `json:"last_modify_string"`
	LastModify       int                   `json:"last_modify"`
	MenuList         []bbsMenuJSONCategory `json:"menu_list"`
	Description      string                `json:"description"`
}

type bbsMenuJSONCategory struct {
	CategoryNumber  string             `json:"category_number"`
	CategoryTotal   int                `json:"category_total"`
	CategoryName    string             `json:"category_name"`
	CategoryContent []bbsMenuJSONBoard `json:"category_content"`
}

type bbsMenuJSONBoard struct {
	CategoryName  string `json:"category_name"`
	DirectoryName string `json:"directory_name"`
	BoardName     string `json:"board_name"`
	CategoryOrder int    `json:"category_order"`
	URL           string `json:"url"`
	Category      int    `json:"category"`
}

func ParseBBSMenuJSON(data []byte) (bm *BBSMenu, err error) {
	var raw bbsMenuJSON
	if err = json.Unmarshal(data, &raw); err != nil {
		err = errors.Join(ErrParseBBSMenuJSON, err)
		return
	}

	bm = &BBSMenu{
		lastModifiedString: raw.LastModifyString,
		lastModified:       raw.LastModify,
		description:        raw.Description,
		categories:         make([]BBSMenuCategory, 0, len(raw.MenuList)),
	}

	for _, rawCategory := range raw.MenuList {
		category := BBSMenuCategory{
			number: rawCategory.CategoryNumber,
			name:   rawCategory.CategoryName,
			total:  rawCategory.CategoryTotal,
			boards: make([]Board, 0, len(rawCategory.CategoryContent)),
		}

		for _, rawBoard := range rawCategory.CategoryContent {
			board := newBoard(rawBoard.BoardName, rawBoard.URL)
			category.boards = append(category.boards, *board)
		}

		bm.categories = append(bm.categories, category)
	}

	return
}
