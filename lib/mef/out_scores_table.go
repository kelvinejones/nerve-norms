package mef

import (
	"encoding/json"
	"errors"

	"gogs.bellstone.ca/james/jitter/lib/mem"

	"github.com/aclements/go-moremath/stats"
)

type OutScoresTable struct {
	Values mem.Column // Usually these are the x-values for a y-mean.
	Scores mem.Column
}

func NewOutScoresTable(norm NormTable, mm *mem.Mem) OutScoresTable {
	lt := mm.LabelledTable(norm.sec, norm.subsec)
	numEl := lt.Len()
	ost := OutScoresTable{
		Values: norm.Values,
		Scores: make(mem.Column, numEl),
	}

	dist := stats.NormalDist{Mu: 0.0, Sigma: 1.0}

	for rowN := 0; rowN < numEl; rowN++ {
		// TODO account for GeometricMean
		if norm.SD[rowN] == 0.0 {
			ost.Scores[rowN] = 0.0
		} else {
			diff := (norm.Mean[rowN] - lt.YColumnAt(rowN)) / norm.SD[rowN]
			if diff > 0 {
				diff *= -1
			}
			ost.Scores[rowN] = 1 - 2*dist.CDF(diff)
		}
	}

	return ost
}

func (ost OutScoresTable) MarshalJSON() ([]byte, error) {
	jt := jsonTable{
		Columns: []string{"scores"},
		Data:    []mem.Column{ost.Scores},
	}

	if ost.Values != nil {
		jt.Columns = append(jt.Columns, "values")
		jt.Data = append(jt.Data, ost.Values)
	}

	return json.Marshal(&jt)
}

func (ost *OutScoresTable) UnmarshalJSON(value []byte) error {
	var jt jsonTable
	err := json.Unmarshal(value, &jt)
	if err != nil {
		return err
	}

	numCol := len(jt.Columns)
	numDat := len(jt.Data)

	if numCol < 1 || numCol > 2 || numDat < 1 || numDat > 2 {
		return errors.New("Incorrect number of OutScoresTable columns in JSON")
	}

	ost.Scores = jt.Data[0]

	if numCol == 2 {
		ost.Values = jt.Data[1]
	}

	return nil
}

type DoubleOutScoresTable struct {
	YOutScores OutScoresTable
	XOutScores OutScoresTable
}

func (ost DoubleOutScoresTable) MarshalJSON() ([]byte, error) {
	jt := jsonTable{
		Columns: []string{"yscores", "yvalues", "xscores", "xvalues"},
		Data:    []mem.Column{ost.YOutScores.Scores, ost.YOutScores.Values, ost.XOutScores.Scores, ost.XOutScores.Values},
	}

	return json.Marshal(&jt)
}

func (ost *DoubleOutScoresTable) UnmarshalJSON(value []byte) error {
	var jt jsonTable
	err := json.Unmarshal(value, &jt)
	if err != nil {
		return err
	}

	if len(jt.Columns) != 4 || len(jt.Data) != 4 {
		return errors.New("Incorrect number of DoubleOutScoresTable columns in JSON")
	}

	ost.YOutScores.Scores = jt.Data[0]
	ost.YOutScores.Values = jt.Data[1]
	ost.XOutScores.Scores = jt.Data[3]
	ost.XOutScores.Values = jt.Data[4]

	return nil
}
