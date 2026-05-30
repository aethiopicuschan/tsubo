package tsubo_test

import (
	"strings"
	"testing"

	"github.com/aethiopicuschan/tsubo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestWalkHTMLNodes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:  "walks all element nodes in preorder",
			input: `<html><body><div><p>Hello</p><a href="https://alpha.example.invalid/">Link</a></div></body></html>`,
			expected: []string{
				"html",
				"head",
				"body",
				"div",
				"p",
				"a",
			},
		},
		{
			name:  "walks nested element nodes",
			input: `<section><article><h1>Title</h1><p>Body</p></article></section>`,
			expected: []string{
				"html",
				"head",
				"body",
				"section",
				"article",
				"h1",
				"p",
			},
		},
		{
			name:  "walks empty document",
			input: ``,
			expected: []string{
				"html",
				"head",
				"body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			doc, err := html.Parse(strings.NewReader(tt.input))
			require.NoError(t, err)

			var got []string
			for node := range tsubo.WalkHTMLNodes(doc) {
				if node.Type == html.ElementNode {
					got = append(got, node.Data)
				}
			}

			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestWalkHTMLNodesStop(t *testing.T) {
	t.Parallel()

	doc, err := html.Parse(
		strings.NewReader(`<div><p>first</p><p>second</p></div>`),
	)
	require.NoError(t, err)

	var got []string
	tsubo.WalkHTMLNodes(doc)(func(node *html.Node) bool {
		if node.Type != html.ElementNode {
			return true
		}

		got = append(got, node.Data)

		return node.Data != "p"
	})

	assert.Equal(t, []string{"html", "head", "body", "div", "p"}, got)
}

func TestWalkHTMLNodesNilRoot(t *testing.T) {
	t.Parallel()

	var got []*html.Node
	for node := range tsubo.WalkHTMLNodes(nil) {
		got = append(got, node)
	}

	assert.Empty(t, got)
}

func TestHTMLNodeText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		selector string
		expected string
	}{
		{
			name:     "returns nested text",
			input:    `<div>Hello <span>world</span>!</div>`,
			selector: "div",
			expected: "Hello world!",
		},
		{
			name:     "returns text from anchor",
			input:    `<a href="https://alpha.example.invalid/">Board A</a>`,
			selector: "a",
			expected: "Board A",
		},
		{
			name:     "concatenates multiple nested text nodes",
			input:    `<p>foo<strong>bar</strong><em>baz</em></p>`,
			selector: "p",
			expected: "foobarbaz",
		},
		{
			name:     "returns empty string when no text exists",
			input:    `<div><br></div>`,
			selector: "div",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			doc, err := html.Parse(strings.NewReader(tt.input))
			require.NoError(t, err)

			node := findElement(t, doc, tt.selector)

			assert.Equal(t, tt.expected, tsubo.HTMLNodeText(node))
		})
	}
}

func TestHTMLNodeTextNilNode(t *testing.T) {
	t.Parallel()

	assert.Empty(t, tsubo.HTMLNodeText(nil))
}

func TestHTMLAttr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		selector string
		key      string
		expected string
	}{
		{
			name:     "returns attribute value",
			input:    `<a href="https://alpha.example.invalid/">Board A</a>`,
			selector: "a",
			key:      "href",
			expected: "https://alpha.example.invalid/",
		},
		{
			name:     "matches key case-insensitively",
			input:    `<a HREF="https://alpha.example.invalid/">Board A</a>`,
			selector: "a",
			key:      "href",
			expected: "https://alpha.example.invalid/",
		},
		{
			name:     "trims spaces",
			input:    `<a href="  https://alpha.example.invalid/  ">Board A</a>`,
			selector: "a",
			key:      "href",
			expected: "https://alpha.example.invalid/",
		},
		{
			name:     "returns empty string when attribute is missing",
			input:    `<a>Board A</a>`,
			selector: "a",
			key:      "href",
			expected: "",
		},
		{
			name:     "returns empty string when node has no attributes",
			input:    `<div>Board A</div>`,
			selector: "div",
			key:      "href",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			doc, err := html.Parse(strings.NewReader(tt.input))
			require.NoError(t, err)

			node := findElement(t, doc, tt.selector)

			assert.Equal(t, tt.expected, tsubo.HTMLAttr(node, tt.key))
		})
	}
}

func TestHTMLAttrNilNode(t *testing.T) {
	t.Parallel()

	assert.Empty(t, tsubo.HTMLAttr(nil, "href"))
}

func findElement(t *testing.T, root *html.Node, tag string) *html.Node {
	t.Helper()

	for node := range tsubo.WalkHTMLNodes(root) {
		if node.Type == html.ElementNode && node.Data == tag {
			return node
		}
	}

	t.Fatalf("element %q not found", tag)
	return nil
}
