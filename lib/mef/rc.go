package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type RCNorm struct {
	Interval mem.Column `json:"interval"`
	MatNorm  `json:"threshChange"`
	mef      *Mef
}

func (norm RCNorm) NColumns() int {
	return len(norm.mef.mems)
}

func (norm RCNorm) NRows() int {
	return len(mem.RCInterval)
}

func (norm RCNorm) Column(i int) mem.Column {
	return norm.mef.mems[i].Sections["RC"].(*mem.RecoveryCycle).YColumn
}

func (norm RCNorm) WasImputed(i int) mem.Column {
	return norm.mef.mems[i].Sections["RC"].(*mem.RecoveryCycle).WasImputed
}

func (mef *Mef) rcNorm() RCNorm {
	norm := RCNorm{
		Interval: mem.RCInterval,
		mef:      mef,
	}
	norm.MatNorm = MatrixNorm(norm)
	return norm
}
