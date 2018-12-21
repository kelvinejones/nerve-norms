package mem

import (
	"errors"
)

type ThresholdIV struct {
	Current         Column
	ThreshReduction Column
}

func (mem *Mem) ThresholdIV() (ThresholdIV, error) {
	tiv := ThresholdIV{}

	sec, err := mem.sectionContainingHeader("THRESHOLD I/V")
	if err != nil {
		return tiv, errors.New("Could not get threshold IV: " + err.Error())
	}

	tiv.Current, err = sec.columnContainsName("Current (%)", 0)
	if err != nil {
		return tiv, errors.New("Could not get threshold IV: " + err.Error())
	}

	tiv.ThreshReduction, err = sec.columnContainsName("Threshold redn. (%)", 0)
	if err != nil {
		return tiv, errors.New("Could not get threshold IV: " + err.Error())
	}

	return tiv, nil
}
