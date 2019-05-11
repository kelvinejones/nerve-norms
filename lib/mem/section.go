package mem

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Column []float64
type Table []Column
type TableSet struct {
	// ColCount is the number of columns in Names and in each Table.
	ColCount int

	// Names are the names of the columns.
	Names []string

	// Tables is a slice of tables (usually just one).
	Tables []Table
}

func (ts TableSet) MarshalJSON() ([]byte, error) {
	data := struct {
		Names []string        `json:"header"`
		Data  json.RawMessage `json:"data"`
	}{}
	data.Names = ts.Names

	var err error
	switch len(ts.Tables) {
	case 0:
		// Do nothing; it's empty
	case 1:
		data.Data, err = json.Marshal(ts.Tables[0])
	default:
		data.Data, err = json.Marshal(ts.Tables)
	}
	if err != nil {
		return nil, err
	}

	return json.Marshal(&data)
}

func (tab Table) MarshalJSON() ([]byte, error) {
	numCols := len(tab)
	if numCols == 0 {
		return []byte(`[[]]`), nil
	}
	numRows := len(tab[0])
	data := make([][]float64, numRows) // It's rows of columns of floats
	for i := range tab[0] {
		data[i] = make([]float64, numCols)
	}

	// Go through the length of the first column (assuming all columns are the same length)
	for colNum, col := range tab {
		for rowNum, val := range col {
			data[rowNum][colNum] = val
		}
	}
	return json.Marshal(&data)
}

type Section struct {
	// Header is the header for the section.
	Header string

	// TableSet is usually just one table, but it might be many
	TableSet

	// ExtraLines are extra lines which couldn't be parsed (e.g. Max CMAP).
	ExtraLines []string
}

func (sec *Section) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name  *string  `json:"name"`
		Data  TableSet `json:"data"`
		Extra []string `json:"extra,omitempty"`
	}{
		Name:  &sec.Header,
		Data:  sec.TableSet,
		Extra: sec.ExtraLines,
	})
}

// columnContainsName returns the first column containing the provided name.
func (sec Section) columnContainsName(name string, table int) (Column, error) {
	if table > len(sec.Tables) {
		return Column{}, errors.New("Attempt to access table out of range in section '" + sec.Header + "'")
	}

	for i, str := range sec.Names {
		if strings.Contains(str, name) {
			return sec.Tables[table][i], nil
		}
	}

	return Column{}, errors.New("Column '" + name + "' was not found in '" + sec.Header + "'")
}

func (tab *Table) appendRow(row []string) error {
	// By this point we know the number of columns matches
	for i, str := range row {
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return errors.New("Could not parse float: '" + str + "'")
		}

		(*tab)[i] = append((*tab)[i], val)
	}
	return nil
}

func (ts *TableSet) appendRow(row []string) error {
	// Format is two characters followed by optional digits, a decimal, and digits
	result := regexp.MustCompile(`^[[:alpha:]]{2}(\d*)\.(\d+)`).FindStringSubmatch(row[0])
	if len(result) != 3 {
		return errors.New("A table row must contain a valid location: '" + row[0] + "'")
	}

	tableNum := 1
	var err error
	if result[1] != "" {
		tableNum, err = strconv.Atoi(result[1])
		if err != nil {
			return errors.New("Table number could not be parsed: " + err.Error())
		}
	}

	// Parse the row number to insure it's valid, but we don't use it
	_, err = strconv.Atoi(result[2])
	if err != nil {
		return errors.New("Row number could not be parsed: " + err.Error())
	}

	for len(ts.Tables) < tableNum {
		ts.Tables = append(ts.Tables, make([]Column, ts.ColCount))
	}

	return ts.Tables[tableNum-1].appendRow(row[1:])
}

func (sec *Section) parse(reader *Reader) error {
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

func (sec Section) String() string {
	return fmt.Sprintf("Section{'%s', %v}", sec.Header, sec.TableSet)
}

func (ts TableSet) String() string {
	if ts.ColCount == 0 {
		return "TableSet{}"
	}
	if ts.Tables == nil || len(ts.Tables) == 0 {
		return fmt.Sprintf("TableSet{empty, %d columns}", ts.ColCount)
	}

	numTables := len(ts.Tables)
	if numTables > 1 {
		numRows := 0
		for _, tab := range ts.Tables {
			if len(tab) == 0 {
				return fmt.Sprintf("TableSet{%d tables, %d x ?}", numTables, ts.ColCount)
			}
			numRows += len(tab[0])
		}
		return fmt.Sprintf("TableSet{%d tables, %dx%d stacked}", numTables, ts.ColCount, numRows)
	}

	// There is only one table
	if len(ts.Tables[0]) == 0 {
		return fmt.Sprintf("TableSet{%d x ?}", ts.ColCount)
	}
	return fmt.Sprintf("TableSet{%dx%d}", ts.ColCount, len(ts.Tables[0][0]))
}

func (col Column) Maximum() float64 {
	max := math.Inf(-1)
	for _, val := range col {
		if val > max {
			max = val
		}
	}
	return max
}

func (col Column) Minimum() float64 {
	min := math.Inf(1)
	for _, val := range col {
		if val < min {
			min = val
		}
	}
	return min
}
