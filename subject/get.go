package subject

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/aethiopicuschan/tsubo/decode"
)

func Get(board string) (subjects []Subject, err error) {
	url, err := url.Parse(board)
	if err != nil {
		return
	}
	url.Path = path.Join(url.Path, "subject.txt")
	res, err := http.Get(url.String())
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("StatusCode was %s.", res.Status))
	}
	str, err := decode.Decode(res.Body)
	if err != nil {
		return
	}
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		subject, err := newSubject(line)
		if err == nil {
			subjects = append(subjects, subject)
		}
	}
	return
}
