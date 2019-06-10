package mef

import (
	"encoding/json"
	"errors"
	"math"

	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type MeanType int

const (
	ArithmeticMean MeanType = iota
	GeometricMean
)

type NormTable struct {
	Values mem.Column // Usually these are the x-values for a y-mean.
	Mean   mem.Column
	SD     mem.Column
	Num    mem.Column
	MeanType
	sec, subsec string
}

func NewNormTable(xv mem.Column, mef *Mef, sec, subsec string, mt MeanType) NormTable {
	norm := NormTable{
		Values:   xv,
		MeanType: mt,
		sec:      sec,
		subsec:   subsec,
	}
	numEl := 0
	for _, mm := range *mef {
		numEl = mm.LabelledTable(norm.sec, norm.subsec).Len()
		if numEl != 0 {
			norm.Mean = make(mem.Column, numEl)
			norm.SD = make(mem.Column, numEl)
			norm.Num = make(mem.Column, numEl)
			break // We really only need to do this once
		}
	}
	if numEl == 0 {
		// If none of the MEM have this value, then we can't construct norms for it.
		return NormTable{}
	}

	// Sum the values
	for _, mm := range *mef {
		lt := mm.LabelledTable(norm.sec, norm.subsec)
		for rowN := 0; rowN < numEl; rowN++ {
			if !lt.WasImputedAt(rowN) {
				norm.Mean[rowN] += norm.forward(lt.YColumnAt(rowN))
				norm.Num[rowN]++
			}
		}
	}

	// Normalize to get mean
	for rowN := range norm.Mean {
		if norm.Num[rowN] == 0 {
			norm.Mean[rowN] = 0
		} else {
			norm.Mean[rowN] /= norm.Num[rowN]
		}
	}

	// Calculate SD
	for _, mm := range *mef {
		lt := mm.LabelledTable(norm.sec, norm.subsec)
		for rowN := 0; rowN < numEl; rowN++ {
			if !lt.WasImputedAt(rowN) {
				norm.SD[rowN] += math.Pow(norm.forward(lt.YColumnAt(rowN))-norm.Mean[rowN], 2)
			}
		}
	}

	// Normalize to get SD
	for rowN := range norm.Mean {
		if norm.Num[rowN] == 0 {
			norm.SD[rowN] = 0
		} else {
			norm.SD[rowN] /= norm.Num[rowN]
		}
		norm.SD[rowN] = norm.reverse(math.Sqrt(norm.SD[rowN]))
		norm.Mean[rowN] = norm.reverse(norm.Mean[rowN])
	}

	return norm
}

func (nt NormTable) LabelledTable(string) mem.LabelledTable {
	wasimp := make([]float64, len(nt.Num))
	hasimp := false
	for i := range nt.Num {
		if nt.Num[i] == 0 {
			wasimp[i] = 1.0
			hasimp = true
		}
	}
	if !hasimp {
		wasimp = nil
	}
	return mem.NewLabTab("", "", nt.Values, nt.Mean, wasimp)
}

func (norm NormTable) forward(val float64) float64 {
	switch norm.MeanType {
	case ArithmeticMean:
		return val
	case GeometricMean:
		return math.Log10(val)
	default:
		return val
	}
}

func (norm NormTable) reverse(val float64) float64 {
	switch norm.MeanType {
	case ArithmeticMean:
		return val
	case GeometricMean:
		return math.Pow(10, val)
	default:
		return val
	}
}

// normJsonTable is used to restructure LabTab data for json.
type normJsonTable struct {
	Columns []string  `json:"columns"`
	Data    mem.Table `json:"data"`
}

func (norm NormTable) MarshalJSON() ([]byte, error) {
	jt := normJsonTable{
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
	var jt normJsonTable
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

type SRNormTable struct {
	YNorm NormTable
	XNorm NormTable
}

func (snt SRNormTable) LabelledTable(string) mem.LabelledTable {
	wasimp := make([]float64, len(snt.XNorm.Num))
	hasimp := false
	for i := range snt.XNorm.Num {
		if snt.XNorm.Num[i] == 0 {
			wasimp[i] = 1.0
			hasimp = true
		}
	}
	if !hasimp {
		wasimp = nil
	}
	return mem.NewLabTab("", "", mem.SRPercentMax, snt.XNorm.Mean, wasimp)
}

func (norm SRNormTable) MarshalJSON() ([]byte, error) {
	jt := normJsonTable{
		Columns: []string{"ymean", "ysd", "ynum", "xmean", "xsd", "xnum"},
		Data:    []mem.Column{norm.YNorm.Mean, norm.YNorm.SD, norm.YNorm.Num, norm.XNorm.Mean, norm.XNorm.SD, norm.XNorm.Num},
	}

	return json.Marshal(&jt)
}

func (norm *SRNormTable) UnmarshalJSON(value []byte) error {
	var jt normJsonTable
	err := json.Unmarshal(value, &jt)
	if err != nil {
		return err
	}

	if len(jt.Columns) != 6 || len(jt.Data) != 6 {
		return errors.New("Incorrect number of SRNormTable columns in JSON")
	}

	norm.YNorm.Mean = jt.Data[0]
	norm.YNorm.SD = jt.Data[1]
	norm.YNorm.Num = jt.Data[3]
	norm.XNorm.Mean = jt.Data[4]
	norm.XNorm.SD = jt.Data[5]
	norm.XNorm.Num = jt.Data[6]

	return nil
}

type TENormTables map[string]NormTable

func NewTENormTables(mef *Mef) TENormTables {
	norm := TENormTables{}
	for _, name := range []string{"h40", "h20", "d40", "d20"} {
		nt := NewNormTable(mem.TEDelay(name), mef, "TE", name, ArithmeticMean)
		if nt.Values != nil {
			// We only add this TE type of it's not zero
			norm[name] = nt
		}
	}
	return norm
}

func (tent TENormTables) LabelledTable(subsec string) mem.LabelledTable {
	return tent[subsec].LabelledTable("")
}
