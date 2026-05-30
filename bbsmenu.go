package tsubo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// BBSMenu represents the structure of a 5ch BBS menu, which contains metadata and a list of categories.
type BBSMenu struct {
	lastModifiedString string
	lastModified       int
	description        string
	bbsmenuURL         string
	categories         []BBSMenuCategory
}

// LastModifiedString returns the last modified string of the BBS menu.
func (bm *BBSMenu) LastModifiedString() string {
	return bm.lastModifiedString
}

// LastModified returns the last modified timestamp of the BBS menu.
func (bm *BBSMenu) LastModified() int {
	return bm.lastModified
}

// Description returns the description of the BBS menu.
func (bm *BBSMenu) Description() string {
	return bm.description
}

// BBSMenuURL returns the URL of the BBS menu.
func (bm *BBSMenu) BBSMenuURL() string {
	return bm.bbsmenuURL
}

// Categories returns a copy of the list of categories in the BBS menu.
func (bm *BBSMenu) Categories() []BBSMenuCategory {
	categories := make([]BBSMenuCategory, len(bm.categories))
	copy(categories, bm.categories)
	return categories
}

// BBSMenuCategory represents a category in the BBS menu, which contains metadata and a list of boards.
type BBSMenuCategory struct {
	number string
	name   string
	total  int
	boards []Board
}

// Number returns the category number of the BBS menu category.
func (c *BBSMenuCategory) Number() string {
	return c.number
}

// Name returns the name of the BBS menu category.
func (c *BBSMenuCategory) Name() string {
	return c.name
}

// Total returns the total number of boards in the BBS menu category.
func (c *BBSMenuCategory) Total() int {
	return c.total
}

// Boards returns a copy of the list of boards in the BBS menu category.
func (c *BBSMenuCategory) Boards() []Board {
	boards := make([]Board, len(c.boards))
	copy(boards, c.boards)
	return boards
}

// BBSMenuFormat represents the format of the BBS menu, which can be auto-detected, JSON, or HTML.
type BBSMenuFormat int

const (
	BBSMenuFormatAuto BBSMenuFormat = iota
	BBSMenuFormatJSON
	BBSMenuFormatHTML
)

// FetchBBSMenu fetches the BBS menu from the specified URL using the client's HTTP method, and returns a BBSMenu instance.
func (c *Client) FetchBBSMenu(ctx context.Context, menuURL string) (bm *BBSMenu, err error) {
	bm, err = FetchBBSMenu(ctx, c.Do, menuURL)
	return
}

// FetchBBSMenu fetches the BBS menu from the specified URL using the provided HTTP client function, and returns a BBSMenu instance.
func FetchBBSMenu(ctx context.Context, do func(req *http.Request) (*http.Response, error), menuURL string) (bm *BBSMenu, err error) {
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, menuURL, nil)
	if err != nil {
		err = errors.Join(ErrCreateBBSMenuRequest, err)
		return
	}

	var res *http.Response
	res, err = do(req)
	if err != nil {
		err = errors.Join(ErrFetchBBSMenu, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		err = errors.Join(
			ErrFetchBBSMenu,
			ErrUnexpectedBBSMenuStatus,
			fmt.Errorf("status: %s", res.Status),
		)
		return
	}

	var body []byte
	body, err = io.ReadAll(res.Body)
	if err != nil {
		err = errors.Join(ErrReadBBSMenuBody, err)
		return
	}

	format := DetectBBSMenuFormat(menuURL, res.Header.Get("Content-Type"), body)

	bm, err = ParseBBSMenu(body, format)
	if err != nil {
		err = errors.Join(ErrParseBBSMenu, err)
		return
	}
	bm.bbsmenuURL = menuURL

	return
}

// DetectBBSMenuFormat detects the format of the BBS menu based on the URL, content type, and body content.
func DetectBBSMenuFormat(menuURL, contentType string, body []byte) BBSMenuFormat {
	contentType = strings.ToLower(contentType)
	menuURL = strings.ToLower(menuURL)

	switch {
	case strings.Contains(contentType, "application/json"):
		return BBSMenuFormatJSON
	case strings.Contains(contentType, "text/html"):
		return BBSMenuFormatHTML
	case strings.HasSuffix(menuURL, ".json"):
		return BBSMenuFormatJSON
	case strings.HasSuffix(menuURL, ".html"), strings.HasSuffix(menuURL, ".htm"):
		return BBSMenuFormatHTML
	}

	trimmed := strings.TrimSpace(string(body))
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		return BBSMenuFormatJSON
	}

	return BBSMenuFormatHTML
}

// ParseBBSMenu parses the BBS menu from the given data using the specified format, and returns a BBSMenu instance.
func ParseBBSMenu(data []byte, format BBSMenuFormat) (bm *BBSMenu, err error) {
	switch format {
	case BBSMenuFormatJSON:
		bm, err = ParseBBSMenuJSON(data)
	case BBSMenuFormatHTML:
		bm, err = ParseBBSMenuHTML(data)
	case BBSMenuFormatAuto:
		bm, err = ParseBBSMenu(
			data,
			DetectBBSMenuFormat("", "", data),
		)
	default:
		err = errors.Join(
			ErrUnknownBBSMenuFormat,
			fmt.Errorf("format: %d", format),
		)
	}

	return
}
