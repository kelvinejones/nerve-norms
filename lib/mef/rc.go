package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type RCNorm struct {
	Interval mem.Column `json:"interval"`
	GenericNorm
}

func rcTable(mData *mem.Mem) *mem.LabelledTable {
	return &mData.Sections["RC"].(*mem.RecoveryCycle).LabelledTable
}

func (mef *Mef) rcNorm() RCNorm {
	norm := RCNorm{
		Interval: mem.RCInterval,
		GenericNorm: GenericNorm{
			mef:  mef,
			ltfm: rcTable,
		},
	}
	norm.MatNorm = MatrixNorm(norm)
	return norm
}
