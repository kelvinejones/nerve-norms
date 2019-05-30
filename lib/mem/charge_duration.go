package mem

import (
	"encoding/json"
	"errors"
)

type ChargeDuration struct {
	Duration     Column
	ThreshCharge Column
	WasImputed   Column
}

func (cd *ChargeDuration) LoadFromMem(mem *rawMem) error {
	cd.Duration = Column([]float64{0.2, 0.4, 0.6, 0.8, 1})

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

	cd.ThreshCharge, err = sec.columnContainsName("Threshold charge (mA.mS)", 0)
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

	cd.WasImputed = cd.ThreshCharge.ImputeWithValue(dur, cd.Duration, 0.0000001, false)

	return nil
}

type jsonChargeDuration struct {
	Columns []string `json:"columns"`
	Data    Table    `json:"data"`
}

func (dat *ChargeDuration) MarshalJSON() ([]byte, error) {
	str := &jsonChargeDuration{
		Columns: []string{"Duration (ms)", "Threshold charge (mA•ms)"},
		Data:    []Column{dat.Duration, dat.ThreshCharge},
	}

	if dat.WasImputed != nil {
		str.Columns = append(str.Columns, "Was Imputed")
		str.Data = append(str.Data, dat.WasImputed)
	}

	return json.Marshal(&str)
}

func (dat *ChargeDuration) UnmarshalJSON(value []byte) error {
	jsDat := jsonChargeDuration{}
	err := json.Unmarshal(value, &jsDat)
	if err != nil {
		return err
	}
	numCol := len(jsDat.Columns)

	if numCol < 2 || numCol > 3 {
		return errors.New("Incorrect number of ChargeDuration columns in JSON")
	}
	if jsDat.Columns[0] != "Duration (ms)" || jsDat.Columns[1] != "Threshold charge (mA•ms)" || (numCol == 3 && jsDat.Columns[2] != "Was Imputed") {
		return errors.New("Incorrect ChargeDuration column names in JSON")
	}

	dat.Duration = jsDat.Data[0]
	dat.ThreshCharge = jsDat.Data[1]
	if numCol == 3 {
		dat.WasImputed = jsDat.Data[2]
	}

	return nil
}

func (cd *ChargeDuration) importOldStyle(threshold Column) error {
	if len(threshold) != len(cd.Duration) {
		return errors.New("Length mis-match in alternative import")
	}
	cd.ThreshCharge = Column(make([]float64, len(threshold)))

	for i := range threshold {
		cd.ThreshCharge[i] = cd.Duration[i] * threshold[i]
	}
	return nil
}
