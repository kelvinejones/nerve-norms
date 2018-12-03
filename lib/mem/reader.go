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
	rd.UnreadString(s)

	return err
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

func (rd *Reader) parseLines(regex *regexp.Regexp, parser Parser) error {
	err := rd.skipNewlines()
	if err != nil {
		return err
	}

	for {
		s, err := rd.ReadString('\n')
		if err != nil {
			return err
		}

		if len(s) == 1 {
			// Done with section; break!
			break
		}
		result := regex.FindStringSubmatch(s)

		err = parser.Parse(result)
		if err != nil {
			return errors.New(err.Error() + ": '" + s + "'")
		}
	}

	return nil
}
