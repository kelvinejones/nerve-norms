package mem

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

type Reader struct {
	reader             *bufio.Reader
	unwrittenString    string
	useUnwrittenString bool
	writeError         bool
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(rd)}
}

func NewStringReader(str string) *Reader {
	return &Reader{reader: bufio.NewReader(strings.NewReader(str))}
}

func (rd *Reader) UnwriteString(str string) {
	// If this is called twice without a call to ReadString in between, the next ReadString call will error.
	if rd.useUnwrittenString {
		rd.writeError = true
	}
	rd.unwrittenString = str
	rd.useUnwrittenString = true
}

func (rd *Reader) ReadString(delim byte) (string, error) {
	if rd.writeError {
		return rd.unwrittenString, errors.New("UnwriteString was called twice before ReadString was called")
	}
	if rd.useUnwrittenString {
		return rd.unwrittenString, nil
	} else {
		return rd.reader.ReadString(delim)
	}
}
