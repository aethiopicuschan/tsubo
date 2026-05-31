package tsubo

import "errors"

// Common
var (
	ErrReadBody      = errors.New("read body")
	ErrCreateRequest = errors.New("create request")
	ErrDecode        = errors.New("decode charset")
)

// BBS
var (
	ErrFetchBBSMenu             = errors.New("fetch bbsmenu")
	ErrUnexpectedBBSMenuStatus  = errors.New("unexpected bbsmenu status")
	ErrParseBBSMenu             = errors.New("parse bbsmenu")
	ErrParseBBSMenuJSON         = errors.New("parse bbsmenu json")
	ErrDetectBBSMenuHTMLCharset = errors.New("detect bbsmenu html charset")
	ErrParseBBSMenuHTML         = errors.New("parse bbsmenu html")
	ErrUnknownBBSMenuFormat     = errors.New("unknown bbsmenu format")
)

// Board / Subject
var (
	ErrCreateSubjectURL        = errors.New("create subject URL")
	ErrFetchSubject            = errors.New("fetch subject")
	ErrUnexpectedSubjectStatus = errors.New("unexpected subject status")
	ErrParseSubject            = errors.New("parse subject")
)
