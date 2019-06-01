package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

func rcTable(mData *mem.Mem) *mem.LabelledTable {
	return &mData.Sections["RC"].(*mem.RecoveryCycle).LabelledTable
}

func (mef *Mef) rcNorm() GenericNorm {
	return NewGenericNorm(mem.RCInterval, rcTable, mef)
}
