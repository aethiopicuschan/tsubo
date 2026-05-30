package tsubo

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

func ParseBBSMenuHTML(data []byte) (bm *BBSMenu, err error) {
	if len(bytes.TrimSpace(data)) == 0 {
		bm = &BBSMenu{
			categories: make([]BBSMenuCategory, 0),
		}
		return
	}

	var reader io.Reader
	reader, err = charset.NewReader(bytes.NewReader(data), "text/html")
	if err != nil {
		err = errors.Join(ErrDetectBBSMenuHTMLCharset, err)
		return
	}

	var doc *html.Node
	doc, err = html.Parse(reader)
	if err != nil {
		err = errors.Join(ErrParseBBSMenuHTML, err)
		return
	}

	bm = &BBSMenu{
		categories: make([]BBSMenuCategory, 0),
	}

	var current *BBSMenuCategory

	for node := range walkHTMLNodes(doc) {
		if node.Type != html.ElementNode {
			continue
		}

		switch strings.ToLower(node.Data) {
		case "b":
			name := strings.TrimSpace(htmlNodeText(node))
			if name == "" {
				continue
			}

			if current != nil {
				current.total = len(current.boards)
				bm.categories = append(bm.categories, *current)
			}

			current = &BBSMenuCategory{
				number: strconv.Itoa(len(bm.categories) + 1),
				name:   name,
				boards: make([]Board, 0),
			}

		case "a":
			if current == nil {
				continue
			}

			href := htmlAttr(node, "href")
			name := strings.TrimSpace(htmlNodeText(node))
			if href == "" || name == "" {
				continue
			}

			board := newBoard(name, href)
			current.boards = append(current.boards, *board)
		}
	}

	if current != nil {
		current.total = len(current.boards)
		bm.categories = append(bm.categories, *current)
	}

	return
}
