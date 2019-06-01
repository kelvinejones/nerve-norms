package mef

import "gogs.bellstone.ca/james/jitter/lib/mem"

type SRNorm struct {
	NormTable
	Cmap NormTable
}

func (mef *Mef) srNorm() SRNorm {
	return SRNorm{
		NormTable: NewNormTable(mem.SRPercentMax, mef, func(mData *mem.Mem) mem.LabelledTable {
			return &mData.Sections["SR"].(*mem.StimResponse).LT
		}),
		Cmap: NewNormTable(nil, mef, func(mData *mem.Mem) mem.LabelledTable {
			return mData.Sections["SR"].(*mem.StimResponse).MaxCmaps.AsLabelledTable()
		}),
	}
}
