package subject

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Subject struct {
	Time   int64
	Title  string
	ResNum int
	Src    string
	BeID   string
	Ikioi  float32
}

// 勢い計算 テスト可能にするため切り出し
func ikioi(res int, time int64, now time.Time) float32 {
	// 「5ちゃんねるからのお知らせ」 などの特殊なスレッドで現在時刻を上回った値が設定されていることがある
	if time > now.Unix() {
		return 0
	}
	return float32(res) / (float32(now.Unix()-time) / 86400)
}

func newSubject(src string) (subject Subject, err error) {
	r := regexp.MustCompile(`([0-9]+\.dat)<>(.+)\s+\(([0-9]+)\)$`)
	if !r.MatchString(src) {
		err = errors.New("illegal source")
		return
	}
	a := r.FindStringSubmatch(src)
	// 元の文字列
	subject.Src = src
	// Unixtime兼DAT
	subject.Time, err = strconv.ParseInt(strings.Split(a[1], ".dat")[0], 10, 64)
	// レス数
	subject.ResNum, err = strconv.Atoi(a[3])
	if err != nil {
		return
	}
	// BeIDとタイトル
	r2 := regexp.MustCompile(`\s+\[([0-9]+)\]$`)
	a2 := r2.FindStringSubmatch(a[2])
	if len(a2) > 0 {
		subject.BeID = a2[1]
		subject.Title = strings.Replace(a[2], fmt.Sprintf("  [%s]", a2[1]), "", -1)
	} else {
		r3 := regexp.MustCompile(`\s$`)
		subject.Title = r3.ReplaceAllString(a[2], "")
	}
	subject.Ikioi = ikioi(subject.ResNum, subject.Time, time.Now())
	return
}
