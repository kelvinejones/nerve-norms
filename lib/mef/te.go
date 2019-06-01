package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type TENorm struct {
	Singles map[string]*teSingle `json:"data"`
}

func (mef *Mef) teNorm() TENorm {
	names := []string{"h40", "h20", "d40", "d20"}
	norm := TENorm{Singles: map[string]*teSingle{}}

	for _, name := range names {
		norm.Singles[name] = &teSingle{
			GenericNorm: GenericNorm{
				XValues: mem.IVCurrent,
				mef:     mef,
				ltfm:    teTableForSection(name),
			},
		}
		norm.Singles[name].CalculateNorms()
	}

	return norm
}

type teSingle struct {
	GenericNorm
}

func teTableForSection(sec string) LabelledTableFromMem {
	return func(mData *mem.Mem) *mem.LabelledTable {
		return (*mData.Sections["TE"].(*mem.ThresholdElectrotonus))[sec]
	}
}
