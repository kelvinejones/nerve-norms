package mem

import (
	"errors"
	"fmt"
)

type ThresholdElectrotonus map[string]*LabelledTable

func TEDelay(teType string) Column {
	switch teType {
	case "h40":
		return Column([]float64{0, 9, 10, 11, 15, 20, 26, 33, 41, 50, 60, 70, 80, 90, 100, 109, 110, 111, 120, 130, 140, 150, 160, 170, 180, 190, 200, 210})
	default:
		return Column([]float64{0, 9, 10, 11, 20, 30, 40, 50, 60, 70, 80, 90, 100, 109, 110, 111, 120, 130, 140, 150, 160, 170, 180, 190, 200, 210})
	}
}

func (te *ThresholdElectrotonus) LoadFromMem(mem *rawMem) error {
	sec, err := mem.sectionContainingHeader("THRESHOLD ELECTROTONUS")
	if err != nil {
		return errors.New("Could not get threshold electrotonus: " + err.Error())
	}

	*te = make(map[string]*LabelledTable, 4)

	for i := range sec.Tables {
		current, err := sec.columnContainsName("Current (%)", i)
		if err != nil {
			return errors.New("Could not get threshold electrotonus: " + err.Error())
		}
		teType := teTypeForCurrent(current)

		pair := LabelledTable{}
		pair.XName = "Delay (ms)"
		pair.YName = "Threshold Reduction (%)"
		pair.XColumn = TEDelay(teType)

		delay, err := sec.columnContainsName("Delay (ms)", i)
		if err != nil {
			return errors.New("Could not get threshold electrotonus: " + err.Error())
		}

		pair.YColumn, err = sec.columnContainsName("Thresh redn. (%)", i)
		if err != nil {
			return errors.New("Could not get threshold electrotonus: " + err.Error())
		}

		pair.WasImputed = pair.YColumn.ImputeWithValue(delay, pair.XColumn, 0.01, false)

		(*te)[teType] = &pair
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
