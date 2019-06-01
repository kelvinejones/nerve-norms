package mef

import (
	"math"

	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type LabelledTableFromMem func(*mem.Mem) *mem.LabelledTable

type GenericNorm struct {
	XValues mem.Column `json:"xvalues,omitempty"`
	Mean    mem.Column `json:"mean"`
	SD      mem.Column `json:"sd"`
	Num     mem.Column `json:"num"`
	ltfm    LabelledTableFromMem
	mef     *Mef
}

func (mat *GenericNorm) CalculateNorms() {
	numEl := len(mat.ltfm(mat.mef.mems[0]).XColumn)
	mat.Mean = make(mem.Column, numEl)
	mat.SD = make(mem.Column, numEl)
	mat.Num = make(mem.Column, numEl)

	// Sum the values
	for colN, mm := range mat.mef.mems {
		col := mat.ltfm(mm).YColumn
		for rowN := range col {
			if !mat.wasImp(colN, rowN) {
				mat.Mean[rowN] += col[rowN]
				mat.Num[rowN]++
			}
		}
	}

	// Normalize to get mean
	for rowN := range mat.Mean {
		mat.Mean[rowN] /= mat.Num[rowN]
	}

	// Calculate SD
	for colN, mm := range mat.mef.mems {
		col := mat.ltfm(mm).YColumn
		for rowN := range col {
			if !mat.wasImp(colN, rowN) {
				mat.SD[rowN] += math.Pow(col[rowN]-mat.Mean[rowN], 2)
			}
		}
	}

	// Normalize to get SD
	for rowN := range mat.Mean {
		mat.SD[rowN] = math.Sqrt(mat.SD[rowN] / mat.Num[rowN])
	}
}

func (mat *GenericNorm) wasImp(colN, rowN int) bool {
	col := mat.ltfm(mat.mef.mems[colN]).WasImputed
	// Yes, this is terrible, but wasImputed is a float column
	return len(col) != 0 && col[rowN] > 0.5
}
