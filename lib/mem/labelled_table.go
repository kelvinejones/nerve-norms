package mem

import (
	"encoding/json"
	"errors"
)

type LabelledTable struct {
	XName      string
	YName      string
	XColumn    Column
	YColumn    Column
	WasImputed Column
}

// jsonTable is used to restructure LabelledTable data for json.
type jsonTable struct {
	Columns []string `json:"columns"`
	Data    Table    `json:"data"`
}

func (lt LabelledTable) MarshalJSON() ([]byte, error) {
	jt := jsonTable{
		Columns: []string{lt.XName, lt.YName},
		Data:    []Column{lt.XColumn, lt.YColumn},
	}

	if lt.WasImputed != nil {
		jt.Columns = append(jt.Columns, "Was Imputed")
		jt.Data = append(jt.Data, lt.WasImputed)
	}

	return json.Marshal(&jt)
}

func (lt *LabelledTable) UnmarshalJSON(value []byte) error {
	var jt jsonTable
	err := json.Unmarshal(value, &jt)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	numCol := len(jt.Columns)
	numDat := len(jt.Data)

	if numCol < 2 || numCol > 3 || numDat < 2 || numDat > 3 {
		return errors.New("Incorrect number of LabelledTable columns in JSON")
	}

	lt.XName = jt.Columns[0]
	lt.YName = jt.Columns[1]
	lt.XColumn = jt.Data[0]
	lt.YColumn = jt.Data[1]

	if numCol == 3 {
		if jt.Columns[2] != "Was Imputed" {
			return errors.New("Incorrect TablelledTable column names in JSON")
		}
		lt.WasImputed = jt.Data[2]
	}

	return nil
}
