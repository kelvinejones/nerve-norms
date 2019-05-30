package mem

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type rawMem struct {
	Header
	Sections []RawSection
	ExcitabilityVariables
}

type Mem struct {
	Header   `json:"header"`
	Sections `json:"sections"`
	ExVars   ExcitabilityVariables `json:"exVars"`
	Settings map[string]string     `json:"settings"`
}

func (mem *rawMem) MarshalJSON() ([]byte, error) {
	str := &Mem{
		Header:   mem.Header,
		Sections: make(Sections),
		ExVars:   mem.ExcitabilityVariables,
		Settings: mem.ExcitabilityVariables.ExcitabilitySettings,
	}

	str.Sections["CD"] = &ChargeDuration{}
	str.Sections["RC"] = &RecoveryCycle{}
	str.Sections["SR"] = &StimResponse{}
	str.Sections["TE"] = &ThresholdElectrotonus{}
	str.Sections["IV"] = &ThresholdIV{}
	for _, sec := range str.Sections {
		if err := sec.LoadFromMem(mem); err != nil {
			return nil, err
		}
	}

	return json.Marshal(str)
}

func (mem *Mem) UnmarshalJSON(value []byte) error {
	mem.Sections = make(Sections) // This is necessary to initialize the nil map...
	// ...and now we want to just to a regular json.Unmarshal, but that would cause recursion...
	type aliasMem *Mem // ...so create an alias...
	mem2 := aliasMem(mem)
	return json.Unmarshal(value, mem2) // ...and now the alias's default Unmarshal does what we want.
}

func Import(data io.Reader) (rawMem, error) {
	reader := NewReader(data)
	mem := rawMem{}
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

func (mem *rawMem) importRawSection(reader *Reader) error {
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

func (mem rawMem) String() string {
	str := "rawMem{\n"
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
func (mem rawMem) sectionContainingHeader(header string) (RawSection, error) {
	for _, sec := range mem.Sections {
		if strings.Contains(strings.Replace(sec.Header, "-", " ", -1), header) {
			return sec, nil
		}
	}
	return RawSection{}, errors.New("MEM does not contain section '" + header + "'")
}
