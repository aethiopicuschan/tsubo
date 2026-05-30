package tsubo

import (
	"io"
	"iter"
	"strings"

	"golang.org/x/net/html"
)

// walkHTMLNodes returns a sequence of all HTML nodes in the document rooted at root, traversed in preorder.
func walkHTMLNodes(root *html.Node) iter.Seq[*html.Node] {
	return func(yield func(*html.Node) bool) {
		var walk func(*html.Node) bool

		walk = func(n *html.Node) bool {
			if n == nil {
				return true
			}

			if !yield(n) {
				return false
			}

			for child := n.FirstChild; child != nil; child = child.NextSibling {
				if !walk(child) {
					return false
				}
			}

			return true
		}

		walk(root)
	}
}

// htmlNodeText returns the concatenated text content of the HTML node n and all its descendants.
func htmlNodeText(n *html.Node) string {
	var builder strings.Builder

	var walk func(*html.Node)

	walk = func(n *html.Node) {
		if n == nil {
			return
		}

		if n.Type == html.TextNode {
			_, _ = io.WriteString(&builder, n.Data)
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(n)

	return builder.String()
}

// htmlAttr returns the value of the attribute with the specified key in the HTML node n, or an empty string if no such attribute exists. The key comparison is case-insensitive, and the returned value is trimmed of leading and trailing whitespace.
func htmlAttr(n *html.Node, key string) (val string) {
	if n == nil {
		return
	}

	for _, attr := range n.Attr {
		if strings.EqualFold(attr.Key, key) {
			val = strings.TrimSpace(attr.Val)
			return
		}
	}

	return
}
