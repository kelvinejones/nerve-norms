package mem

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type Reader struct {
	reader          *bufio.Reader
	unreadString    string
	useUnreadString bool
	writeError      bool
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(rd)}
}

func NewStringReader(str string) *Reader {
	return &Reader{reader: bufio.NewReader(strings.NewReader(str))}
}

func (rd *Reader) UnreadString(str string) {
	// If this is called twice without a call to ReadString in between, the next ReadString call will error.
	if rd.useUnreadString {
		rd.writeError = true
	}
	rd.unreadString = str
	rd.useUnreadString = true
}

func (rd *Reader) ReadString(delim byte) (string, error) {
	if rd.writeError {
		return rd.unreadString, errors.New("UnreadString was called twice before ReadString was called")
	}
	if rd.useUnreadString {
		return rd.unreadString, nil
	} else {
		return rd.reader.ReadString(delim)
	}
}
