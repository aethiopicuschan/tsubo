package decode_test

import (
	"testing"

	"github.com/aethiopicuschan/tsubo/decode"
)

func TestUnescapeHtml(t *testing.T) {
	str := decode.UnescapeHtml("&quot;アレ&quot;みたくなる")
	if str != "\"アレ\"みたくなる" {
		t.Errorf("want \"\"アレ\"みたくなる\", got \"%s\"", str)
	}
}
