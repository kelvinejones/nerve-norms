package mem

import (
	"errors"
)

type RecoveryCycle struct{ LabTab }

var RCInterval = Column([]float64{2, 2.5, 3.2, 4, 5, 6.3, 7.9, 10, 13, 18, 24, 32, 42, 56, 75, 100, 140, 200})

func RCLabelledTable(mem *Mem) LabelledTable {
	return &mem.Sections["RC"].(*RecoveryCycle).LabTab
}

func (rc *RecoveryCycle) LoadFromMem(mem *rawMem) error {
	rc.xname = "Interval (ms)"
	rc.yname = "Threshold change (%)"
	rc.xcol = RCInterval

	sec, err := mem.sectionContainingHeader("RECOVERY CYCLE")
	if err != nil {
		return errors.New("Could not get recovery cycle: " + err.Error())
	}

	interval, err := sec.columnContainsName("Interval (ms)", 0)
	if err != nil {
		return errors.New("Could not get recovery cycle: " + err.Error())
	}

	rc.ycol, err = sec.columnContainsName("Threshold change (%)", 0)
	if err != nil {
		return errors.New("Could not get recovery cycle: " + err.Error())
	}

	rc.wasimp = rc.ycol.ImputeWithValue(interval, rc.xcol, 0.000001, true)

	return nil
}
