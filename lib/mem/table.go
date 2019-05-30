package mem

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type Table []Column
type TableSet struct {
	// ColCount is the number of columns in Names and in each Table.
	ColCount int

	// Names are the names of the columns.
	Names []string

	// Tables is a slice of tables (usually just one).
	Tables []Table

	// Abbreviation is the prefix of each import row.
	Abbreviation string
}

// MarshalJSON swaps the order of rows and columns.
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

// UnmarshalJSON swaps the order of rows and columns.
func (tab *Table) UnmarshalJSON(value []byte) error {
	rows := []Column{}
	err := json.Unmarshal(value, &rows)
	if err != nil {
		return err
	}

	numRows := len(rows)
	if numRows == 0 {
		// Empty table
		return nil
	}

	for i := range rows[0] {
		thisCol := make(Column, numRows)
		*tab = append(*tab, thisCol)
		for j, val := range rows {
			thisCol[j] = val[i]
		}
	}

	return nil
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
	result := regexp.MustCompile(`^([[:alpha:]]{2})(\d*)\.(\d+)`).FindStringSubmatch(row[0])
	if len(result) != 4 {
		return errors.New("A table row must contain a valid location: '" + row[0] + "'")
	}

	if ts.Abbreviation != result[1] {
		if ts.Abbreviation == "" {
			ts.Abbreviation = result[1]
		} else {
			return errors.New("The table's rows don't have matching prefixes: '" + ts.Abbreviation + "' and '" + result[1] + "'")
		}
	}

	tableNum := 1
	var err error
	if result[2] != "" {
		tableNum, err = strconv.Atoi(result[2])
		if err != nil {
			return errors.New("Table number could not be parsed: " + err.Error())
		}
	}

	// Parse the row number to insure it's valid, but we don't use it
	_, err = strconv.Atoi(result[3])
	if err != nil {
		return errors.New("Row number could not be parsed: " + err.Error())
	}

	for len(ts.Tables) < tableNum {
		ts.Tables = append(ts.Tables, make([]Column, ts.ColCount))
	}

	return ts.Tables[tableNum-1].appendRow(row[1:])
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
