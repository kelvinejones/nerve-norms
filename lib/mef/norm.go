package mef

import "gogs.bellstone.ca/james/jitter/lib/mem"

type Norm struct {
	CDNorm NormTable `json:"cd"`
	RCNorm NormTable `json:"rc"`
	SRNorm `json:"sr"`
	IVNorm NormTable            `json:"iv"`
	TENorm map[string]NormTable `json:"te"`
}

func (mef *Mef) Norm() Norm {
	mef = mef.FilteredMef()

	norm := Norm{
		CDNorm: NewNormTable(mem.CDDuration, mef, mem.CDLabelledTable),
		IVNorm: NewNormTable(mem.IVCurrent, mef, mem.IVLabelledTable),
		RCNorm: NewNormTable(mem.RCInterval, mef, mem.RCLabelledTable),
		SRNorm: mef.srNorm(),
		TENorm: map[string]NormTable{},
	}

	for _, name := range []string{"h40", "h20", "d40", "d20"} {
		norm.TENorm[name] = NewNormTable(mem.TEDelay(name), mef, mem.TELabelledTable(name))
	}

	return norm
}
