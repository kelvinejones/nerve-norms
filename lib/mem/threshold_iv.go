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

func (tiv *ThresholdIV) LoadFromMem(mem *rawMem) error {
	tiv.Current = Column([]float64{50, 40, 30, 20, 10, 0, -10, -20, -30, -40, -50, -60, -70, -80, -90, -100})

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

func (dat *ThresholdIV) MarshalJSON() ([]byte, error) {
	str := &struct {
		Columns []string `json:"columns"`
		Data    Table    `json:"data"`
	}{
		Columns: []string{"Current (%)", "Threshold Reduction (%)"},
		Data:    []Column{dat.Current, dat.ThreshReduction},
	}

	if dat.WasImputed != nil {
		str.Columns = append(str.Columns, "Was Imputed")
		str.Data = append(str.Data, dat.WasImputed)
	}

	return json.Marshal(&str)
}
