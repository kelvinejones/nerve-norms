package mef

import "gogs.bellstone.ca/james/jitter/lib/mem"

type SRNorm struct {
	NormTable
	Cmap NormTable
}

func (mef *Mef) srNorm() SRNorm {
	return SRNorm{
		NormTable: NewNormTable(mem.SRPercentMax, mef, func(mData *mem.Mem) *mem.LabelledTable {
			return &mData.Sections["SR"].(*mem.StimResponse).LT
		}),
		Cmap: NewNormTable(nil, mef, func(mData *mem.Mem) *mem.LabelledTable {
			cmap, err := mData.Sections["SR"].(*mem.StimResponse).MaxCmaps.AsFloat()
			lt := mem.LabelledTable{
				XName:   "Time (ms)",
				YName:   "CMAP",
				XColumn: mem.Column{1},
				YColumn: mem.Column{cmap},
			}
			if err != nil {
				lt.WasImputed = mem.Column{1}
			}
			return &lt
		}),
	}
}
