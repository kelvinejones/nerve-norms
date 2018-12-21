package mem

import (
	"errors"
)

type RecoveryCycle struct {
	Interval     Column
	ThreshChange Column
}

func (mem *Mem) RecoveryCycle() (RecoveryCycle, error) {
	tiv := RecoveryCycle{}

	sec, err := mem.sectionContainingHeader("RECOVERY CYCLE")
	if err != nil {
		return tiv, errors.New("Could not get recovery cycle: " + err.Error())
	}

	tiv.Interval, err = sec.columnContainsName("Interval (ms)", 0)
	if err != nil {
		return tiv, errors.New("Could not get recovery cycle: " + err.Error())
	}

	tiv.ThreshChange, err = sec.columnContainsName("Threshold change (%)", 0)
	if err != nil {
		return tiv, errors.New("Could not get recovery cycle: " + err.Error())
	}

	return tiv, nil
}
