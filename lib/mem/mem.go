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
	Settings map[string]string `json:"settings"`
}

func (mem *Mem) LabelledTable(name, subsec string) LabelledTable {
	return mem.Sections[name].LabelledTable(subsec)
}

func (mem *rawMem) AsMem() (*Mem, error) {
	trueMem := &Mem{
		Header:   mem.Header,
		Sections: make(Sections),
		Settings: mem.ExcitabilityVariables.ExcitabilitySettings,
	}

	trueMem.Sections["CD"] = newCD()
	trueMem.Sections["RC"] = newRC()
	trueMem.Sections["SR"] = newSR()
	trueMem.Sections["TE"] = newTE()
	trueMem.Sections["IV"] = newIV()
	trueMem.Sections["ExVars"] = newExVar()
	for name, sec := range trueMem.Sections {
		if err := sec.LoadFromMem(mem); err != nil {
			if _, ok := err.(MissingSection); ok {
				// It's okay if a section is missing, but it should be removed
				delete(trueMem.Sections, name)
			} else {
				return nil, err
			}
		}
	}
	return trueMem, nil
}

func (mem *rawMem) MarshalJSON() ([]byte, error) {
	trueMem, err := mem.AsMem()
	if err != nil {
		return nil, err
	}
	return json.Marshal(trueMem)
}

func (mem *Mem) UnmarshalJSON(value []byte) error {
	mem.Sections = make(Sections) // This is necessary to initialize the nil map...
	// ...and now we want to just to a regular json.Unmarshal, but that would cause recursion...
	type aliasMem *Mem // ...so create an alias...
	mem2 := aliasMem(mem)
	return json.Unmarshal(value, mem2) // ...and now the alias's default Unmarshal does what we want.
}

func Import(data io.Reader) (*Mem, error) {
	reader := NewReader(data)
	mem := rawMem{}
	mem.ExcitabilityVariables.Values = make(map[int]float64)

	err := mem.Header.Parse(reader)
	if err != nil {
		return nil, err
	}

	for err == nil {
		err = mem.importRawSection(reader)
	}

	if err != io.EOF && err != nil {
		return nil, fmt.Errorf("Error encountered at line %d: %s", reader.GetLastLineNumber(), err.Error())
	}

	return mem.AsMem()
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
