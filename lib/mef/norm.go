package mef

import (
	"math"

	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type Norm struct {
	CDNorm     NormTable    `json:"CD"`
	RCNorm     NormTable    `json:"RC"`
	SRNorm     SRNormTable  `json:"SR"`
	SRelNorm   NormTable    `json:"SRel"`
	IVNorm     NormTable    `json:"IV"`
	TENorm     TENormTables `json:"TE"`
	ExVarsNorm NormTable    `json:"ExVars"`
}

func (mef *Mef) Norm() Norm {
	norm := Norm{
		CDNorm:     NewNormTable(mem.CDDuration, mef, "CD", "", ArithmeticMean),
		IVNorm:     NewNormTable(mem.IVCurrent, mef, "IV", "", ArithmeticMean),
		RCNorm:     NewNormTable(mem.RCInterval, mef, "RC", "", ArithmeticMean),
		ExVarsNorm: NewNormTable(mem.ExVarIndices, mef, "ExVars", "", ArithmeticMean),
		SRNorm: SRNormTable{
			XNorm: NewNormTable(nil, mef, "SR", "calculatedX", GeometricMean),
			YNorm: NewNormTable(nil, mef, "SR", "calculatedY", GeometricMean),
		},
		SRelNorm: NewNormTable(mem.SRPercentMax, mef, "SR", "relative", ArithmeticMean),
		TENorm:   NewTENormTables(mef),
	}

	return norm
}

func (mef *Mef) Mean() *mem.Mem {
	norm := mef.Norm()
	memData := &mem.Mem{
		Header: mem.Header{
			File: "Calculated",
			Name: "Mean",
		},
		Sections: mem.Sections{
			"CD":     norm.CDNorm.asSection(),
			"RC":     norm.RCNorm.asSection(),
			"SR":     norm.SRNorm.asSection(),
			"TE":     norm.TENorm.asSection(),
			"IV":     norm.IVNorm.asSection(),
			"ExVars": norm.ExVarsNorm.asSection(),
		},
	}
	return memData
}

type OutScores struct {
	CDOutScores     OutScoresTable       `json:"CD"`
	RCOutScores     OutScoresTable       `json:"RC"`
	SROutScores     DoubleOutScoresTable `json:"SR"`
	SRelOutScores   OutScoresTable       `json:"SRel"`
	IVOutScores     OutScoresTable       `json:"IV"`
	TEOutScores     TEOutScores          `json:"TE"`
	ExVarsOutScores OutScoresTable       `json:"ExVars"`
	Overall         float64              `json:"Overall"`
}

func (norm *Norm) OutlierScores(mm *mem.Mem) OutScores {
	os := OutScores{
		CDOutScores:     NewOutScoresTable(norm.CDNorm, mm),
		IVOutScores:     NewOutScoresTable(norm.IVNorm, mm),
		RCOutScores:     NewOutScoresTable(norm.RCNorm, mm),
		ExVarsOutScores: NewOutScoresTable(norm.ExVarsNorm, mm),
		SROutScores:     NewDoubleOutScoresTable(norm.SRNorm.XNorm, norm.SRNorm.YNorm, mm),
		SRelOutScores:   NewOutScoresTable(norm.SRelNorm, mm),
		TEOutScores:     NewTEOutScores(norm.TENorm, mm),
	}

	os.Overall = nonZeroOverall([]float64{
		os.CDOutScores.Overall,
		os.RCOutScores.Overall,
		os.SROutScores.Overall,
		os.SRelOutScores.Overall,
		os.IVOutScores.Overall,
		os.TEOutScores.Overall,
		os.ExVarsOutScores.Overall,
	})

	return os
}

func nonZeroOverall(scores []float64) float64 {
	num := 0
	overall := 1.0
	for _, val := range scores {
		if val != 0 {
			num++
			overall *= val
		}
	}
	return math.Pow(overall, 1.0/float64(num))
}
