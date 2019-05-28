package mem

import (
	"errors"
	"fmt"
)

type ThresholdIV struct {
	Current         Column
	ThreshReduction Column
	WasImputed      Column
}

func (mem *Mem) ThresholdIV() (ThresholdIV, error) {
	tiv := ThresholdIV{Current: Column([]float64{50, 40, 30, 20, 10, 0, -10, -20, -30, -40, -50, -60, -70, -80, -90, -100})}

	sec, err := mem.sectionContainingHeader("THRESHOLD I/V")
	if err != nil {
		return tiv, errors.New("Could not get threshold IV: " + err.Error())
	}

	curr, err := sec.columnContainsName("Current (%)", 0)
	if err != nil {
		return tiv, errors.New("Could not get threshold IV: " + err.Error())
	}

	tiv.ThreshReduction, err = sec.columnContainsName("Threshold redn. (%)", 0)
	if err != nil {
		return tiv, errors.New("Could not get threshold IV: " + err.Error())
	}

	old := tiv.ThreshReduction
	tiv.WasImputed = tiv.ThreshReduction.ImputeWithValue(curr, tiv.Current, 0.01)
	if tiv.WasImputed != nil {
		fmt.Println("Imputed TIV:", old, tiv.ThreshReduction)
	}

	return tiv, nil
}
