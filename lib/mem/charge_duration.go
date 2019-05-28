package mem

import (
	"errors"
)

type ChargeDuration struct {
	Duration     Column
	ThreshCharge Column
}

func (mem *Mem) ChargeDuration() (ChargeDuration, error) {
	cd := ChargeDuration{Duration: Column([]float64{0.2, 0.4, 0.6, 0.8, 1})}

	sec, err := mem.sectionContainingHeader("CHARGE DURATION")
	if err != nil {
		return cd, errors.New("Could not get charge duration: " + err.Error())
	}

	dur, err := sec.columnContainsName("Duration (ms)", 0)
	if err != nil {
		return cd, errors.New("Could not get charge duration: " + err.Error())
	}
	if !cd.Duration.Equals(dur, 0.0000001) {
		return cd, errors.New("File contains invalid cd Duration")
	}

	cd.ThreshCharge, err = sec.columnContainsName("Threshold charge (mA.mS)", 0)
	if err != nil {
		return cd, errors.New("Could not get charge duration: " + err.Error())
	}

	return cd, nil
}
