package mef

import "gogs.bellstone.ca/james/jitter/lib/mem"

type Norm struct {
	CDNorm GenericNorm `json:"cd"`
	RCNorm GenericNorm `json:"rc"`
	SRNorm `json:"sr"`
	IVNorm GenericNorm            `json:"iv"`
	TENorm map[string]GenericNorm `json:"te"`
}

type SRNorm struct {
	GenericNorm
	Cmap GenericNorm
}

func (mef *Mef) Norm() Norm {
	mef = mef.FilteredMef()

	norm := Norm{
		CDNorm: NewGenericNorm(mem.CDDuration, mef, func(mData *mem.Mem) *mem.LabelledTable {
			return &mData.Sections["CD"].(*mem.ChargeDuration).LabelledTable
		}),
		IVNorm: NewGenericNorm(mem.IVCurrent, mef, func(mData *mem.Mem) *mem.LabelledTable {
			return &mData.Sections["IV"].(*mem.ThresholdIV).LabelledTable
		}),
		RCNorm: NewGenericNorm(mem.RCInterval, mef, func(mData *mem.Mem) *mem.LabelledTable {
			return &mData.Sections["RC"].(*mem.RecoveryCycle).LabelledTable
		}),
		TENorm: map[string]GenericNorm{},
		SRNorm: SRNorm{
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
		},
	}

	for _, name := range []string{"h40", "h20", "d40", "d20"} {
		norm.TENorm[name] = NewGenericNorm(mem.TEDelay(name), mef, func(mData *mem.Mem) *mem.LabelledTable {
			return (*mData.Sections["TE"].(*mem.ThresholdElectrotonus))[name]
		})
	}

	return norm
}
