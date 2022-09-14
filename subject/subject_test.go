package subject_test

import (
	"testing"
	"time"

	"github.com/aethiopicuschan/tsubo/subject"
)

func TestNewSubject(t *testing.T) {
	// 普通のソース
	sj, err := subject.NewSubject("1660389180.dat<>(ヽ´ん`)「目を閉じて始めよう」  [697453962] (6)")
	if err != nil {
		t.Errorf("%s", err)
	}
	if sj.BeID != "697453962" {
		t.Errorf("BeID is wrong, want \"697453962\", got %s", sj.BeID)
	}
	if sj.Title != "(ヽ´ん`)「目を閉じて始めよう」" {
		t.Errorf("Title is wrong, want \"(ヽ´ん`)「目を閉じて始めよう」\", got \"%s\"", sj.Title)
	}
	if sj.Time != 1660389180 {
		t.Errorf("Time is wrong, want 1660389180, got %d", sj.Time)
	}
	if sj.ResNum != 6 {
		t.Errorf("ResNum is wrong, want 6, got %d", sj.ResNum)
	}
	// 怪しいソース
	sj, err = subject.NewSubject("9990000000.dat<>123.dat(81)[123456789]  [987654321] (100)")
	if err != nil {
		t.Errorf("%s", err)
	}
	if sj.BeID != "987654321" {
		t.Errorf("BeID is wrong, want \"987654321\", got %s", sj.BeID)
	}
	if sj.Title != "123.dat(81)[123456789]" {
		t.Errorf("Title is wrong, want \"123.dat[123456789]\", got \"%s\"", sj.Title)
	}
	if sj.Time != 9990000000 {
		t.Errorf("Time is wrong, want 9990000000, got %d", sj.Time)
	}
	if sj.ResNum != 100 {
		t.Errorf("ResNum is wrong, want 100, got %d", sj.ResNum)
	}
	// Beなし
	sj, err = subject.NewSubject("1663155465.dat<>OPが良曲すぎた皆忘れてそうなアニメ  (133)")
	if err != nil {
		t.Errorf("%s", err)
	}
	if sj.BeID != "" {
		t.Errorf("BeID is wrong, want \"\", got %s", sj.BeID)
	}
	if sj.Title != "OPが良曲すぎた皆忘れてそうなアニメ" {
		t.Errorf("Title is wrong, want \"OPが良曲すぎた皆忘れてそうなアニメ\", got \"%s\"", sj.Title)
	}
	if sj.Time != 1663155465 {
		t.Errorf("Time is wrong, want 1663155465, got %d", sj.Time)
	}
	if sj.ResNum != 133 {
		t.Errorf("ResNum is wrong, want 133, got %d", sj.ResNum)
	}
}

func TestIkioi(t *testing.T) {
	ikioi := subject.Ikioi(139, 1660986000, time.Unix(1660989360, 0))
	if ikioi != 3574.285645 {
		t.Errorf("Ikioi is wrong, want 3574.285645, got %f", ikioi)
	}
	ikioi = subject.Ikioi(6, 9245000000, time.Unix(1660989360, 0))
	if ikioi != 0 {
		t.Errorf("Ikioi is wrong, want 0, got %f", ikioi)
	}
}
