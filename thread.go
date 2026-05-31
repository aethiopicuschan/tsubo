package tsubo

import (
	"strconv"
	"time"
)

// Thread represents a thread in a 5ch board, containing its key, title, response count, BE ID, and metadata.
type Thread struct {
	key      string
	title    string
	resCount int
	beID     string
	metadata []ThreadMetadata
}

// ThreadMetadata represents a key-value pair of metadata associated with a thread, along with its raw representation.
type ThreadMetadata struct {
	key   string
	value string
	raw   string
}

// Key returns the key of the thread.
func (t *Thread) Key() string {
	return t.key
}

// Title returns the title of the thread.
func (t *Thread) Title() string {
	return t.title
}

// ResCount returns the response count of the thread.
func (t *Thread) ResCount() int {
	return t.resCount
}

// BeID returns the BE ID of the thread.
func (t *Thread) BeID() string {
	return t.beID
}

// Metadata returns a copy of the list of metadata associated with the thread.
func (t *Thread) Metadata() []ThreadMetadata {
	metadata := make([]ThreadMetadata, len(t.metadata))
	copy(metadata, t.metadata)
	return metadata
}

// Momentum returns the number of responses per day.
//
// The thread key is treated as the Unix timestamp when the thread was created.
// If the key is invalid, resCount is negative, or the thread age is not positive,
// Momentum returns 0.
func (t *Thread) Momentum(now time.Time) float64 {
	createdAt, err := strconv.ParseInt(t.key, 10, 64)
	if err != nil {
		return 0
	}

	if t.resCount < 0 {
		return 0
	}

	age := now.Sub(time.Unix(createdAt, 0))
	if age <= 0 {
		return 0
	}

	return float64(t.resCount) / age.Hours() * 24
}

// Key returns the key of the thread metadata.
func (m ThreadMetadata) Key() string {
	return m.key
}

// Value returns the value of the thread metadata.
func (m ThreadMetadata) Value() string {
	return m.value
}

// Raw returns the raw representation of the thread metadata.
func (m ThreadMetadata) Raw() string {
	return m.raw
}
