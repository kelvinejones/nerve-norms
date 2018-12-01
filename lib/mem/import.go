package mem

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Import(data io.Reader) (Mem, error) {
	reader := bufio.NewReader(data)
	mem := Mem{}
	var err error

	err = parseHeader(reader, &mem.Header)
	if err != nil {
		return mem, err
	}

	err = parseStimResponse(reader, &mem.StimResponse)
	if err != nil {
		return mem, err
	}

	return mem, nil
}

func skipNewlines(reader *bufio.Reader) (string, error) {
	s := "\n"
	var err error
	for s == "\n" && err == nil {
		s, err = reader.ReadString('\n')
	}

	return s, err
}

func skipUntilContains(reader *bufio.Reader, search string) (string, error) {
	s, err := skipNewlines(reader)
	if err == nil && !strings.Contains(s, search) {
		err = errors.New("Could not find '" + search + "'")
	}
	return s, err
}

func parseHeader(reader *bufio.Reader, header *Header) error {
	return parseLines(reader, headerRegex, header)
}

func parseStimResponse(reader *bufio.Reader, sr *StimResponse) error {
	// Find section header
	s, err := skipUntilContains(reader, "STIMULUS-RESPONSE DATA")
	if err != nil {
		return err
	}

	// Find some random string that's there
	s, err = skipUntilContains(reader, "Values are those recorded")
	if err != nil {
		return err
	}

	// Find Max CMAP
	s, err = skipNewlines(reader)
	if err != nil {
		return err
	}
	n, err := fmt.Sscanf(s, " Max CMAP  1 ms =  %f mV\n", &sr.MaxCmap)
	if n != 1 || err != nil {
		return errors.New("Could not find Max CMAP: " + s)
	}

	s, err = skipUntilContains(reader, "% Max               	Stimulus")
	if err != nil {
		return err
	}

	return nil
}

var headerRegex = regexp.MustCompile(`^\s+([^:]+):\s+(.*)`)

func (header *Header) Parse(result []string) error {
	if len(result) != 3 {
		return errors.New("Incorrect header line length")
	}

	var err error
	val := result[2]
	switch result[1] {
	case "NC/disease":
		if val == "NC" {
			header.NormalControl = true
		}
		// TODO update for other options? Currently disease is the default, which I suppose excludes an uncertain ones from the control database
	case "Sex":
		switch val {
		case "M":
			header.Sex = MaleSex
		case "F":
			header.Sex = FemaleSex
		default:
			header.Sex = UnknownSex
		}
	case "Temperature":
		temp, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return err
		}
		header.Temperature = float32(temp)
	case "Age":
		header.Age, err = strconv.Atoi(val)
		if err != nil {
			return err
		}
	case "Date":
		layout := "2/1/06"
		header.Date, err = time.Parse(layout, val)
		if err != nil {
			return err
		}
	case "Start time":
		layout := "2/1/06 15:04:05"
		header.StartTime, err = time.Parse(layout, "2/1/06 "+val)
		if err != nil {
			return err
		}
	case "File":
		header.File = val
	case "Name":
		header.Name = val
	case "Protocol":
		header.Protocol = val
	case "S/R sites":
		header.SRSites = val
	case "Operator":
		header.Operator = val
	case "Comments":
		header.Comment = val
	}

	return nil
}

type Parser interface {
	Parse([]string) error
}

func parseLines(reader *bufio.Reader, regex *regexp.Regexp, parser Parser) error {
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if len(s) == 1 {
			// Done with section; break!
			break
		}
		result := headerRegex.FindStringSubmatch(s)

		err = parser.Parse(result)
		if err != nil {
			return errors.New(err.Error() + ": '" + s + "'")
		}
	}

	return nil
}
