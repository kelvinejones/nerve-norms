package mem

import (
	"errors"
)

type ThresholdIV struct{ LabTab }

var IVCurrent = Column([]float64{50, 40, 30, 20, 10, 0, -10, -20, -30, -40, -50, -60, -70, -80, -90, -100})

func IVLabelledTable(mem *Mem) LabelledTable {
	return &mem.Sections["IV"].(*ThresholdIV).LabTab
}

func (tiv *ThresholdIV) LoadFromMem(mem *rawMem) error {
	tiv.xname = "Current (%)"
	tiv.yname = "Threshold Reduction (%)"
	tiv.xcol = IVCurrent

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

	tiv.ycol, err = sec.columnContainsName("Threshold redn. (%)", 0)
	if err != nil {
		// Try alternative spelling
		tiv.ycol, err = sec.columnContainsName("Threshold change (%)", 0)
		if err != nil {
			return errors.New("Could not get threshold IV: " + err.Error())
		}
	}

	tiv.wasimp = tiv.ycol.ImputeWithValue(curr, tiv.xcol, 0.01, false)

	return nil
}
