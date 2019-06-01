package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type CDNorm struct {
	Duration mem.Column `json:"duration"`
	GenericNorm
}

func cdTable(mData *mem.Mem) *mem.LabelledTable {
	return &mData.Sections["CD"].(*mem.ChargeDuration).LabelledTable
}

func (mef *Mef) cdNorm() CDNorm {
	norm := CDNorm{
		Duration: mem.CDDuration,
		GenericNorm: GenericNorm{
			mef:  mef,
			ltfm: cdTable,
		},
	}
	norm.MatNorm = MatrixNorm(norm)
	return norm
}
