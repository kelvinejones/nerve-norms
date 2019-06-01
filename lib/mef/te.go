package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type TENorm struct {
	Singles map[string]*teSingle `json:"data"`
	mef     *Mef
}

type teSingle struct {
	Delay   mem.Column `json:"delay"`
	MatNorm `json:"threshReduction"`
	section string
	mef     *Mef
}

func (norm teSingle) NColumns() int {
	return len(norm.mef.mems)
}

func (norm teSingle) NRows() int {
	return len(mem.TEDelay(norm.section))
}

func (norm teSingle) Column(i int) mem.Column {
	return (*norm.mef.mems[i].Sections["TE"].(*mem.ThresholdElectrotonus))[norm.section].YColumn
}

func (norm teSingle) WasImputed(i int) mem.Column {
	return (*norm.mef.mems[i].Sections["TE"].(*mem.ThresholdElectrotonus))[norm.section].WasImputed
}

func (mef *Mef) teNorm() TENorm {
	names := []string{"h40", "h20", "d40", "d20"}
	norm := TENorm{
		mef:     mef,
		Singles: map[string]*teSingle{},
	}

	for _, name := range names {
		norm.Singles[name] = &teSingle{
			Delay:   mem.TEDelay(name),
			section: name,
			mef:     mef,
		}
		norm.Singles[name].MatNorm = MatrixNorm(*norm.Singles[name])
	}

	return norm
}
