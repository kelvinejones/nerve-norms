package mef

import (
	"math"

	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type NormTable struct {
	XValues mem.Column `json:"xvalues,omitempty"`
	Mean    mem.Column `json:"mean"`
	SD      mem.Column `json:"sd"`
	Num     mem.Column `json:"num"`
}

func NewNormTable(xv mem.Column, mef *Mef, sec, subsec string) NormTable {
	lt := mef.mems[0].LabelledTable(sec, subsec)
	numEl := lt.Len()
	norm := NormTable{
		XValues: xv,
		Mean:    make(mem.Column, numEl),
		SD:      make(mem.Column, numEl),
		Num:     make(mem.Column, numEl),
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
