package tsubo

import "errors"

var (
	ErrCreateBBSMenuRequest     = errors.New("create bbsmenu request")
	ErrFetchBBSMenu             = errors.New("fetch bbsmenu")
	ErrUnexpectedBBSMenuStatus  = errors.New("unexpected bbsmenu status")
	ErrReadBBSMenuBody          = errors.New("read bbsmenu body")
	ErrParseBBSMenu             = errors.New("parse bbsmenu")
	ErrParseBBSMenuJSON         = errors.New("parse bbsmenu json")
	ErrDetectBBSMenuHTMLCharset = errors.New("detect bbsmenu html charset")
	ErrParseBBSMenuHTML         = errors.New("parse bbsmenu html")
	ErrUnknownBBSMenuFormat     = errors.New("unknown bbsmenu format")
)
