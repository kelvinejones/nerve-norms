package mef

import (
	"encoding/json"
	"errors"
	"math"

	"gogs.bellstone.ca/james/jitter/lib/mem"

	"github.com/aclements/go-moremath/stats"
)

type OutScoresTable struct {
	Values  mem.Column // Usually these are the x-values for a y-mean.
	Scores  mem.Column
	Overall float64
}

var dist = stats.NormalDist{Mu: 0.0, Sigma: 1.0}

func NewOutScoresTable(norm NormTable, mm *mem.Mem) OutScoresTable {
	lt := mm.LabelledTable(norm.sec, norm.subsec)
	numEl := lt.Len()
	ost := OutScoresTable{
		Values:  norm.Values,
		Scores:  make(mem.Column, numEl),
		Overall: 1,
	}

	numScored := 0
	for rowN := 0; rowN < numEl; rowN++ {
		diff := norm.numSD(rowN, lt.YColumnAt(rowN))
		if diff > 0 {
			diff *= -1
		}
		ost.Scores[rowN] = 2 * dist.CDF(diff)
		if lt.IncludeOutlierScore(rowN) {
			ost.Overall *= ost.Scores[rowN]
			numScored++
		}
	}
	ost.Overall = math.Pow(ost.Overall, 1.0/float64(numScored))

	return ost
}

func (norm NormTable) numSD(rowN int, val float64) float64 {
	if norm.SD[rowN] == 0.0 {
		return 0.0
	}
	switch norm.MeanType {
	case ArithmeticMean:
		return (norm.Mean[rowN] - val) / norm.SD[rowN]
	case GeometricMean:
		return (math.Log10(norm.Mean[rowN]) - math.Log10(val)) / norm.SD[rowN]
	default:
		return 0.0
	}
}

// osJsonTable is used to restructure LabTab data for json.
type osJsonTable struct {
	Columns []string  `json:"columns"`
	Data    mem.Table `json:"data"`
	Overall float64   `json:"Overall"`
}

func (ost OutScoresTable) MarshalJSON() ([]byte, error) {
	jt := osJsonTable{
		Columns: []string{"scores"},
		Data:    []mem.Column{ost.Scores},
		Overall: ost.Overall,
	}

	if ost.Values != nil {
		jt.Columns = append(jt.Columns, "values")
		jt.Data = append(jt.Data, ost.Values)
	}

	return json.Marshal(&jt)
}

func (ost *OutScoresTable) UnmarshalJSON(value []byte) error {
	var jt osJsonTable
	err := json.Unmarshal(value, &jt)
	if err != nil {
		return err
	}

	numCol := len(jt.Columns)
	numDat := len(jt.Data)

	if numCol < 1 || numCol > 2 || numDat < 1 || numDat > 2 {
		return errors.New("Incorrect number of OutScoresTable columns in JSON")
	}

	ost.Overall = jt.Overall
	ost.Scores = jt.Data[0]

	if numCol == 2 {
		ost.Values = jt.Data[1]
	}

	return nil
}

type DoubleOutScoresTable struct {
	YOutScores OutScoresTable
	XOutScores OutScoresTable
	Overall    float64
}

func NewDoubleOutScoresTable(norm1, norm2 NormTable, mm *mem.Mem) DoubleOutScoresTable {
	dost := DoubleOutScoresTable{
		XOutScores: NewOutScoresTable(norm1, mm),
		YOutScores: NewOutScoresTable(norm2, mm),
	}
	dost.Overall = math.Sqrt(dost.XOutScores.Overall * dost.YOutScores.Overall)
	return dost
}

func (ost DoubleOutScoresTable) MarshalJSON() ([]byte, error) {
	jt := osJsonTable{
		Columns: []string{"yscores", "yvalues", "xscores", "xvalues"},
		Data:    []mem.Column{ost.YOutScores.Scores, ost.YOutScores.Values, ost.XOutScores.Scores, ost.XOutScores.Values},
		Overall: ost.Overall,
	}

	return json.Marshal(&jt)
}

func (ost *DoubleOutScoresTable) UnmarshalJSON(value []byte) error {
	var jt osJsonTable
	err := json.Unmarshal(value, &jt)
	if err != nil {
		return err
	}

	if len(jt.Columns) != 4 || len(jt.Data) != 4 {
		return errors.New("Incorrect number of DoubleOutScoresTable columns in JSON")
	}

	ost.Overall = jt.Overall
	ost.YOutScores.Scores = jt.Data[0]
	ost.YOutScores.Values = jt.Data[1]
	ost.XOutScores.Scores = jt.Data[3]
	ost.XOutScores.Values = jt.Data[4]

	return nil
}

type TEOutScores struct {
	OutScores map[string]OutScoresTable
	Overall   float64
}

func NewTEOutScores(norm map[string]NormTable, mm *mem.Mem) TEOutScores {
	teos := TEOutScores{
		OutScores: map[string]OutScoresTable{},
		Overall:   1,
	}
	for _, name := range []string{"h40", "h20", "d40", "d20"} {
		teos.OutScores[name] = NewOutScoresTable(norm[name], mm)
		teos.Overall *= teos.OutScores[name].Overall
	}
	teos.Overall = math.Pow(teos.Overall, 0.25)
	return teos
}
