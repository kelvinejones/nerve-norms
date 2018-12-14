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
		str := ""
		var err error

		// Skip past blank lines
		for str == "" {
			if str, err = rd.reader.ReadString('\n'); err != nil {
				if err == io.EOF {
					rd.isEof = true
				}
				return str, err
			}
			str = strings.TrimSuffix(str, "\r\n")
		}
		return str, err
	}
}

// ReadLineExtractingString expects to receive a regex which finds a single string
func (rd *Reader) ReadLineExtractingString(regstring string) (string, error) {
	err := rd.skipNewlines()
	if err != nil {
		return "", err
	}

	s, err := rd.ReadLine()
	if err != nil {
		return "", err
	}

	result := regexp.MustCompile(regstring).FindStringSubmatch(s)
	if len(result) != 2 {
		return "", fmt.Errorf("Incorrect ReadLineExtractingString length (%d) for '"+regstring+"'", len(result))
	}

	return result[1], nil
}

func (rd *Reader) skipNewlines() error {
	s := ""
	var err error
	for s == "" && err == nil {
		s, err = rd.ReadLine()
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

	s, err := rd.ReadLine()
	if err == nil && !strings.Contains(s, search) {
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
		err := rd.skipNewlines()
		if err != nil {
			// We reached EOF
			return nil
		}

		s, err := rd.ReadLine()
		if err != nil {
			// We reached EOF
			return nil
		}

		result := parser.ParseRegex().FindStringSubmatch(s)

		err = parser.ParseLine(result)
		if err != nil {
			// The string couldn't be parsed.
			prefix := parser.LinePrefix()
			if prefix != "" && strings.HasPrefix(s, prefix) {
				// The line has the correct prefix, but it couldn't be parsed.
				// We'll skip this line and keep going.
				fmt.Println("WARNING: Line \"" + s + "\" was skipped: " + err.Error())
				continue
			}

			// Otherwise, this still isn't error; it just means we're done parsing this regex.
			return rd.UnreadString(s)
		}
	}

	return nil
}
