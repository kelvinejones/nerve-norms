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

var RCInterval = Column([]float64{2, 2.5, 3.2, 4, 5, 6.3, 7.9, 10, 13, 18, 24, 32, 42, 56, 75, 100, 140, 200})

func (rc *RecoveryCycle) LoadFromMem(mem *rawMem) error {
	rc.Interval = RCInterval

	sec, err := mem.sectionContainingHeader("RECOVERY CYCLE")
	if err != nil {
		return errors.New("Could not get recovery cycle: " + err.Error())
	}

	interval, err := sec.columnContainsName("Interval (ms)", 0)
	if err != nil {
		return errors.New("Could not get recovery cycle: " + err.Error())
	}

	rc.ThreshChange, err = sec.columnContainsName("Threshold change (%)", 0)
	if err != nil {
		return errors.New("Could not get recovery cycle: " + err.Error())
	}

	rc.WasImputed = rc.ThreshChange.ImputeWithValue(interval, rc.Interval, 0.000001, true)

	return nil
}

type jsonRecoveryCycle struct {
	Columns []string `json:"columns"`
	Data    Table    `json:"data"`
}

func (dat *RecoveryCycle) MarshalJSON() ([]byte, error) {
	str := &jsonRecoveryCycle{
		Columns: []string{"Interval (ms)", "Threshold change (%)"},
		Data:    []Column{dat.Interval, dat.ThreshChange},
	}

	if dat.WasImputed != nil {
		str.Columns = append(str.Columns, "Was Imputed")
		str.Data = append(str.Data, dat.WasImputed)
	}

	return json.Marshal(&str)
}

func (dat *RecoveryCycle) UnmarshalJSON(value []byte) error {
	jsDat := jsonRecoveryCycle{}
	err := json.Unmarshal(value, &jsDat)
	if err != nil {
		return err
	}
	numCol := len(jsDat.Columns)

	if numCol < 2 || numCol > 3 {
		return errors.New("Incorrect number of RecoveryCycle columns in JSON")
	}
	if jsDat.Columns[0] != "Interval (ms)" || jsDat.Columns[1] != "Threshold change (%)" || (numCol == 3 && jsDat.Columns[2] != "Was Imputed") {
		return errors.New("Incorrect RecoveryCycle column names in JSON")
	}

	dat.Interval = jsDat.Data[0]
	dat.ThreshChange = jsDat.Data[1]
	if numCol == 3 {
		dat.WasImputed = jsDat.Data[2]
	}

	return nil
}
