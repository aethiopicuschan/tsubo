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

// == Board ==
var BoardNameFromURL = boardNameFromURL
var ParseBoardNameHTML = parseBoardNameHTML
var NormalizeBoardName = normalizeBoardName

// == Subject ==
var ParseSubjectTitleAndMetadata = parseSubjectTitleAndMetadata
var ParseSubjectResCount = parseSubjectResCount
var ParseSubjectBeID = parseSubjectBeID

// == Thread ==
func NewThreadForTest(key string, title string, resCount int, beID string, metadata []ThreadMetadata) Thread {
	return Thread{
		key:      key,
		title:    title,
		resCount: resCount,
		beID:     beID,
		metadata: metadata,
	}
}
