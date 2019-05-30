package mef

import (
	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type SRNorm struct {
	PercentMax mem.Column `json:"percentMax"`
	MatNorm    `json:"stimulus"`
	Cmap       maxCmapNorm
	mef        *Mef
}

func (norm SRNorm) NColumns() int {
	return len(norm.mef.mems)
}

func (norm SRNorm) NRows() int {
	return len(mem.SRPercentMax)
}

func (norm SRNorm) Column(i int) mem.Column {
	return norm.mef.mems[i].Sections["SR"].(*mem.StimResponse).Stimulus
}

func (norm SRNorm) WasImputed(i int) mem.Column {
	return norm.mef.mems[i].Sections["SR"].(*mem.StimResponse).WasImputed
}

func (mef *Mef) srNorm() SRNorm {
	norm := SRNorm{
		PercentMax: mem.SRPercentMax,
		mef:        mef,
		Cmap:       maxCmapNorm{mef: mef},
	}
	norm.MatNorm = MatrixNorm(norm)
	norm.Cmap.calcNorm()
	return norm
}

type maxCmapNorm struct {
	Mean float64
	SD   float64
	Num  float64
	mef  *Mef
}

func (norm maxCmapNorm) NColumns() int {
	return len(norm.mef.mems)
}

func (norm maxCmapNorm) NRows() int {
	return 1
}

func (norm maxCmapNorm) Column(i int) mem.Column {
	cmap, _ := norm.mef.mems[i].Sections["SR"].(*mem.StimResponse).MaxCmaps.AsFloat()
	return mem.Column{cmap}
}

func (norm maxCmapNorm) WasImputed(i int) mem.Column {
	_, err := norm.mef.mems[i].Sections["SR"].(*mem.StimResponse).MaxCmaps.AsFloat()
	wasImp := mem.Column(nil)
	if err != nil {
		wasImp = mem.Column{1.0}
	}
	return wasImp
}

func (mcn *maxCmapNorm) calcNorm() {
	matN := MatrixNorm(mcn)
	mcn.Mean = matN.Mean[0]
	mcn.SD = matN.SD[0]
	mcn.Num = matN.Num[0]
}
