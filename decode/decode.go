package decode

import (
	"html"
	"io"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func Decode(src io.Reader) (result string, err error) {
	raw, err := io.ReadAll(transform.NewReader(src, japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return
	}
	result = html.UnescapeString(string(raw))
	return
}
