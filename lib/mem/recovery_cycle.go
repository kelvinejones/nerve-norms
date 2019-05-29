package mem

import (
	"encoding/json"
	"errors"
)

type RecoveryCycle struct {
	Interval     Column
	ThreshChange Column
	WasImputed   Column
}

func (mem *Mem) RecoveryCycle() (RecoveryCycle, error) {
	rc := RecoveryCycle{Interval: Column([]float64{2, 2.5, 3.2, 4, 5, 6.3, 7.9, 10, 13, 18, 24, 32, 42, 56, 75, 100, 140, 200})}

	sec, err := mem.sectionContainingHeader("RECOVERY CYCLE")
	if err != nil {
		return rc, errors.New("Could not get recovery cycle: " + err.Error())
	}

	interval, err := sec.columnContainsName("Interval (ms)", 0)
	if err != nil {
		return rc, errors.New("Could not get recovery cycle: " + err.Error())
	}

	rc.ThreshChange, err = sec.columnContainsName("Threshold change (%)", 0)
	if err != nil {
		return rc, errors.New("Could not get recovery cycle: " + err.Error())
	}

	rc.WasImputed = rc.ThreshChange.ImputeWithValue(interval, rc.Interval, 0.000001)

	return rc, nil
}

func (dat *RecoveryCycle) MarshalJSON() ([]byte, error) {
	str := &struct {
		Columns []string `json:"columns"`
		Data    Table    `json:"data"`
	}{
		Columns: []string{"Interval (ms)", "Threshold change (%)"},
		Data:    []Column{dat.Interval, dat.ThreshChange},
	}

	if dat.WasImputed != nil {
		str.Columns = append(str.Columns, "Was Imputed")
		str.Data = append(str.Data, dat.WasImputed)
	}

	return json.Marshal(&str)
}
