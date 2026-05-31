package tsubo_test

import (
	"strconv"
	"testing"
	"time"

	"github.com/aethiopicuschan/tsubo"
	"github.com/stretchr/testify/assert"
)

func TestThreadMomentum(t *testing.T) {
	t.Parallel()

	now := time.Date(
		2025,
		time.January,
		2,
		0,
		0,
		0,
		0,
		time.UTC,
	)

	tests := []struct {
		name     string
		thread   tsubo.Thread
		expected float64
	}{
		{
			name: "100 responses in one day",
			thread: tsubo.NewThreadForTest(
				strconv.FormatInt(now.Add(-24*time.Hour).Unix(), 10),
				"Thread A",
				100,
				"",
				nil,
			),
			expected: 100,
		},
		{
			name: "200 responses in two days",
			thread: tsubo.NewThreadForTest(
				strconv.FormatInt(now.Add(-48*time.Hour).Unix(), 10),
				"Thread A",
				200,
				"",
				nil,
			),
			expected: 100,
		},
		{
			name: "50 responses in half day",
			thread: tsubo.NewThreadForTest(
				strconv.FormatInt(now.Add(-12*time.Hour).Unix(), 10),
				"Thread A",
				50,
				"",
				nil,
			),
			expected: 100,
		},
		{
			name: "zero responses",
			thread: tsubo.NewThreadForTest(
				strconv.FormatInt(now.Add(-24*time.Hour).Unix(), 10),
				"Thread A",
				0,
				"",
				nil,
			),
			expected: 0,
		},
		{
			name: "invalid key",
			thread: tsubo.NewThreadForTest(
				"invalid",
				"Thread A",
				100,
				"",
				nil,
			),
			expected: 0,
		},
		{
			name: "future thread",
			thread: tsubo.NewThreadForTest(
				strconv.FormatInt(now.Add(24*time.Hour).Unix(), 10),
				"Thread A",
				100,
				"",
				nil,
			),
			expected: 0,
		},
		{
			name: "negative response count",
			thread: tsubo.NewThreadForTest(
				strconv.FormatInt(now.Add(-24*time.Hour).Unix(), 10),
				"Thread A",
				-1,
				"",
				nil,
			),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.InDelta(
				t,
				tt.expected,
				tt.thread.Momentum(now),
				0.0001,
			)
		})
	}
}
