package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type RCNorm struct{ GenericNorm }

func rcTable(mData *mem.Mem) *mem.LabelledTable {
	return &mData.Sections["RC"].(*mem.RecoveryCycle).LabelledTable
}

func (mef *Mef) rcNorm() RCNorm {
	norm := RCNorm{
		GenericNorm: GenericNorm{
			XValues: mem.RCInterval,
			mef:     mef,
			ltfm:    rcTable,
		},
	}
	norm.MatNorm = MatrixNorm(norm)
	return norm
}
