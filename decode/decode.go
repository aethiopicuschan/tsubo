package decode

import (
	"io"
	"io/ioutil"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func Decode(src io.ReadCloser) (result string, err error) {
	raw, err := ioutil.ReadAll(transform.NewReader(src, japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return
	}
	result = unescapeHtml(decodeNCR(string(raw)))
	return
}
