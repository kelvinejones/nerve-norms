package mef

import "gogs.bellstone.ca/james/jitter/lib/mem"

type SRNorm struct {
	NormTable
	Cmap NormTable
}

func (mef *Mef) srNorm() SRNorm {
	return SRNorm{
		NormTable: NewNormTable(mem.SRPercentMax, mef, "SR", ""),
		Cmap:      NewNormTable(nil, mef, "SR", "CMAP"),
	}
}
