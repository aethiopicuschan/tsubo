package tsubo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Board represents a board in the 5ch menu.
type Board struct {
	name string
	url  string
}

// NewBoard creates a new Board instance with the given name and URL.
func NewBoard(name, url string) *Board {
	return &Board{
		name: name,
		url:  url,
	}
}

// Name returns the name of the board.
func (b *Board) Name() string {
	return b.name
}

// URL returns the URL of the board.
func (b *Board) URL() string {
	return b.url
}

func (c *Client) FetchSubject(ctx context.Context, board Board) (subject *Subject, err error) {
	subject, err = FetchSubject(ctx, c.Do, board)
	return
}

func FetchSubject(ctx context.Context, do func(req *http.Request) (*http.Response, error), board Board) (subject *Subject, err error) {
	subjectURL, err := url.JoinPath(board.URL(), "subject.txt")
	if err != nil {
		err = errors.Join(
			ErrCreateSubjectURL,
			fmt.Errorf("board URL: %s", board.URL()),
			err,
		)
		return
	}
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, subjectURL, nil)
	if err != nil {
		err = errors.Join(ErrCreateRequest, err)
		return
	}

	var res *http.Response
	res, err = do(req)
	if err != nil {
		err = errors.Join(ErrFetchSubject, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		err = errors.Join(
			ErrFetchSubject,
			ErrUnexpectedSubjectStatus,
			fmt.Errorf("status: %s", res.Status),
		)
		return
	}

	var body []byte
	body, err = io.ReadAll(res.Body)
	if err != nil {
		err = errors.Join(ErrReadBody, err)
		return
	}

	subject, err = ParseSubject(body, board)
	if err != nil {
		err = errors.Join(ErrParseSubject, err)
		return
	}

	return
}
