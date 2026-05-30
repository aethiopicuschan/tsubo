package tsubo

import "golang.org/x/net/html"

var WalkHTMLNodes = walkHTMLNodes
var HTMLNodeText = htmlNodeText
var HTMLAttr = htmlAttr

func NewHTMLNodeForTest(nodeType html.NodeType, data string, attrs ...html.Attribute) *html.Node {
	return &html.Node{
		Type: nodeType,
		Data: data,
		Attr: attrs,
	}
}
