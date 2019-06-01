package mef

import "gogs.bellstone.ca/james/jitter/lib/mem"

type Norm struct {
	CDNorm GenericNorm `json:"cd"`
	RCNorm GenericNorm `json:"rc"`
	SRNorm `json:"sr"`
	IVNorm GenericNorm `json:"iv"`
	TENorm `json:"te"`
}

func (mef *Mef) Norm() Norm {
	mef = mef.FilteredMef()

	return Norm{
		CDNorm: mef.cdNorm(),
		RCNorm: mef.rcNorm(),
		SRNorm: mef.srNorm(),
		IVNorm: mef.ivNorm(),
		TENorm: mef.teNorm(),
	}
}

func (mef *Mef) cdNorm() GenericNorm {
	return NewGenericNorm(mem.CDDuration, mef, func(mData *mem.Mem) *mem.LabelledTable {
		return &mData.Sections["CD"].(*mem.ChargeDuration).LabelledTable
	})
}

func (mef *Mef) ivNorm() GenericNorm {
	return NewGenericNorm(mem.IVCurrent, mef, func(mData *mem.Mem) *mem.LabelledTable {
		return &mData.Sections["IV"].(*mem.ThresholdIV).LabelledTable
	})
}

func (mef *Mef) rcNorm() GenericNorm {
	return NewGenericNorm(mem.RCInterval, mef, func(mData *mem.Mem) *mem.LabelledTable {
		return &mData.Sections["RC"].(*mem.RecoveryCycle).LabelledTable
	})
}

type SRNorm struct {
	GenericNorm
	Cmap GenericNorm
}

func (mef *Mef) srNorm() SRNorm {
	return SRNorm{
		GenericNorm: NewGenericNorm(mem.SRPercentMax, mef, func(mData *mem.Mem) *mem.LabelledTable {
			return &mData.Sections["SR"].(*mem.StimResponse).LT
		}),
		Cmap: NewGenericNorm(nil, mef, func(mData *mem.Mem) *mem.LabelledTable {
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

type TENorm map[string]GenericNorm

func (mef *Mef) teNorm() TENorm {
	norm := TENorm(map[string]GenericNorm{})
	for _, name := range []string{"h40", "h20", "d40", "d20"} {
		norm[name] = NewGenericNorm(mem.TEDelay(name), mef, func(mData *mem.Mem) *mem.LabelledTable {
			return (*mData.Sections["TE"].(*mem.ThresholdElectrotonus))[name]
		})
	}
	return norm
}
