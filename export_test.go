package tsubo

import "golang.org/x/net/html"

// == HTML ==
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

// == Charset ==
var IsUTF8 = isUTF8
var IsShiftJIS = isShiftJIS
var DecodeText = decodeText
var DecodeShiftJIS = decodeShiftJIS

// == Subject ==
var ParseSubjectTitleAndMetadata = parseSubjectTitleAndMetadata
var ParseSubjectResCount = parseSubjectResCount
var ParseSubjectBeID = parseSubjectBeID
