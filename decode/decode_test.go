package decode_test

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"testing/iotest"

	"github.com/aethiopicuschan/tsubo/decode"
	"github.com/motemen/go-testutil/dataloc"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func TestDcode(t *testing.T) {
	testcases := []struct {
		name      string
		src       string
		expect    string
		expectErr bool
	}{
		{
			name:   "emoji",
			src:    "家康の関東移封&#9876;左遷どころか大当たりだった&#127919;",
			expect: "家康の関東移封⚔左遷どころか大当たりだった🎯",
		},
		{
			name:   "html",
			src:    "&quot;アレ&quot;みたくなる",
			expect: `"アレ"みたくなる`,
		},
		{
			name:      "error",
			expectErr: true,
		},
	}

	e := japanese.ShiftJIS.NewEncoder()
	for _, testcase := range testcases {
		var reader io.Reader
		if testcase.expectErr {
			reader = iotest.ErrReader(errors.New("test"))
		} else {
			sjis, _, _ := transform.String(e, string(testcase.src))
			reader = bytes.NewReader([]byte(sjis))
		}
		got, err := decode.Decode(reader)
		if testcase.expectErr {
			if err == nil {
				t.Error("want err")
			}
		} else {
			if err != nil {
				t.Errorf("not want err, got \"%s\"", err)
			}
		}

		if got != testcase.expect {
			t.Errorf("want %s, got %s, test case at %s", testcase.expect, got, dataloc.L(testcase.name))
		}
	}
}
