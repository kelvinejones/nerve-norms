package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type TENorm struct {
	mef     *Mef
	singles map[string]*teSingle
}

type teSingle struct {
	Delay   mem.Column
	section string
	MatNorm
	mef *Mef
}

func (norm teSingle) NColumns() int {
	return len(norm.mef.mems)
}

func (norm teSingle) NRows() int {
	return len(mem.TEDelay(norm.section))
}

func (norm teSingle) Column(i int) mem.Column {
	return norm.mef.mems[i].Sections["TE"].(*mem.ThresholdElectrotonus).Data[norm.section].ThreshReduction
}

func (norm teSingle) WasImputed(i int) mem.Column {
	return norm.mef.mems[i].Sections["TE"].(*mem.ThresholdElectrotonus).Data[norm.section].WasImputed
}

func (mef *Mef) teNorm() TENorm {
	norm := TENorm{
		mef: mef,
		singles: map[string]*teSingle{
			"h40": &teSingle{
				Delay:   mem.TEDelay("h40"),
				section: "h40",
				mef:     mef,
			},
			"h20": &teSingle{
				Delay:   mem.TEDelay("h20"),
				section: "h20",
				mef:     mef,
			},
			"d40": &teSingle{
				Delay:   mem.TEDelay("d40"),
				section: "d40",
				mef:     mef,
			},
			"d20": &teSingle{
				Delay:   mem.TEDelay("d20"),
				section: "d20",
				mef:     mef,
			},
		},
	}

	for i := range norm.singles {
		norm.singles[i].MatNorm = MatrixNorm(*norm.singles[i])
	}

	return norm
}
