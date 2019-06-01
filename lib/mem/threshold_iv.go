package mem

import (
	"errors"
)

type ThresholdIV struct{ LabelledTable }

var IVCurrent = Column([]float64{50, 40, 30, 20, 10, 0, -10, -20, -30, -40, -50, -60, -70, -80, -90, -100})

func (tiv *ThresholdIV) LoadFromMem(mem *rawMem) error {
	tiv.XName = "Current (%)"
	tiv.YName = "Threshold Reduction (%)"
	tiv.XColumn = IVCurrent

	sec, err := mem.sectionContainingHeader("THRESHOLD I/V")
	if err != nil {
		// Sometimes an old format spelled this incorrectly
		sec, err = mem.sectionContainingHeader("THESHOLD I/V")
		if err != nil {
			return errors.New("Could not get threshold IV: " + err.Error())
		}
	}

	curr, err := sec.columnContainsName("Current (%)", 0)
	if err != nil {
		return errors.New("Could not get threshold IV: " + err.Error())
	}

	tiv.YColumn, err = sec.columnContainsName("Threshold redn. (%)", 0)
	if err != nil {
		// Try alternative spelling
		tiv.YColumn, err = sec.columnContainsName("Threshold change (%)", 0)
		if err != nil {
			return errors.New("Could not get threshold IV: " + err.Error())
		}
	}

	tiv.WasImputed = tiv.YColumn.ImputeWithValue(curr, tiv.XColumn, 0.01, false)

	return nil
}
