package mem

import (
	"errors"
)

type RecoveryCycle struct {
	Interval     Column
	ThreshChange Column
}

func (mem *Mem) RecoveryCycle() (RecoveryCycle, error) {
	rc := RecoveryCycle{Interval: Column([]float64{2, 2.5, 3.2, 4, 5, 6.3, 7.9, 10, 13, 18, 24, 32, 42, 56, 75, 100, 140, 200})}

	sec, err := mem.sectionContainingHeader("RECOVERY CYCLE")
	if err != nil {
		return rc, errors.New("Could not get recovery cycle: " + err.Error())
	}

	interval, err := sec.columnContainsName("Interval (ms)", 0)
	if err != nil {
		return rc, errors.New("Could not get recovery cycle: " + err.Error())
	}
	if !rc.Interval.Equals(interval, 0.000001) {
		return rc, errors.New("File contains invalid rc Interval")
	}

	rc.ThreshChange, err = sec.columnContainsName("Threshold change (%)", 0)
	if err != nil {
		return rc, errors.New("Could not get recovery cycle: " + err.Error())
	}

	return rc, nil
}
