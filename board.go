package tsubo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Board represents a board in the 5ch menu.
type Board struct {
	name string
	url  string
}

// FetchBoard fetches the board information from the given board URL and returns a Board instance.
func (c *Client) FetchBoard(ctx context.Context, boardURL string) (board *Board, err error) {
	board, err = FetchBoard(ctx, c.Do, boardURL)
	return
}

// FetchBoard fetches the board information from the given board URL using the provided HTTP client and returns a Board instance.
func FetchBoard(ctx context.Context, do func(req *http.Request) (*http.Response, error), boardURL string) (board *Board, err error) {
	name, err := FetchBoardName(ctx, do, boardURL)
	if err != nil {
		name, err = boardNameFromURL(boardURL)
		if err != nil {
			err = errors.Join(ErrCreateBoard, err)
			return
		}
	}

	board = NewBoard(name, boardURL)

	return
}

// fallback to extract board name from URL path if fetching board name from HTML fails
func boardNameFromURL(boardURL string) (name string, err error) {
	u, err := url.Parse(boardURL)
	if err != nil {
		err = errors.Join(ErrCreateBoardURL, err)
		return
	}

	name = strings.Trim(u.Path, "/")
	if name == "" {
		err = errors.Join(
			ErrCreateBoardURL,
			fmt.Errorf("board path is empty: %s", boardURL),
		)
		return
	}

	parts := strings.Split(name, "/")
	name = parts[len(parts)-1]

	return
}

// FetchBoardName fetches the board name from the given board URL using the provided HTTP client function.
func FetchBoardName(ctx context.Context, do func(req *http.Request) (*http.Response, error), boardURL string) (name string, err error) {
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, boardURL, nil)
	if err != nil {
		err = errors.Join(ErrCreateRequest, err)
		return
	}

	var res *http.Response
	res, err = do(req)
	if err != nil {
		err = errors.Join(ErrFetchBoardName, err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		err = errors.Join(
			ErrFetchBoardName,
			ErrUnexpectedBoardStatus,
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

	name, err = parseBoardNameHTML(body)
	if err != nil {
		err = errors.Join(ErrParseBoardName, err)
		return
	}

	return
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
