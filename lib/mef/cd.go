package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type CDNorm struct {
	Duration mem.Column `json:"duration"`
	MatNorm  `json:"threshCharge"`
	mef      *Mef
}

func (norm CDNorm) NColumns() int {
	return len(norm.mef.mems)
}

func (norm CDNorm) NRows() int {
	return len(mem.CDDuration)
}

func (norm CDNorm) Column(i int) mem.Column {
	return norm.mef.mems[i].Sections["CD"].(*mem.ChargeDuration).ThreshCharge
}

func (norm CDNorm) WasImputed(i int) mem.Column {
	return norm.mef.mems[i].Sections["CD"].(*mem.ChargeDuration).WasImputed
}

func (mef *Mef) cdNorm() CDNorm {
	norm := CDNorm{
		Duration: mem.CDDuration,
		mef:      mef,
	}
	norm.MatNorm = MatrixNorm(norm)
	return norm
}
