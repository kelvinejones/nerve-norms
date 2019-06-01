package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type TENorm map[string]GenericNorm

func (mef *Mef) teNorm() TENorm {
	names := []string{"h40", "h20", "d40", "d20"}
	norm := TENorm(map[string]GenericNorm{})

	for _, name := range names {
		norm[name] = NewGenericNorm(mem.TEDelay(name), teTableForSection(name), mef)
	}

	return norm
}

func teTableForSection(sec string) LabelledTableFromMem {
	return func(mData *mem.Mem) *mem.LabelledTable {
		return (*mData.Sections["TE"].(*mem.ThresholdElectrotonus))[sec]
	}
}
