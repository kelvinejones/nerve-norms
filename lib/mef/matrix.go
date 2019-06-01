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
}

func NewGenericNorm(xv mem.Column, ltfm LabelledTableFromMem, mef *Mef) GenericNorm {
	numEl := len(ltfm(mef.mems[0]).XColumn)
	norm := GenericNorm{
		XValues: xv,
		Mean:    make(mem.Column, numEl),
		SD:      make(mem.Column, numEl),
		Num:     make(mem.Column, numEl),
	}

	// Sum the values
	for _, mm := range mef.mems {
		col := ltfm(mm).YColumn
		for rowN := range col {
			if !norm.wasImp(ltfm(mm), rowN) {
				norm.Mean[rowN] += col[rowN]
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
		col := ltfm(mm).YColumn
		for rowN := range col {
			if !norm.wasImp(ltfm(mm), rowN) {
				norm.SD[rowN] += math.Pow(col[rowN]-norm.Mean[rowN], 2)
			}
		}
	}

	// Normalize to get SD
	for rowN := range norm.Mean {
		norm.SD[rowN] = math.Sqrt(norm.SD[rowN] / norm.Num[rowN])
	}

	return norm
}

func (norm *GenericNorm) wasImp(lt *mem.LabelledTable, rowN int) bool {
	col := lt.WasImputed
	// Yes, this is terrible, but wasImputed is a float column
	return len(col) != 0 && col[rowN] > 0.5
}
