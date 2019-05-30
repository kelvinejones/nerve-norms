package mem

import (
	"encoding/json"
	"errors"
	"fmt"
)

type TEPair struct {
	Delay           Column
	ThreshReduction Column
	WasImputed      Column
}

type ThresholdElectrotonus struct {
	Data map[string]*TEPair
}

var TEDelay = Column([]float64{0, 9, 10, 11, 15, 20, 26, 30, 33, 30, 41, 50, 60, 70, 80, 90, 100, 109, 110, 111, 120, 130, 140, 150, 160, 170, 180, 190, 200, 210})

func (te *ThresholdElectrotonus) LoadFromMem(mem *rawMem) error {
	sec, err := mem.sectionContainingHeader("THRESHOLD ELECTROTONUS")
	if err != nil {
		return errors.New("Could not get threshold electrotonus: " + err.Error())
	}

	te.Data = make(map[string]*TEPair, 4)

	for i := range sec.Tables {
		current, err := sec.columnContainsName("Current (%)", i)
		if err != nil {
			return errors.New("Could not get threshold electrotonus: " + err.Error())
		}
		teType := teTypeForCurrent(current)

		pair := TEPair{Delay: TEDelay}

		delay, err := sec.columnContainsName("Delay (ms)", i)
		if err != nil {
			return errors.New("Could not get threshold electrotonus: " + err.Error())
		}

		pair.ThreshReduction, err = sec.columnContainsName("Thresh redn. (%)", i)
		if err != nil {
			return errors.New("Could not get threshold electrotonus: " + err.Error())
		}

		pair.WasImputed = pair.ThreshReduction.ImputeWithValue(delay, pair.Delay, 0.01, false)

		te.Data[teType] = &pair
	}

	return nil
}

func teTypeForCurrent(current Column) string {
	// This is a quick and simple way to parse the data we expect to see
	max := current.Maximum()
	min := current.Minimum()
	switch {
	case max > 34 && max < 46 && min == 0:
		return "h40"
	case max > 14 && max < 26 && min == 0:
		return "h20"
	case max == 0 && min < -34 && min > -46:
		return "d40"
	case max == 0 && min < -14 && min > -26:
		return "d20"
	case max == 0 && min < -64 && min > -76:
		return "d70"
	case max == 0 && min < -94 && min > -106:
		return "d100"
	default:
		fmt.Printf("TE contained unexpected current [%f, %f]\n", min, max)
		return ""
	}
}

type jsonThresholdElectrotonus struct {
	Columns []string         `json:"columns"`
	Data    map[string]Table `json:"data"`
}

func (dat *ThresholdElectrotonus) MarshalJSON() ([]byte, error) {
	str := &jsonThresholdElectrotonus{
		Columns: []string{"Delay (ms)", "Threshold Reduction (%)"},
		Data:    make(map[string]Table, 4),
	}

	somethingWasImputed := false
	for key, val := range dat.Data {
		str.Data[key] = []Column{val.Delay, val.ThreshReduction}
		if val.WasImputed != nil {
			somethingWasImputed = true
		}
	}

	if somethingWasImputed {
		str.Columns = append(str.Columns, "Was Imputed")
		for key, val := range dat.Data {
			str.Data[key] = append(str.Data[key], val.WasImputed)
		}
	}

	return json.Marshal(&str)
}

func (dat *ThresholdElectrotonus) UnmarshalJSON(value []byte) error {
	jsDat := jsonThresholdElectrotonus{}
	err := json.Unmarshal(value, &jsDat)
	if err != nil {
		return err
	}
	numCol := len(jsDat.Columns)

	if numCol < 2 || numCol > 3 {
		return errors.New("Incorrect number of ThresholdElectrotonus columns in JSON")
	}
	if jsDat.Columns[0] != "Delay (ms)" || jsDat.Columns[1] != "Threshold Reduction (%)" || (numCol == 3 && jsDat.Columns[2] != "Was Imputed") {
		return errors.New("Incorrect ThresholdElectrotonus column names in JSON")
	}

	dat.Data = make(map[string]*TEPair, 4)
	for key := range jsDat.Data {
		dat.Data[key] = tePairFromTable(jsDat.Data[key])
	}

	return nil
}

func tePairFromTable(tab Table) *TEPair {
	tep := &TEPair{}
	tep.Delay = tab[0]
	tep.ThreshReduction = tab[1]
	if len(tab) == 3 {
		tep.WasImputed = tab[2]
	}
	return tep
}
