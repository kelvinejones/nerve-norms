package mem

import (
	"errors"
)

type ThresholdIV struct {
	Current         Column
	ThreshReduction Column
}

func (mem *Mem) ThresholdIV() (ThresholdIV, error) {
	tiv := ThresholdIV{Current: Column([]float64{-100, -90, -80, -70, -60, -50, -40, -30, -20, -10, 0, 10, 20, 30, 40, 50})}

	sec, err := mem.sectionContainingHeader("THRESHOLD I/V")
	if err != nil {
		return tiv, errors.New("Could not get threshold IV: " + err.Error())
	}

	curr, err := sec.columnContainsName("Current (%)", 0)
	if err != nil {
		return tiv, errors.New("Could not get threshold IV: " + err.Error())
	}
	if !tiv.Current.Equals(curr, 0.01) {
		return tiv, errors.New("File contains invalid tiv Current")
	}

	tiv.ThreshReduction, err = sec.columnContainsName("Threshold redn. (%)", 0)
	if err != nil {
		return tiv, errors.New("Could not get threshold IV: " + err.Error())
	}

	return tiv, nil
}
