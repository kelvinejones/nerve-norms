package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type IVNorm struct {
	Interval mem.Column
	MatNorm
	mef *Mef
}

func (norm IVNorm) NColumns() int {
	return len(norm.mef.mems)
}

func (norm IVNorm) NRows() int {
	return len(mem.IVCurrent)
}

func (norm IVNorm) Column(i int) mem.Column {
	return norm.mef.mems[i].Sections["IV"].(*mem.ThresholdIV).ThreshReduction
}

func (norm IVNorm) WasImputed(i int) mem.Column {
	return norm.mef.mems[i].Sections["IV"].(*mem.ThresholdIV).WasImputed
}

func (mef *Mef) ivNorm() IVNorm {
	norm := IVNorm{
		Interval: mem.IVCurrent,
		mef:      mef,
	}
	norm.MatNorm = MatrixNorm(norm)
	return norm
}