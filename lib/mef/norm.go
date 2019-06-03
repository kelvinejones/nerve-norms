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
	norm := Norm{
		CDNorm: NewNormTable(mem.CDDuration, mef, "CD", "", ArithmeticMean),
		IVNorm: NewNormTable(mem.IVCurrent, mef, "IV", "", ArithmeticMean),
		RCNorm: NewNormTable(mem.RCInterval, mef, "RC", "", ArithmeticMean),
		SRNorm: DoubleNormTable{
			XNorm: NewNormTable(nil, mef, "SR", "calculatedX", GeometricMean),
			YNorm: NewNormTable(nil, mef, "SR", "calculatedY", GeometricMean),
		},
		SRelNorm: NewNormTable(mem.SRPercentMax, mef, "SR", "relative", ArithmeticMean),
		TENorm:   map[string]NormTable{},
	}

	for _, name := range []string{"h40", "h20", "d40", "d20"} {
		norm.TENorm[name] = NewNormTable(mem.TEDelay(name), mef, "TE", name, ArithmeticMean)
	}

	return norm
}
