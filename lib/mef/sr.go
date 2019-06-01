package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type SRNorm struct {
	GenericNorm
	Cmap GenericNorm
}

func srTable(mData *mem.Mem) *mem.LabelledTable {
	return &mData.Sections["SR"].(*mem.StimResponse).LT
}

func (mef *Mef) srNorm() SRNorm {
	norm := SRNorm{
		GenericNorm: GenericNorm{XValues: mem.SRPercentMax},
	}
	norm.CalculateNorms(srTable, mef)
	norm.Cmap.CalculateNorms(maxCmapTable, mef)
	return norm
}

func maxCmapTable(mData *mem.Mem) *mem.LabelledTable {
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
}
