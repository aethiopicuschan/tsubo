package tsubo

import (
	"bytes"
	"io"
	"unicode/utf8"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// isUTF8 reports whether data is valid UTF-8.
func isUTF8(data []byte) bool {
	return utf8.Valid(data)
}

// isShiftJIS reports whether data can be interpreted as Shift_JIS.
//
// Note that ASCII-only data is valid as both UTF-8 and Shift_JIS.
func isShiftJIS(data []byte) bool {
	for i := 0; i < len(data); i++ {
		b := data[i]

		switch {
		case b <= 0x7F:
			continue

		case 0xA1 <= b && b <= 0xDF:
			continue

		case (0x81 <= b && b <= 0x9F) ||
			(0xE0 <= b && b <= 0xFC):
			if i+1 >= len(data) {
				return false
			}

			c := data[i+1]
			if !(0x40 <= c && c <= 0xFC && c != 0x7F) {
				return false
			}

			i++

		default:
			return false
		}
	}

	return true
}

// DecodeText decodes data into UTF-8 text.
//
// If data is already valid UTF-8, it is returned as-is.
// Otherwise, data is decoded as Shift_JIS.
func decodeText(data []byte) (decoded []byte, err error) {
	if isUTF8(data) {
		decoded = data
		return
	}

	decoded, err = decodeShiftJIS(data)
	return
}

func decodeShiftJIS(data []byte) (decoded []byte, err error) {
	reader := transform.NewReader(
		bytes.NewReader(data),
		japanese.ShiftJIS.NewDecoder(),
	)

	decoded, err = io.ReadAll(reader)
	return
}
