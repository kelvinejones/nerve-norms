package mef

import (
	"math"

	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type LabelledTableFromMem func(*mem.Mem) *mem.LabelledTable

type GenericNorm struct {
	XValues mem.Column `json:"xvalues,omitempty"`
	MatNorm `json:"norms"`
	ltfm    LabelledTableFromMem
	mef     *Mef
}

type MatNorm struct {
	Mean mem.Column `json:"mean"`
	SD   mem.Column `json:"sd"`
	Num  mem.Column `json:"num"`
}

func (mat *GenericNorm) CalculateNorms() {
	numEl := len(mat.ltfm(mat.mef.mems[0]).XColumn)
	mat.MatNorm.Mean = make(mem.Column, numEl)
	mat.MatNorm.SD = make(mem.Column, numEl)
	mat.MatNorm.Num = make(mem.Column, numEl)

	// Sum the values
	for colN, mm := range mat.mef.mems {
		col := mat.ltfm(mm).YColumn
		for rowN := range col {
			if !mat.wasImp(colN, rowN) {
				mat.MatNorm.Mean[rowN] += col[rowN]
				mat.MatNorm.Num[rowN]++
			}
		}
	}

	// Normalize to get mean
	for rowN := range mat.MatNorm.Mean {
		mat.MatNorm.Mean[rowN] /= mat.MatNorm.Num[rowN]
	}

	// Calculate SD
	for colN, mm := range mat.mef.mems {
		col := mat.ltfm(mm).YColumn
		for rowN := range col {
			if !mat.wasImp(colN, rowN) {
				mat.MatNorm.SD[rowN] += math.Pow(col[rowN]-mat.MatNorm.Mean[rowN], 2)
			}
		}
	}

	// Normalize to get SD
	for rowN := range mat.MatNorm.Mean {
		mat.MatNorm.SD[rowN] = math.Sqrt(mat.MatNorm.SD[rowN] / mat.MatNorm.Num[rowN])
	}
}

func (mat *GenericNorm) wasImp(colN, rowN int) bool {
	col := mat.ltfm(mat.mef.mems[colN]).WasImputed
	// Yes, this is terrible, but wasImputed is a float column
	return len(col) != 0 && col[rowN] > 0.5
}
