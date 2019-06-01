package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

func cdTable(mData *mem.Mem) *mem.LabelledTable {
	return &mData.Sections["CD"].(*mem.ChargeDuration).LabelledTable
}

func (mef *Mef) cdNorm() GenericNorm {
	return NewGenericNorm(mem.CDDuration, cdTable, mef)
}
