package mem

import (
	"encoding/json"
	"errors"
)

type LabelledTable interface {
	XName() string
	YName() string
	XColumnAt(int) float64
	YColumnAt(int) float64
	HasImputed() bool
	WasImputedAt(int) bool
	Len() int
}

type LabelledTableFromMem func(*Mem) LabelledTable

type LabTab struct {
	xname  string
	yname  string
	xcol   Column
	ycol   Column
	wasimp Column
}

func (lt LabTab) XName() string {
	return lt.xname
}

func (lt LabTab) YName() string {
	return lt.yname
}

func (lt LabTab) XColumnAt(idx int) float64 {
	return lt.xcol[idx]
}

func (lt LabTab) YColumnAt(idx int) float64 {
	return lt.ycol[idx]
}

func (lt LabTab) HasImputed() bool {
	return lt.wasimp != nil
}

func (lt LabTab) WasImputedAt(idx int) bool {
	return len(lt.wasimp) != 0 && lt.wasimp[idx] > 0.5
}

func (lt LabTab) Len() int {
	return len(lt.xcol)
}

// jsonTable is used to restructure LabTab data for json.
type jsonTable struct {
	Columns []string `json:"columns"`
	Data    Table    `json:"data"`
}

func (lt LabTab) MarshalJSON() ([]byte, error) {
	jt := jsonTable{
		Columns: []string{lt.xname, lt.yname},
		Data:    []Column{lt.xcol, lt.ycol},
	}

	if lt.wasimp != nil {
		jt.Columns = append(jt.Columns, "Was Imputed")
		jt.Data = append(jt.Data, lt.wasimp)
	}

	return json.Marshal(&jt)
}

func (lt *LabTab) UnmarshalJSON(value []byte) error {
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

	lt.xname = jt.Columns[0]
	lt.yname = jt.Columns[1]
	lt.xcol = jt.Data[0]
	lt.ycol = jt.Data[1]

	if numCol == 3 {
		if jt.Columns[2] != "Was Imputed" {
			return errors.New("Incorrect TablelledTable column names in JSON")
		}
		lt.wasimp = jt.Data[2]
	}

	return nil
}
