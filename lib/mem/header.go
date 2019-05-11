package mem

import (
	"errors"
	"fmt"
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
	File          string    `json:"file"`
	Name          string    `json:"name"`
	Protocol      string    `json:"protocol"`
	Date          time.Time `json:"date"`
	StartTime     time.Time // TODO get rid of this field; merge into Date
	Age           int       `json:"age"`
	Sex           `json:"sex"`
	Temperature   float64 `json:"temperature"`
	SRSites       string  `json:"srSites"`
	NormalControl bool    `json:"normalControl"`
	Operator      string  `json:"operator"`
	Comment       string  `json:"comment"`
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
		header.Temperature, err = strconv.ParseFloat(extractPotentialDoubleValue(val), 64)
		if err != nil {
			fmt.Println("WARNING: Line \"" + result[0] + "\" may have imported incorrectly: " + err.Error())
		}
	case "Age":
		header.Age, err = strconv.Atoi(extractPotentialDoubleValue(val))
		if err != nil {
			fmt.Println("WARNING: Line \"" + result[0] + "\" may have imported incorrectly: " + err.Error())
		}
	case "Date":
		layout := "2/1/06"
		header.Date, err = time.Parse(layout, extractPotentialDoubleValue(val))
		if err != nil {
			fmt.Println("WARNING: Line \"" + result[0] + "\" may have imported incorrectly: " + err.Error())
		}
	case "Start time":
		layout := "2/1/06 15:04:05"
		header.StartTime, err = time.Parse(layout, "2/1/06 "+extractPotentialDoubleValue(val))
		if err != nil {
			fmt.Println("WARNING: Line \"" + result[0] + "\" may have imported incorrectly: " + err.Error())
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

// extractPotentialDoubleValue will check if the input string was two parameters separated by " / ". If so, only the first is used.
func extractPotentialDoubleValue(str string) string {
	if strings.Contains(str, " / ") {
		strs := strings.SplitN(str, " / ", 2)
		return strs[0]
	} else {
		// There's no slash splitting it.
		return str
	}
}
