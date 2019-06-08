package mef

import (
	"math"

	"gogs.bellstone.ca/james/jitter/lib/mem"
)

type Norm struct {
	CDNorm     *NormTable            `json:"CD"`
	RCNorm     *NormTable            `json:"RC"`
	SRNorm     *DoubleNormTable      `json:"SR"`
	SRelNorm   *NormTable            `json:"SRel"`
	IVNorm     *NormTable            `json:"IV"`
	TENorm     map[string]*NormTable `json:"TE"`
	ExVarsNorm *NormTable            `json:"ExVars"`
}

func (mef *Mef) Norm() Norm {
	norm := Norm{
		CDNorm:     NewNormTable(mem.CDDuration, mef, "CD", "", ArithmeticMean),
		IVNorm:     NewNormTable(mem.IVCurrent, mef, "IV", "", ArithmeticMean),
		RCNorm:     NewNormTable(mem.RCInterval, mef, "RC", "", ArithmeticMean),
		ExVarsNorm: NewNormTable(mem.ExVarIndices, mef, "ExVars", "", ArithmeticMean),
		SRNorm: &DoubleNormTable{
			XNorm: NewNormTable(nil, mef, "SR", "calculatedX", GeometricMean),
			YNorm: NewNormTable(nil, mef, "SR", "calculatedY", GeometricMean),
		},
		SRelNorm: NewNormTable(mem.SRPercentMax, mef, "SR", "relative", ArithmeticMean),
		TENorm:   map[string]*NormTable{},
	}

	for _, name := range []string{"h40", "h20", "d40", "d20"} {
		nt := NewNormTable(mem.TEDelay(name), mef, "TE", name, ArithmeticMean)
		if nt != nil {
			// We only add this TE type of it's not zero
			norm.TENorm[name] = nt
		}
	}

	return norm
}

type OutScores struct {
	CDOutScores     *OutScoresTable       `json:"CD"`
	RCOutScores     *OutScoresTable       `json:"RC"`
	SROutScores     *DoubleOutScoresTable `json:"SR"`
	SRelOutScores   *OutScoresTable       `json:"SRel"`
	IVOutScores     *OutScoresTable       `json:"IV"`
	TEOutScores     *TEOutScores          `json:"TE"`
	ExVarsOutScores *OutScoresTable       `json:"ExVars"`
	Overall         float64               `json:"Overall"`
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
