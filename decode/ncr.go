package decode

import (
	"fmt"
	"regexp"
	"strconv"
)

func decodeNCR(source string) (result string) {
	re := regexp.MustCompile("&#[0-9]+;|&#x[0-9a-fA-F]+;")
	return re.ReplaceAllStringFunc(source, func(match string) string {
		var cp int64
		if match[2] == 0x78 {
			cp, _ = strconv.ParseInt(match[3:len(match)-1], 16, 32)
		} else {
			cp, _ = strconv.ParseInt(match[2:len(match)-1], 10, 32)
		}

		return fmt.Sprintf("%c", cp)
	})
}
