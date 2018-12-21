package mem

import (
	"errors"
)

type ChargeDuration struct {
	Duration     Column
	ThreshCharge Column
}

func (mem *Mem) ChargeDuration() (ChargeDuration, error) {
	tiv := ChargeDuration{}

	sec, err := mem.sectionContainingHeader("CHARGE DURATION")
	if err != nil {
		return tiv, errors.New("Could not get charge duration: " + err.Error())
	}

	tiv.Duration, err = sec.columnContainsName("Duration (ms)", 0)
	if err != nil {
		return tiv, errors.New("Could not get charge duration: " + err.Error())
	}

	tiv.ThreshCharge, err = sec.columnContainsName("Threshold charge (mA.mS)", 0)
	if err != nil {
		return tiv, errors.New("Could not get charge duration: " + err.Error())
	}

	return tiv, nil
}
