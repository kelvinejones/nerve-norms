package mem

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Mem struct {
	Header
	Sections []RawSection
	ExcitabilityVariables
}

func (mem *Mem) MarshalJSON() ([]byte, error) {
	str := &struct {
		Header   *Header                `json:"header"`
		Sections map[string]*RawSection `json:"sections"`
		ExVars   map[int]float64        `json:"exVars"`
		Settings map[string]string      `json:"settings"`
	}{
		Header:   &mem.Header,
		Sections: make(map[string]*RawSection),
		ExVars:   mem.ExcitabilityVariables.Values,
		Settings: mem.ExcitabilityVariables.ExcitabilitySettings,
	}
	for i, sec := range mem.Sections {
		str.Sections[sec.Abbreviation] = &(mem.Sections[i])
	}
	return json.Marshal(str)
}

// Verify verifies that all of the required data is here.
func (mem *Mem) Verify() error {
	if _, err := mem.ChargeDuration(); err != nil {
		return err
	}
	if _, err := mem.RecoveryCycle(); err != nil {
		return err
	}
	if _, err := mem.StimulusResponse(); err != nil {
		return err
	}
	if _, err := mem.ThresholdElectrotonus(); err != nil {
		return err
	}
	if _, err := mem.ThresholdIV(); err != nil {
		return err
	}
	return nil
}

func Import(data io.Reader) (Mem, error) {
	reader := NewReader(data)
	mem := Mem{}
	mem.ExcitabilityVariables.Values = make(map[int]float64)

	err := mem.Header.Parse(reader)
	if err != nil {
		return mem, err
	}

	for err == nil {
		err = mem.importRawSection(reader)
	}

	if err != io.EOF && err != nil {
		return mem, fmt.Errorf("Error encountered at line %d: %s", reader.GetLastLineNumber(), err.Error())
	}

	return mem, nil
}

func (mem *Mem) importRawSection(reader *Reader) error {
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

	sec := RawSection{Header: strings.TrimSpace(str)}
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
func (mem Mem) sectionContainingHeader(header string) (RawSection, error) {
	for _, sec := range mem.Sections {
		if strings.Contains(strings.Replace(sec.Header, "-", " ", -1), header) {
			return sec, nil
		}
	}
	return RawSection{}, errors.New("MEM does not contain section '" + header + "'")
}
