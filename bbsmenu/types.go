package bbsmenu

type CategoryContent struct {
	BoardName     string `json:"board_name"`
	URL           string `json:"url"`
	CategoryOrder int    `json:"category_order"`
	Category      int    `json:"category"`
	DirectoryName string `json:"directory_name"`
	CategoryName  string `json:"category_name"`
}

type CategoryContents []CategoryContent

type Menu struct {
	CategoryContent CategoryContents `json:"category_content"`
	CategoryTotal   int              `json:"category_total"`
	CategoryNumber  string           `json:"category_number"`
	CategoryName    string           `json:"category_name"`
}

type MenuList []Menu

type BBSMenu struct {
	Description      string   `json:"description"`
	LastModify       int64    `json:"last_modify"`
	MenuList         MenuList `json:"menu_list"`
	LastModifyString string   `json:"last_modify_string"`
}
