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

func (norm GenericNorm) NColumns() int {
	return len(norm.mef.mems)
}

func (norm GenericNorm) NRows() int {
	return len(norm.ltfm(norm.mef.mems[0]).XColumn)
}

func (norm GenericNorm) Column(i int) mem.Column {
	return norm.ltfm(norm.mef.mems[i]).YColumn
}

func (norm GenericNorm) WasImputed(i int) mem.Column {
	return norm.ltfm(norm.mef.mems[i]).WasImputed
}

type MatNorm struct {
	Mean mem.Column `json:"mean"`
	SD   mem.Column `json:"sd"`
	Num  mem.Column `json:"num"`
}

func (mat *GenericNorm) CalculateNorms() {
	numEl := mat.NRows()
	mat.MatNorm.Mean = make(mem.Column, numEl)
	mat.MatNorm.SD = make(mem.Column, numEl)
	mat.MatNorm.Num = make(mem.Column, numEl)

	// Sum the values
	for colN := 0; colN < mat.NColumns(); colN++ {
		col := mat.Column(colN)
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
	for colN := 0; colN < mat.NColumns(); colN++ {
		col := mat.Column(colN)
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
	col := mat.WasImputed(colN)
	// Yes, this is terrible, but wasImputed is a float column
	return len(col) != 0 && col[rowN] > 0.5
}
