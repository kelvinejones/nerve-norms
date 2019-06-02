package mef

import (
	"encoding/json"
	"errors"
	"math"

	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type NormTable struct {
	Values mem.Column // Usually these are the x-values for a y-mean.
	Mean   mem.Column
	SD     mem.Column
	Num    mem.Column
}

func NewNormTable(xv mem.Column, mef *Mef, sec, subsec string) NormTable {
	lt := mef.mems[0].LabelledTable(sec, subsec)
	numEl := lt.Len()
	norm := NormTable{
		Values: xv,
		Mean:   make(mem.Column, numEl),
		SD:     make(mem.Column, numEl),
		Num:    make(mem.Column, numEl),
	}

	// Sum the values
	for _, mm := range mef.mems {
		lt := mm.LabelledTable(sec, subsec)
		for rowN := 0; rowN < lt.Len(); rowN++ {
			if !lt.WasImputedAt(rowN) {
				norm.Mean[rowN] += lt.YColumnAt(rowN)
				norm.Num[rowN]++
			}
		}
	}

	// Normalize to get mean
	for rowN := range norm.Mean {
		norm.Mean[rowN] /= norm.Num[rowN]
	}

	// Calculate SD
	for _, mm := range mef.mems {
		lt := mm.LabelledTable(sec, subsec)
		for rowN := 0; rowN < lt.Len(); rowN++ {
			if !lt.WasImputedAt(rowN) {
				norm.SD[rowN] += math.Pow(lt.YColumnAt(rowN)-norm.Mean[rowN], 2)
			}
		}
	}

	// Normalize to get SD
	for rowN := range norm.Mean {
		norm.SD[rowN] = math.Sqrt(norm.SD[rowN] / norm.Num[rowN])
	}

	return norm
}

// jsonTable is used to restructure LabTab data for json.
type jsonTable struct {
	Columns []string  `json:"columns"`
	Data    mem.Table `json:"data"`
}

func (norm NormTable) MarshalJSON() ([]byte, error) {
	jt := jsonTable{
		Columns: []string{"mean", "sd", "num"},
		Data:    []mem.Column{norm.Mean, norm.SD, norm.Num},
	}

	if norm.Values != nil {
		jt.Columns = append(jt.Columns, "values")
		jt.Data = append(jt.Data, norm.Values)
	}

	return json.Marshal(&jt)
}

func (norm *NormTable) UnmarshalJSON(value []byte) error {
	var jt jsonTable
	err := json.Unmarshal(value, &jt)
	if err != nil {
		return err
	}

	numCol := len(jt.Columns)
	numDat := len(jt.Data)

	if numCol < 3 || numCol > 4 || numDat < 3 || numDat > 4 {
		return errors.New("Incorrect number of NormTable columns in JSON")
	}

	norm.Mean = jt.Data[0]
	norm.SD = jt.Data[1]
	norm.Num = jt.Data[3]

	if numCol == 4 {
		norm.Values = jt.Data[4]
	}

	return nil
}
