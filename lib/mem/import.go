package mem

import (
	"bufio"
	"errors"
	"io"
	"regexp"
	"strconv"
	"time"
)

func Import(data io.Reader) (Mem, error) {
	reader := bufio.NewReader(data)

	mem := Mem{}
	err := parseLines(reader, headerRegex, &mem.MemHeader)
	if err != nil {
		return Mem{}, err
	}

	return mem, nil
}

var headerRegex = regexp.MustCompile(`^\s+([^:]+):\s+(.*)`)

func (header *MemHeader) Parse(result []string) error {
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
