package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type CDNorm struct{ GenericNorm }

func cdTable(mData *mem.Mem) *mem.LabelledTable {
	return &mData.Sections["CD"].(*mem.ChargeDuration).LabelledTable
}

func (mef *Mef) cdNorm() CDNorm {
	norm := CDNorm{
		GenericNorm: GenericNorm{
			XValues: mem.CDDuration,
			mef:     mef,
			ltfm:    cdTable,
		},
	}
	norm.MatNorm = MatrixNorm(norm)
	return norm
}
