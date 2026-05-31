package tsubo_test

import (
	"strings"
	"testing"

	"github.com/aethiopicuschan/tsubo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/encoding/japanese"
)

func TestParseSubject(t *testing.T) {
	t.Parallel()

	board := *tsubo.NewBoard(
		"Board A",
		"https://alpha.example.invalid/board-a/",
	)

	data := []byte(strings.Join([]string{
		"1234567890.dat<>Thread A (10)",
		"1234567891.dat<>Thread B &amp; C (20)",
		"1234567892.dat<>Thread C [123456789] (30)",
		"invalid line",
		"",
	}, "\n"))

	subject, err := tsubo.ParseSubject(data, board)

	require.NoError(t, err)
	require.NotNil(t, subject)

	threads := subject.Threads()
	require.Len(t, threads, 3)

	assert.Equal(t, "1234567890", threads[0].Key())
	assert.Equal(t, "Thread A", threads[0].Title())
	assert.Equal(t, 10, threads[0].ResCount())
	assert.Empty(t, threads[0].BeID())
	assert.Empty(t, threads[0].Metadata())

	assert.Equal(t, "1234567891", threads[1].Key())
	assert.Equal(t, "Thread B & C", threads[1].Title())
	assert.Equal(t, 20, threads[1].ResCount())

	assert.Equal(t, "1234567892", threads[2].Key())
	assert.Equal(t, "Thread C", threads[2].Title())
	assert.Equal(t, 30, threads[2].ResCount())
	assert.Equal(t, "123456789", threads[2].BeID())

	metadata := threads[2].Metadata()
	require.Len(t, metadata, 1)
	assert.Equal(t, "be_id", metadata[0].Key())
	assert.Equal(t, "123456789", metadata[0].Value())
	assert.Equal(t, "[123456789]", metadata[0].Raw())
}

func TestParseSubjectShiftJIS(t *testing.T) {
	t.Parallel()

	board := *tsubo.NewBoard(
		"Board A",
		"https://alpha.example.invalid/board-a/",
	)

	encoded, err := japanese.ShiftJIS.NewEncoder().Bytes(
		[]byte("1234567890.dat<>スレッドA (10)\n"),
	)
	require.NoError(t, err)

	subject, err := tsubo.ParseSubject(encoded, board)

	require.NoError(t, err)
	require.NotNil(t, subject)

	threads := subject.Threads()
	require.Len(t, threads, 1)

	assert.Equal(t, "1234567890", threads[0].Key())
	assert.Equal(t, "スレッドA", threads[0].Title())
	assert.Equal(t, 10, threads[0].ResCount())
}

func TestSubjectThreadsReturnsCopy(t *testing.T) {
	t.Parallel()

	board := *tsubo.NewBoard(
		"Board A",
		"https://alpha.example.invalid/board-a/",
	)

	subject, err := tsubo.ParseSubject(
		[]byte("1234567890.dat<>Thread A (10)\n"),
		board,
	)
	require.NoError(t, err)

	threads := subject.Threads()
	require.Len(t, threads, 1)

	threads[0] = tsubo.Thread{}

	assert.Equal(t, "Thread A", subject.Threads()[0].Title())
}

func TestParseSubjectTitleAndMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		expectedTitle string
		expectedCount int
		expectedBeID  string
	}{
		{
			name:          "title and response count",
			input:         "Thread A (10)",
			expectedTitle: "Thread A",
			expectedCount: 10,
		},
		{
			name:          "title response count and be id",
			input:         `Thread A [123456789] (10)`,
			expectedTitle: "Thread A",
			expectedCount: 10,
			expectedBeID:  "123456789",
		},
		{
			name:          "title with html escaped be id",
			input:         `Thread &amp; A [123456789] (10)`,
			expectedTitle: "Thread &amp; A",
			expectedCount: 10,
			expectedBeID:  "123456789",
		},
		{
			name:          "title without response count",
			input:         "Thread A",
			expectedTitle: "Thread A",
			expectedCount: 0,
		},
		{
			name:          "short bracket number stays title",
			input:         "Thread A [123] (10)",
			expectedTitle: "Thread A [123]",
			expectedCount: 10,
		},
		{
			name:          "non numeric bracket stays title",
			input:         "Thread A [abc123] (10)",
			expectedTitle: "Thread A [abc123]",
			expectedCount: 10,
		},
		{
			name:          "invalid response count stays title",
			input:         "Thread A (abc)",
			expectedTitle: "Thread A (abc)",
			expectedCount: 0,
		},
		{
			name:          "parentheses in title",
			input:         "Thread (A) (10)",
			expectedTitle: "Thread (A)",
			expectedCount: 10,
		},
		{
			name:          "be id before malformed response count is not parsed",
			input:         "Thread A [123456789] (abc)",
			expectedTitle: "Thread A [123456789] (abc)",
			expectedCount: 0,
			expectedBeID:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			title, count, beID, metadata := tsubo.ParseSubjectTitleAndMetadata(tt.input)

			assert.Equal(t, tt.expectedTitle, title)
			assert.Equal(t, tt.expectedCount, count)
			assert.Equal(t, tt.expectedBeID, beID)

			if tt.expectedBeID == "" {
				assert.Empty(t, metadata)
				return
			}

			require.Len(t, metadata, 1)
			assert.Equal(t, "be_id", metadata[0].Key())
			assert.Equal(t, tt.expectedBeID, metadata[0].Value())
			assert.Equal(t, "["+tt.expectedBeID+"]", metadata[0].Raw())
		})
	}
}

func TestParseSubjectResCount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		expectedTitle string
		expectedCount int
	}{
		{
			name:          "valid",
			input:         "Thread A (10)",
			expectedTitle: "Thread A",
			expectedCount: 10,
		},
		{
			name:          "no count",
			input:         "Thread A",
			expectedTitle: "Thread A",
		},
		{
			name:          "invalid count",
			input:         "Thread A (abc)",
			expectedTitle: "Thread A (abc)",
		},
		{
			name:          "no space before count",
			input:         "Thread A(10)",
			expectedTitle: "Thread A(10)",
		},
		{
			name:          "parentheses in title",
			input:         "Thread (A) (10)",
			expectedTitle: "Thread (A)",
			expectedCount: 10,
		},
		{
			name:          "trims title",
			input:         "  Thread A   (10)",
			expectedTitle: "Thread A",
			expectedCount: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			title, count := tsubo.ParseSubjectResCount(tt.input)

			assert.Equal(t, tt.expectedTitle, title)
			assert.Equal(t, tt.expectedCount, count)
		})
	}
}

func TestParseSubjectBeID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		expectedTitle string
		expectedBeID  string
		expectedRaw   string
	}{
		{
			name:          "valid",
			input:         "Thread A [123456789]",
			expectedTitle: "Thread A",
			expectedBeID:  "123456789",
			expectedRaw:   "[123456789]",
		},
		{
			name:          "short number is ignored",
			input:         "Thread A [123]",
			expectedTitle: "Thread A [123]",
		},
		{
			name:          "non numeric is ignored",
			input:         "Thread A [abc123]",
			expectedTitle: "Thread A [abc123]",
		},
		{
			name:          "middle be like text is ignored",
			input:         "Thread [123456789] A",
			expectedTitle: "Thread [123456789] A",
		},
		{
			name:          "trims title",
			input:         "  Thread A   [123456789]  ",
			expectedTitle: "Thread A",
			expectedBeID:  "123456789",
			expectedRaw:   "[123456789]",
		},
		{
			name:          "empty",
			input:         "",
			expectedTitle: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			title, beID, raw := tsubo.ParseSubjectBeID(tt.input)

			assert.Equal(t, tt.expectedTitle, title)
			assert.Equal(t, tt.expectedBeID, beID)
			assert.Equal(t, tt.expectedRaw, raw)
		})
	}
}
