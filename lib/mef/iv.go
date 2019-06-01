package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type IVNorm struct {
	Current mem.Column `json:"current"`
	GenericNorm
}

func ivTable(mData *mem.Mem) *mem.LabelledTable {
	return &mData.Sections["IV"].(*mem.ThresholdIV).LabelledTable
}

func (mef *Mef) ivNorm() IVNorm {
	norm := IVNorm{
		Current: mem.IVCurrent,
		GenericNorm: GenericNorm{
			mef:  mef,
			ltfm: ivTable,
		},
	}
	norm.MatNorm = MatrixNorm(norm)
	return norm
}
