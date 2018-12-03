package mem

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strings"
)

type Reader struct {
	reader          *bufio.Reader
	unreadString    string
	useUnreadString bool
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(rd)}
}

func NewStringReader(str string) *Reader {
	return &Reader{reader: bufio.NewReader(strings.NewReader(str))}
}

func (rd *Reader) UnreadString(str string) error {
	// If this is called twice without a call to ReadString in between, the next ReadString call will error.
	if rd.useUnreadString {
		return errors.New("UnreadString was called twice before ReadString was called")
	}
	rd.unreadString = str
	rd.useUnreadString = true
	return nil
}

func (rd *Reader) ReadString(delim byte) (string, error) {
	if rd.useUnreadString {
		rd.useUnreadString = false
		return rd.unreadString, nil
	} else {
		return rd.reader.ReadString(delim)
	}
}

func (rd *Reader) skipNewlines() error {
	s := "\n"
	var err error
	for s == "\n" && err == nil {
		s, err = rd.ReadString('\n')
	}
	if err != nil {
		return err
	}

	return rd.UnreadString(s)
}

func (rd *Reader) skipPast(search string) error {
	err := rd.skipNewlines()
	if err != nil {
		return err
	}

	s, err := rd.ReadString('\n')
	if err == nil && !strings.Contains(s, search) {
		err = errors.New("Could not find '" + search + "'")
	}
	return err
}

type Parser interface {
	Parse([]string) error
}

// parseLines keeps reading as long as the input is empty lines and Parse doesn't return an error.
// Once Parse can't read a line, we're done (but it's not an error). Back up one line and return.
// Even EOF isn't an error here. We let the caller decide that.
func (rd *Reader) parseLines(regex *regexp.Regexp, parser Parser) error {
	for {
		err := rd.skipNewlines()
		if err != nil {
			// We reached EOF
			return nil
		}

		s, err := rd.ReadString('\n')
		if err != nil {
			// We reached EOF
			return nil
		}

		result := regex.FindStringSubmatch(s)

		err = parser.Parse(result)
		if err != nil {
			// The string couldn't be parsed. That's not an error; it just means we're done parsing this regex.
			return rd.UnreadString(s)
		}
	}

	return nil
}
