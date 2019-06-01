package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

func ivTable(mData *mem.Mem) *mem.LabelledTable {
	return &mData.Sections["IV"].(*mem.ThresholdIV).LabelledTable
}

func (mef *Mef) ivNorm() GenericNorm {
	return NewGenericNorm(mem.IVCurrent, ivTable, mef)
}
