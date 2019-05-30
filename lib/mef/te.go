package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type TENorm struct {
	Interval mem.Column
	mef      *Mef
	singles  map[string]*teSingle
}

type teSingle struct {
	section string
	MatNorm
	mef *Mef
}

func (norm teSingle) NColumns() int {
	return len(norm.mef.mems)
}

func (norm teSingle) NRows() int {
	return len(mem.TEDelay)
}

func (norm teSingle) Column(i int) mem.Column {
	return norm.mef.mems[i].Sections["TE"].(*mem.ThresholdElectrotonus).GetPair(norm.section).ThreshReduction
}

func (norm teSingle) WasImputed(i int) mem.Column {
	return norm.mef.mems[i].Sections["TE"].(*mem.ThresholdElectrotonus).GetPair(norm.section).WasImputed
}

func (mef *Mef) teNorm() TENorm {
	norm := TENorm{
		Interval: mem.TEDelay,
		mef:      mef,
		singles: map[string]*teSingle{
			"h40": &teSingle{section: "h40", mef: mef},
			"h20": &teSingle{section: "h20", mef: mef},
			"d40": &teSingle{section: "d40", mef: mef},
			"d20": &teSingle{section: "d20", mef: mef},
		},
	}

	for i := range norm.singles {
		norm.singles[i].MatNorm = MatrixNorm(*norm.singles[i])
	}

	return norm
}
