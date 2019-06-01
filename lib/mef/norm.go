package mef

import "gogs.bellstone.ca/james/jitter/lib/mem"

type Norm struct {
	CDNorm NormTable `json:"cd"`
	RCNorm NormTable `json:"rc"`
	SRNorm `json:"sr"`
	IVNorm NormTable            `json:"iv"`
	TENorm map[string]NormTable `json:"te"`
}

type SRNorm struct {
	NormTable
	Cmap NormTable
}

func (mef *Mef) Norm() Norm {
	mef = mef.FilteredMef()

	norm := Norm{
		CDNorm: NewNormTable(mem.CDDuration, mef, func(mData *mem.Mem) *mem.LabelledTable {
			return &mData.Sections["CD"].(*mem.ChargeDuration).LabelledTable
		}),
		IVNorm: NewNormTable(mem.IVCurrent, mef, func(mData *mem.Mem) *mem.LabelledTable {
			return &mData.Sections["IV"].(*mem.ThresholdIV).LabelledTable
		}),
		RCNorm: NewNormTable(mem.RCInterval, mef, func(mData *mem.Mem) *mem.LabelledTable {
			return &mData.Sections["RC"].(*mem.RecoveryCycle).LabelledTable
		}),
		TENorm: map[string]NormTable{},
		SRNorm: SRNorm{
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
		},
	}

	for _, name := range []string{"h40", "h20", "d40", "d20"} {
		norm.TENorm[name] = NewNormTable(mem.TEDelay(name), mef, func(mData *mem.Mem) *mem.LabelledTable {
			return (*mData.Sections["TE"].(*mem.ThresholdElectrotonus))[name]
		})
	}

	return norm
}
