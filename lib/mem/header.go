package mem

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Sex int

const (
	UnknownSex Sex = iota
	FemaleSex
	MaleSex
)

type Header struct {
	File      string
	Name      string
	Protocol  string
	Date      time.Time
	StartTime time.Time // TODO get rid of this field; merge into Date
	Age       int
	Sex
	Temperature   float64
	SRSites       string
	NormalControl bool
	Operator      string
	Comment       string
}

func (header *Header) Parse(reader *Reader) error {
	return reader.parseLines(header)
}

func (header Header) String() string {
	return "Header{File{\"" + header.File + "\"}, Name{\"" + header.Name + "\"} }"
}

func (header Header) ParseRegex() *regexp.Regexp {
	return regexp.MustCompile(`^\s+([^:]+):\s*(.*)`)
}

func (header *Header) ParseLine(result []string) error {
	if len(result) != 3 {
		return errors.New("Incorrect header line length")
	}

	var err error
	val := strings.TrimSpace(result[2])
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
		temp, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		header.Temperature = temp
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
