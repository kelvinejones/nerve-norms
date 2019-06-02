package mef

import "gogs.bellstone.ca/james/jitter/lib/mem"

type Norm struct {
	CDNorm   NormTable            `json:"CD"`
	RCNorm   NormTable            `json:"RC"`
	SRNorm   DoubleNormTable      `json:"SR"`
	SRelNorm NormTable            `json:"SRel"`
	IVNorm   NormTable            `json:"IV"`
	TENorm   map[string]NormTable `json:"TE"`
}

func (mef *Mef) Norm() Norm {
	mef = mef.FilteredMef()

	norm := Norm{
		CDNorm: NewNormTable(mem.CDDuration, mef, "CD", ""),
		IVNorm: NewNormTable(mem.IVCurrent, mef, "IV", ""),
		RCNorm: NewNormTable(mem.RCInterval, mef, "RC", ""),
		SRNorm: DoubleNormTable{
			XNorm: NewNormTable(nil, mef, "SR", "calculatedX"),
			YNorm: NewNormTable(nil, mef, "SR", "calculatedY"),
		},
		SRelNorm: NewNormTable(mem.SRPercentMax, mef, "SR", "relative"),
		TENorm:   map[string]NormTable{},
	}

	for _, name := range []string{"h40", "h20", "d40", "d20"} {
		norm.TENorm[name] = NewNormTable(mem.TEDelay(name), mef, "TE", name)
	}

	return norm
}
