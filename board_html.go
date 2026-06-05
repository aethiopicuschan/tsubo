package tsubo

import (
	"errors"
	"strings"

	"golang.org/x/net/html"
)

func parseBoardNameHTML(data []byte) (name string, err error) {
	data, err = decodeText(data)
	if err != nil {
		err = errors.Join(ErrDecode, err)
		return
	}

	var doc *html.Node
	doc, err = html.Parse(strings.NewReader(string(data)))
	if err != nil {
		err = errors.Join(ErrParseBoardHTML, err)
		return
	}

	for node := range walkHTMLNodes(doc) {
		if node.Type != html.ElementNode {
			continue
		}

		if !strings.EqualFold(node.Data, "title") {
			continue
		}

		name = normalizeBoardName(htmlNodeText(node))
		if name != "" {
			return
		}
	}

	err = ErrBoardNameNotFound
	return
}

func normalizeBoardName(name string) string {
	name = strings.TrimSpace(name)

	suffixes := []string{
		"＠5ちゃんねる",
		"＠5ch掲示板",
		" - 5ちゃんねる掲示板",
		" - 5ch",
		"＠掲示板",
	}

	for _, suffix := range suffixes {
		name = strings.TrimSpace(strings.TrimSuffix(name, suffix))
	}

	return name
}
