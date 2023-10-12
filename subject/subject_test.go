package subject_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/aethiopicuschan/tsubo/subject"
	"github.com/motemen/go-testutil/dataloc"
)

func assertSubject(a subject.Subject, b subject.Subject) error {
	if a.BeID != b.BeID {
		return fmt.Errorf("BeID is wrong, want %s, got %s", a.BeID, b.BeID)
	}
	if a.Title != b.Title {
		return fmt.Errorf("Title is wrong, want %s, got %s", a.Title, b.Title)
	}
	if a.Time != b.Time {
		return fmt.Errorf("Time is wrong, want %d, got %d", a.Time, b.Time)
	}
	if a.ResNum != b.ResNum {
		return fmt.Errorf("ResNum is wrong, want %d, got %d", a.ResNum, b.ResNum)
	}
	return nil
}

func TestNewSubject(t *testing.T) {
	testcases := []struct {
		name      string
		src       string
		expect    subject.Subject
		expectErr bool
	}{
		{
			name: "normal",
			src:  "1660389180.dat<>(ヽ´ん`)「目を閉じて始めよう」  [697453962] (6)",
			expect: subject.Subject{
				BeID:   "697453962",
				Title:  "(ヽ´ん`)「目を閉じて始めよう」",
				Time:   1660389180,
				ResNum: 6,
			},
		},
		{
			name: "confused",
			src:  "9990000000.dat<>123.dat(81)[123456789]  [987654321] (100)",
			expect: subject.Subject{
				BeID:   "987654321",
				Title:  "123.dat(81)[123456789]",
				Time:   9990000000,
				ResNum: 100,
			},
		},
		{
			name: "without_be",
			src:  "1663155465.dat<>OPが良曲すぎた皆忘れてそうなアニメ  (133)",
			expect: subject.Subject{
				BeID:   "",
				Title:  "OPが良曲すぎた皆忘れてそうなアニメ",
				Time:   1663155465,
				ResNum: 133,
			},
		},
		{
			name:      "invalid_time",
			src:       "invalid.dat<>OPが良曲すぎた皆忘れてそうなアニメ  (133)",
			expectErr: true,
		},
		{
			name:      "invalid_resnum",
			src:       "1663155465.dat<>OPが良曲すぎた皆忘れてそうなアニメ  (invalid)",
			expectErr: true,
		},
	}

	for _, testcase := range testcases {
		got, err := subject.NewSubject(testcase.src)
		if testcase.expectErr {
			if err == nil {
				t.Error("want err")
			}
		} else {
			if err != nil {
				t.Errorf("not want err, got \"%s\"", err)
			}
		}

		if err := assertSubject(testcase.expect, got); err != nil {
			t.Errorf("%s, test case at %s", err, dataloc.L(testcase.name))
		}
	}
}

func TestIkioi(t *testing.T) {
	testcases := []struct {
		name   string
		res    int
		time   int64
		now    int64
		expect float64
	}{
		{
			name:   "normal",
			res:    139,
			time:   1660986000,
			now:    1660989360,
			expect: 139.0 / ((1660989360 - 1660986000) / 86400.0),
		},
		{
			name:   "future",
			res:    6,
			time:   9245000000,
			now:    1660989360,
			expect: 0,
		},
	}

	for _, testcase := range testcases {
		got := subject.Ikioi(testcase.res, testcase.time, time.Unix(testcase.now, 0))
		if got != testcase.expect {
			t.Errorf("want %f, got %f, test case at %s", testcase.expect, got, dataloc.L(testcase.name))
		}
	}
}
