package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type IVNorm struct{ GenericNorm }

func ivTable(mData *mem.Mem) *mem.LabelledTable {
	return &mData.Sections["IV"].(*mem.ThresholdIV).LabelledTable
}

func (mef *Mef) ivNorm() IVNorm {
	norm := IVNorm{
		GenericNorm: GenericNorm{
			XValues: mem.IVCurrent,
			mef:     mef,
			ltfm:    ivTable,
		},
	}
	norm.CalculateNorms()
	return norm
}
