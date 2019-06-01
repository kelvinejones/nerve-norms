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
		CDNorm: NewNormTable(mem.CDDuration, mef, func(mData *mem.Mem) mem.LabelledTable {
			return &mData.Sections["CD"].(*mem.ChargeDuration).LabTab
		}),
		IVNorm: NewNormTable(mem.IVCurrent, mef, func(mData *mem.Mem) mem.LabelledTable {
			return &mData.Sections["IV"].(*mem.ThresholdIV).LabTab
		}),
		RCNorm: NewNormTable(mem.RCInterval, mef, func(mData *mem.Mem) mem.LabelledTable {
			return &mData.Sections["RC"].(*mem.RecoveryCycle).LabTab
		}),
		SRNorm: mef.srNorm(),
		TENorm: map[string]NormTable{},
	}

	for _, name := range []string{"h40", "h20", "d40", "d20"} {
		norm.TENorm[name] = NewNormTable(mem.TEDelay(name), mef, func(mData *mem.Mem) mem.LabelledTable {
			return (*mData.Sections["TE"].(*mem.ThresholdElectrotonus))[name]
		})
	}

	return norm
}
