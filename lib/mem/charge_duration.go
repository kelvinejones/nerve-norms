package mem

import (
	"errors"
)

type ChargeDuration struct{ LabelledTable }

var CDDuration = Column([]float64{0.2, 0.4, 0.6, 0.8, 1})

func (cd *ChargeDuration) LoadFromMem(mem *rawMem) error {
	cd.XName = "Duration (ms)"
	cd.YName = "Threshold charge (mA•ms)"
	cd.XColumn = CDDuration

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

	cd.YColumn, err = sec.columnContainsName("Threshold charge (mA.mS)", 0)
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

	cd.WasImputed = cd.YColumn.ImputeWithValue(dur, cd.XColumn, 0.0000001, false)

	return nil
}

func (cd *ChargeDuration) importOldStyle(threshold Column) error {
	if len(threshold) != len(cd.XColumn) {
		return errors.New("Length mis-match in alternative import")
	}
	cd.YColumn = Column(make([]float64, len(threshold)))

	for i := range threshold {
		cd.YColumn[i] = cd.XColumn[i] * threshold[i]
	}
	return nil
}
