package mem

import (
	"errors"
)

type ChargeDuration struct{ LabTab }

var CDDuration = Column([]float64{0.2, 0.4, 0.6, 0.8, 1})

func (cd *ChargeDuration) LoadFromMem(mem *rawMem) error {
	cd.xname = "Duration (ms)"
	cd.yname = "Threshold charge (mA•ms)"
	cd.xcol = CDDuration

	sec, err := mem.sectionContainingHeader("CHARGE DURATION")
	if err != nil {
		return errors.New("Could not get charge duration: " + err.Error())
	}

	dur, err := sec.columnContainsName("Duration (ms)", 0)
	if err != nil {
		// For some reason this column sometimes has the wrong name in older files
		dur, err = sec.columnContainsName("Current (%)", 0)
		if err != nil {
			return errors.New("Could not get charge duration: " + err.Error())
		}
	}

	cd.ycol, err = sec.columnContainsName("Threshold charge (mA.mS)", 0)
	if err != nil {
		// Some old formats use this mis-labeled column that must be converted
		threshold, err := sec.columnContainsName("Threshold change (%)", 0)
		if err != nil {
			return errors.New("Could not get charge duration: " + err.Error())
		}

		err = cd.importOldStyle(threshold)
		if err != nil {
			return errors.New("Could not get charge duration: " + err.Error())
		}
	}

	cd.wasimp = cd.ycol.ImputeWithValue(dur, cd.xcol, 0.0000001, false)

	return nil
}

func (cd *ChargeDuration) importOldStyle(threshold Column) error {
	if len(threshold) != len(cd.xcol) {
		return errors.New("Length mis-match in alternative import")
	}
	cd.ycol = Column(make([]float64, len(threshold)))

	for i := range threshold {
		cd.ycol[i] = cd.xcol[i] * threshold[i]
	}
	return nil
}
