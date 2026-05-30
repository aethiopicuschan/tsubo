package tsubo_test

import (
	"testing"

	"github.com/aethiopicuschan/tsubo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/encoding/japanese"
)

func TestIsUTF8(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{
			name:     "ascii",
			data:     []byte("Thread A"),
			expected: true,
		},
		{
			name:     "utf8 japanese",
			data:     []byte("スレッド"),
			expected: true,
		},
		{
			name: "shift jis japanese",
			data: []byte{
				0x83, 0x58,
				0x83, 0x8C,
				0x83, 0x62,
				0x83, 0x68,
			},
			expected: false,
		},
		{
			name:     "invalid utf8",
			data:     []byte{0x80},
			expected: false,
		},
		{
			name:     "empty",
			data:     []byte{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, tsubo.IsUTF8(tt.data))
		})
	}
}

func TestIsShiftJIS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     []byte
		expected bool
	}{
		{
			name:     "ascii is also shift jis compatible",
			data:     []byte("Thread A"),
			expected: true,
		},
		{
			name: "shift jis japanese",
			data: []byte{
				0x83, 0x58,
				0x83, 0x8C,
				0x83, 0x62,
				0x83, 0x68,
			},
			expected: true,
		},
		{
			name:     "half width katakana",
			data:     []byte{0xB1, 0xB2, 0xB3},
			expected: true,
		},
		{
			name:     "invalid lead byte",
			data:     []byte{0x80},
			expected: false,
		},
		{
			name:     "missing trail byte",
			data:     []byte{0x82},
			expected: false,
		},
		{
			name:     "invalid trail byte",
			data:     []byte{0x82, 0x7F},
			expected: false,
		},
		{
			name:     "empty",
			data:     []byte{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, tsubo.IsShiftJIS(tt.data))
		})
	}
}

func TestDecodeText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		data     []byte
		expected string
	}{
		{
			name:     "ascii",
			data:     []byte("Thread A"),
			expected: "Thread A",
		},
		{
			name:     "utf8 japanese",
			data:     []byte("スレッド"),
			expected: "スレッド",
		},
		{
			name: "shift jis japanese",
			data: []byte{
				0x83, 0x58,
				0x83, 0x8C,
				0x83, 0x62,
				0x83, 0x68,
			},
			expected: "スレッド",
		},
		{
			name:     "empty",
			data:     []byte{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			decoded, err := tsubo.DecodeText(tt.data)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(decoded))
		})
	}
}

func TestDecodeTextWithShiftJISEncoder(t *testing.T) {
	t.Parallel()

	encoded, err := japanese.ShiftJIS.NewEncoder().Bytes(
		[]byte("カテゴリA"),
	)
	require.NoError(t, err)

	decoded, err := tsubo.DecodeText(encoded)

	require.NoError(t, err)
	assert.Equal(t, "カテゴリA", string(decoded))
}
