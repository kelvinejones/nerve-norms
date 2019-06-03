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

type transform func(float64) float64

func NewExpNormTable(xv mem.Column, mef *Mef, sec, subsec string) NormTable {
	return newNormTable(xv, mef, sec, subsec, math.Log10, func(x float64) float64 {
		return math.Pow(10, x)
	})
}

func NewNormTable(xv mem.Column, mef *Mef, sec, subsec string) NormTable {
	return newNormTable(xv, mef, sec, subsec, nil, nil)
}

func newNormTable(xv mem.Column, mef *Mef, sec, subsec string, forward, reverse transform) NormTable {
	norm := NormTable{Values: xv}
	for _, mm := range *mef {
		lt := mm.LabelledTable(sec, subsec)
		numEl := lt.Len()
		norm.Mean = make(mem.Column, numEl)
		norm.SD = make(mem.Column, numEl)
		norm.Num = make(mem.Column, numEl)
	}

	// Sum the values
	for _, mm := range *mef {
		lt := mm.LabelledTable(sec, subsec)
		for rowN := 0; rowN < lt.Len(); rowN++ {
			if !lt.WasImputedAt(rowN) {
				val := lt.YColumnAt(rowN)
				if forward != nil {
					val = forward(val)
				}
				norm.Mean[rowN] += val
				norm.Num[rowN]++
			}
		}
	}

	// Normalize to get mean
	for rowN := range norm.Mean {
		norm.Mean[rowN] /= norm.Num[rowN]
	}

	// Calculate SD
	for _, mm := range *mef {
		lt := mm.LabelledTable(sec, subsec)
		for rowN := 0; rowN < lt.Len(); rowN++ {
			if !lt.WasImputedAt(rowN) {
				val := lt.YColumnAt(rowN)
				if forward != nil {
					val = forward(val)
				}
				norm.SD[rowN] += math.Pow(val-norm.Mean[rowN], 2)
			}
		}
	}

	// Normalize to get SD
	for rowN := range norm.Mean {
		norm.SD[rowN] = math.Sqrt(norm.SD[rowN] / norm.Num[rowN])
		if reverse != nil {
			norm.SD[rowN] = reverse(norm.SD[rowN])
		}
	}

	if reverse != nil {
		for rowN := range norm.Mean {
			norm.Mean[rowN] = reverse(norm.Mean[rowN])
		}
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

type DoubleNormTable struct {
	YNorm NormTable
	XNorm NormTable
}

func (norm DoubleNormTable) MarshalJSON() ([]byte, error) {
	jt := jsonTable{
		Columns: []string{"ymean", "ysd", "ynum", "xmean", "xsd", "xnum"},
		Data:    []mem.Column{norm.YNorm.Mean, norm.YNorm.SD, norm.YNorm.Num, norm.XNorm.Mean, norm.XNorm.SD, norm.XNorm.Num},
	}

	return json.Marshal(&jt)
}

func (norm *DoubleNormTable) UnmarshalJSON(value []byte) error {
	var jt jsonTable
	err := json.Unmarshal(value, &jt)
	if err != nil {
		return err
	}

	if len(jt.Columns) != 6 || len(jt.Data) != 6 {
		return errors.New("Incorrect number of DoubleNormTable columns in JSON")
	}

	norm.YNorm.Mean = jt.Data[0]
	norm.YNorm.SD = jt.Data[1]
	norm.YNorm.Num = jt.Data[3]
	norm.XNorm.Mean = jt.Data[4]
	norm.XNorm.SD = jt.Data[5]
	norm.XNorm.Num = jt.Data[6]

	return nil
}
