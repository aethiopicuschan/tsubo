package tsubo

import (
	"errors"
	"html"
	"regexp"
	"strconv"
	"strings"
)

var subjectBeIDPattern = regexp.MustCompile(`\s+\[(\d{6,})\]$`)

// Subject represents a subject of a board, which contains multiple threads.
type Subject struct {
	board   Board
	threads []Thread
}

// Threads returns a copy of the list of threads in the subject.
func (s *Subject) Threads() []Thread {
	threads := make([]Thread, len(s.threads))
	copy(threads, s.threads)
	return threads
}

// ParseSubject parses the subject data and returns a Subject instance. It takes the raw data of the subject and the corresponding board as input.
func ParseSubject(data []byte, board Board) (subject *Subject, err error) {
	data, err = decodeText(data)
	if err != nil {
		err = errors.Join(ErrDecode, err)
		return
	}

	subject = &Subject{
		board:   board,
		threads: make([]Thread, 0),
	}

	for line := range strings.Lines(string(data)) {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		key, rest, ok := strings.Cut(line, ".dat<>")
		if !ok {
			continue
		}

		title, resCount, beID, metadata := parseSubjectTitleAndMetadata(rest)

		subject.threads = append(subject.threads, Thread{
			key:      key,
			title:    html.UnescapeString(title),
			resCount: resCount,
			beID:     beID,
			metadata: metadata,
		})
	}

	return
}

// parseSubjectTitleAndMetadata parses the title, response count, beID, and metadata from the given string. It returns the title, response count, beID, and a list of metadata.
func parseSubjectTitleAndMetadata(s string) (title string, resCount int, beID string, metadata []ThreadMetadata) {
	s = strings.TrimSpace(s)

	s, resCount = parseSubjectResCount(s)

	var rawBeID string
	s, beID, rawBeID = parseSubjectBeID(s)
	if beID != "" {
		metadata = append(metadata, ThreadMetadata{
			key:   "be_id",
			value: beID,
			raw:   rawBeID,
		})
	}

	title = strings.TrimSpace(s)

	return
}

// parseSubjectResCount parses the response count from the given string. It returns the title without the response count and the response count itself. If the response count cannot be parsed, it returns the original string as the title and 0 as the response count.
func parseSubjectResCount(s string) (title string, resCount int) {
	s = strings.TrimSpace(s)

	open := strings.LastIndex(s, " (")
	close := strings.LastIndex(s, ")")

	if open < 0 || close != len(s)-1 || open >= close {
		title = s
		return
	}

	count, err := strconv.Atoi(s[open+2 : close])
	if err != nil {
		title = s
		return
	}

	title = strings.TrimSpace(s[:open])
	resCount = count

	return
}

// parseSubjectBeID parses the beID from the given string. It returns the title without the beID, the beID itself, and the raw string of the beID including brackets. If the beID cannot be parsed, it returns the original string as the title and empty strings for the beID and raw string.
func parseSubjectBeID(s string) (title string, beID string, raw string) {
	s = strings.TrimSpace(s)

	matches := subjectBeIDPattern.FindStringSubmatchIndex(s)
	if matches == nil {
		title = s
		return
	}

	raw = strings.TrimSpace(s[matches[0]:matches[1]])
	beID = s[matches[2]:matches[3]]
	title = strings.TrimSpace(s[:matches[0]])

	return
}
