package tsubo

// Board represents a board in the 5ch menu.
type Board struct {
	name string
	url  string
}

// NewBoard creates a new Board instance with the given name and URL.
func NewBoard(name, url string) *Board {
	return &Board{
		name: name,
		url:  url,
	}
}

// Name returns the name of the board.
func (b *Board) Name() string {
	return b.name
}

// URL returns the URL of the board.
func (b *Board) URL() string {
	return b.url
}
