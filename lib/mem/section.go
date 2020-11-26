package mem

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type loadableSection interface {
	LoadFromMem(mem *rawMem) error
	Section
}

type Section interface {
	LabelledTable(string) LabelledTable
}

type Sections map[string]Section

func (secs *Sections) UnmarshalJSON(value []byte) error {
	rawSecs := map[string]json.RawMessage{}
	err := json.Unmarshal(value, &rawSecs)
	if err != nil {
		return err
	}

	for key, val := range rawSecs {
		var sec Section
		switch key {
		case "CD":
			sec = &ChargeDuration{}
			err = json.Unmarshal(val, sec)
		case "RC":
			sec = &RecoveryCycle{}
			err = json.Unmarshal(val, sec)
		case "SR":
			sec = &StimResponse{}
			err = json.Unmarshal(val, sec)
		case "TE":
			sec = &ThresholdElectrotonus{}
			err = json.Unmarshal(val, sec)
		case "IV":
			sec = &ThresholdIV{}
			err = json.Unmarshal(val, sec)
		case "ExVars":
			sec = &ExcitabilityVariablesSection{}
			err = json.Unmarshal(val, sec)
		}
		if err != nil {
			return err
		}
		(*secs)[key] = sec
	}

	return nil
}

//RawSection A section of text defined by what type of variables they are
type RawSection struct {
	// Header is the header for the section.
	Header string

	// TableSet is usually just one table, but it might be many
	TableSet

	// ExtraLines are extra lines which couldn't be parsed (e.g. Max CMAP).
	ExtraLines []string
}

func (sec *RawSection) MarshalJSON() ([]byte, error) {
	str := &struct {
		Name    *string         `json:"name"`
		Columns []string        `json:"columnNames"`
		Data    json.RawMessage `json:"data"`
		Extra   []string        `json:"extra,omitempty"`
	}{
		Name:    &sec.Header,
		Columns: sec.TableSet.Names,
		Extra:   sec.ExtraLines,
	}

	var err error
	switch len(sec.TableSet.Tables) {
	case 0:
		// Do nothing; it's empty
	case 1:
		str.Data, err = json.Marshal(sec.TableSet.Tables[0])
	default:
		str.Data, err = json.Marshal(sec.TableSet.Tables)
	}
	if err != nil {
		return nil, err
	}

	return json.Marshal(&str)
}

// columnContainsName returns the first column containing the provided name.
func (sec RawSection) columnContainsName(name string, table int) (Column, error) {
	if table >= len(sec.Tables) {
		return Column{}, errors.New("Attempt to access table out of range in section '" + sec.Header + "'")
	}

	for i, str := range sec.Names {
		if strings.Contains(str, name) {
			return sec.Tables[table][i], nil
		}
	}

	return Column{}, errors.New("Column '" + name + "' was not found in '" + sec.Header + "'")
}

func (sec *RawSection) parse(reader *Reader) error {
	// Keep parsing extra lines until we get a valid table header
	for sec.ColCount == 0 {
		str, err := reader.skipEmptyLines()
		if err != nil {
			return err
		}

		if strings.Contains(str, "\t") {
			// A tab indicates it's a table header. I hope.
			sec.Names = splitColumns(str)
			if !rowIsHeader(sec.Names) {
				return errors.New("Could not parse header row: '" + str + "'")
			}
			sec.Names = sec.Names[1:] // Delete the empty first column
			sec.ColCount = len(sec.Names)
		} else {
			sec.ExtraLines = append(sec.ExtraLines, strings.TrimSpace(str))
		}
	}

	// Now that there's a valid table header, parse the remaining lines
	str, err := reader.skipEmptyLines()
	cols := splitColumns(str)
	for err == nil {
		// Parse a line
		if len(cols) != sec.ColCount+1 {
			// Row doesn't have correct number of columns, so assume we don't parse it
			break
		}
		sec.TableSet.appendRow(cols)

		str, err = reader.skipEmptyLines()
		cols = splitColumns(str)
	}
	if err != nil {
		return err
	}

	// The most recent line isn't what we want. Put it back.
	reader.UnreadString(str)

	return nil
}

func splitColumns(str string) []string {
	columns := strings.Split(str, "\t")
	for i, col := range columns {
		columns[i] = strings.TrimSpace(col)
	}
	return columns
}

func rowIsHeader(cols []string) bool {
	// The statement is ordered to prevent panics while checking for two things: first column is empty string, and there are 2+ columns.
	return !(len(cols) < 1 || strings.TrimSpace(cols[0]) != "" || len(cols) < 2)
}

func (sec RawSection) String() string {
	return fmt.Sprintf("RawSection{'%s', %v}", sec.Header, sec.TableSet)
}
