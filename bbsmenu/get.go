package bbsmenu

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

// 引数の例: https://menu.5ch.net/
func Get(menu string) (bm BBSMenu, err error) {
	url, err := url.Parse(menu)
	if err != nil {
		return
	}
	url.Path = path.Join(url.Path, "bbsmenu.json")
	res, err := http.Get(url.String())
	if err != nil {
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("StatusCode was %s.", res.Status))
		return
	}
	body, err := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &bm)
	return
}
