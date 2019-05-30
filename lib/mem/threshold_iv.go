package mem

import (
	"encoding/json"
	"errors"
)

type ThresholdIV struct {
	Current         Column
	ThreshReduction Column
	WasImputed      Column
}

var IVCurrent = Column([]float64{50, 40, 30, 20, 10, 0, -10, -20, -30, -40, -50, -60, -70, -80, -90, -100})

func (tiv *ThresholdIV) LoadFromMem(mem *rawMem) error {
	tiv.Current = IVCurrent

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

	tiv.ThreshReduction, err = sec.columnContainsName("Threshold redn. (%)", 0)
	if err != nil {
		// Try alternative spelling
		tiv.ThreshReduction, err = sec.columnContainsName("Threshold change (%)", 0)
		if err != nil {
			return errors.New("Could not get threshold IV: " + err.Error())
		}
	}

	tiv.WasImputed = tiv.ThreshReduction.ImputeWithValue(curr, tiv.Current, 0.01, false)

	return nil
}

type jsonThresholdIV struct {
	Columns []string `json:"columns"`
	Data    Table    `json:"data"`
}

func (dat *ThresholdIV) MarshalJSON() ([]byte, error) {
	str := &jsonThresholdIV{
		Columns: []string{"Current (%)", "Threshold Reduction (%)"},
		Data:    []Column{dat.Current, dat.ThreshReduction},
	}

	if dat.WasImputed != nil {
		str.Columns = append(str.Columns, "Was Imputed")
		str.Data = append(str.Data, dat.WasImputed)
	}

	return json.Marshal(&str)
}

func (dat *ThresholdIV) UnmarshalJSON(value []byte) error {
	jsDat := jsonThresholdIV{}
	err := json.Unmarshal(value, &jsDat)
	if err != nil {
		return err
	}
	numCol := len(jsDat.Columns)

	if numCol < 2 || numCol > 3 {
		return errors.New("Incorrect number of ThresholdIV columns in JSON")
	}
	if jsDat.Columns[0] != "Current (%)" || jsDat.Columns[1] != "Threshold Reduction (%)" || (numCol == 3 && jsDat.Columns[2] != "Was Imputed") {
		return errors.New("Incorrect ThresholdIV column names in JSON")
	}

	dat.Current = jsDat.Data[0]
	dat.ThreshReduction = jsDat.Data[1]
	if numCol == 3 {
		dat.WasImputed = jsDat.Data[2]
	}

	return nil
}
