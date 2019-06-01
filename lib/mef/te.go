package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type TENorm map[string]GenericNorm

func (mef *Mef) teNorm() TENorm {
	norm := TENorm(map[string]GenericNorm{})
	for _, name := range []string{"h40", "h20", "d40", "d20"} {
		teTableFunc := func(mData *mem.Mem) *mem.LabelledTable {
			return (*mData.Sections["TE"].(*mem.ThresholdElectrotonus))[name]
		}
		norm[name] = NewGenericNorm(mem.TEDelay(name), teTableFunc, mef)
	}
	return norm
}
