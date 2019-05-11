package mem

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type Mem struct {
	Header                `json:"header"`
	Sections              []Section `json:"sections"`
	ExcitabilityVariables `json:"exVars"`
}

func Import(data io.Reader) (Mem, error) {
	reader := NewReader(data)
	mem := Mem{}
	mem.ExcitabilityVariables.Values = make(map[string]float64)

	err := mem.Header.Parse(reader)
	if err != nil {
		return mem, err
	}

	for err == nil {
		err = mem.importSection(reader)
	}

	if err != io.EOF && err != nil {
		return mem, fmt.Errorf("Error encountered at line %d: %s", reader.GetLastLineNumber(), err.Error())
	}

	return mem, nil
}

func (mem *Mem) importSection(reader *Reader) error {
	str, err := reader.skipEmptyLines()
	if err != nil {
		return err
	}

	if len(str) < 2 || str[0] != ' ' {
		return errors.New("Could not parse invalid section: '" + str + "'")
	}

	if strings.Contains(str, "DERIVED EXCITABILITY VARIABLES") {
		return mem.ExcitabilityVariables.Parse(reader)
	}

	sec := Section{Header: strings.TrimSpace(str)}
	err = sec.parse(reader)
	if err != nil {
		return err
	}
	mem.Sections = append(mem.Sections, sec)

	return nil
}

func (mem Mem) String() string {
	str := "Mem{\n"
	str += "\t" + mem.Header.String() + ",\n"
	for _, sec := range mem.Sections {
		str += "\t" + sec.String() + ",\n"
	}
	str += "\t" + mem.ExcitabilityVariables.String() + ",\n"
	str += "}"
	return str
}

// sectionContainingHeader returns a section containing the provided header.
// Dashes are replaced with spaces for a slightly less sensitive search.
func (mem Mem) sectionContainingHeader(header string) (Section, error) {
	for _, sec := range mem.Sections {
		if strings.Contains(strings.Replace(sec.Header, "-", " ", -1), header) {
			return sec, nil
		}
	}
	return Section{}, errors.New("MEM does not contain section '" + header + "'")
}
