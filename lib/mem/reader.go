package mem

import (
	"bufio"
	"io"
	"strings"
)

type Reader struct {
	reader *bufio.Reader
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(rd)}
}

func NewStringReader(str string) *Reader {
	return &Reader{reader: bufio.NewReader(strings.NewReader(str))}
}

func (rd *Reader) ReadString(delim byte) (string, error) {
	return rd.reader.ReadString(delim)
}
