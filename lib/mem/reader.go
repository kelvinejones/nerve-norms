package mem

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Reader struct {
	reader          *bufio.Reader
	unreadString    string
	useUnreadString bool
	isEof           bool
	lastReadLineNum int
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{reader: bufio.NewReader(rd)}
}

func NewStringReader(str string) *Reader {
	return &Reader{reader: bufio.NewReader(strings.NewReader(str))}
}

func (rd *Reader) UnreadString(str string) error {
	// If this is called twice without a call to ReadLine in between, the next ReadLine call will error.
	if rd.useUnreadString {
		return errors.New("UnreadString was called twice before ReadLine was called")
	}
	rd.unreadString = str
	rd.useUnreadString = true
	return nil
}

func (rd *Reader) ReadLine() (string, error) {
	if rd.useUnreadString {
		if rd.isEof {
			return "", io.EOF
		}
		rd.useUnreadString = false
		return rd.unreadString, nil
	} else {
		str, err := rd.reader.ReadString('\n')
		if err == io.EOF {
			rd.isEof = true
		}
		rd.lastReadLineNum++
		return strings.TrimSuffix(str, "\r\n"), err
	}
}

func (rd *Reader) GetLastLineNumber() int {
	return rd.lastReadLineNum
}

// ReadLineExtractingString expects to receive a regex which finds a single string
func (rd *Reader) ReadLineExtractingString(regstring string) (string, error) {
	s, err := rd.skipNewlines()
	if err != nil {
		return s, err
	}

	result := regexp.MustCompile(regstring).FindStringSubmatch(s)
	if len(result) != 2 {
		rd.UnreadString(s)
		return "", fmt.Errorf("Incorrect ReadLineExtractingString length (%d) for '"+regstring+"'", len(result))
	}

	return result[1], nil
}

func (rd *Reader) skipNewlines() (string, error) {
	s := ""
	var err error
	for s == "" && err == nil {
		s, err = rd.ReadLine()
	}

	return s, err
}

func (rd *Reader) skipPast(search string) error {
	s, err := rd.skipNewlines()
	if err != nil {
		return err
	}

	if !strings.Contains(s, search) {
		err = errors.New("Could not find '" + search + "'" + " in line: " + s)
		rd.UnreadString(s)
	}
	return err
}

// parseLines keeps reading as long as the input is empty lines and Parse doesn't return an error.
// Once Parse can't read a line, we're done (but it's not an error). Back up one line and return.
// Even EOF isn't an error here. We let the caller decide that.
func (rd *Reader) parseLines(parser LineParser) error {
	for {
		s, err := rd.skipNewlines()
		if err != nil {
			// We reached EOF
			return nil
		}

		if parser.ParseLine(parser.ParseRegex().FindStringSubmatch(s)) != nil {
			// The string couldn't be parsed. This isn't an error;
			// it just means we're done parsing this regex.
			return rd.UnreadString(s)
		}
	}
}
