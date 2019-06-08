package mem

import (
	"errors"
	"fmt"
)

type ThresholdElectrotonus map[string]*LabTab

func TEDelay(teType string) Column {
	switch teType {
	case "h40":
		return Column([]float64{0, 9, 10, 11, 15, 20, 26, 33, 41, 50, 60, 70, 80, 90, 100, 109, 110, 111, 120, 130, 140, 150, 160, 170, 180, 190, 200, 210})
	default:
		return Column([]float64{0, 9, 10, 11, 20, 30, 40, 50, 60, 70, 80, 90, 100, 109, 110, 111, 120, 130, 140, 150, 160, 170, 180, 190, 200, 210})
	}
}

func teLabTabForType(teType string, tableNum int) *LabTab {
	return &LabTab{
		section:   "THRESHOLD ELECTROTONUS",
		xname:     "Delay (ms)",
		yname:     "Thresh redn. (%)",
		xcol:      TEDelay(teType),
		precision: 0.01,
		tableNum:  tableNum,
	}
}

func newTE() *ThresholdElectrotonus {
	te := ThresholdElectrotonus(make(map[string]*LabTab, 4))
	return &te
}

func (te *ThresholdElectrotonus) LoadFromMem(mem *rawMem) error {
	sec, err := mem.sectionContainingHeader("THRESHOLD ELECTROTONUS")
	if err != nil {
		return MissingSection(errors.New("Could not get threshold electrotonus: " + err.Error()))
	}

	for i := range sec.Tables {
		current, err := sec.columnContainsName("Current (%)", i)
		if err != nil {
			return errors.New("Could not get threshold electrotonus: " + err.Error())
		}
		teType := teTypeForCurrent(current)
		lt := teLabTabForType(teType, i)
		err = lt.LoadFromMem(mem)
		if err != nil {
			if _, ok := err.(MissingSection); ok {
				continue // this is okay, but don't add it
			} else {
				return errors.New("Could not load threshold electrotonus: " + err.Error())
			}
		}
		(*te)[teType] = lt
	}

	return nil
}

func (te *ThresholdElectrotonus) LabelledTable(subsec string) LabelledTable {
	lt, ok := (*te)[subsec]
	if !ok {
		return &emptyLT{}
	}
	return lt
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
	case max == 0 && min < -94 && min > -108:
		return "d100"
	default:
		fmt.Printf("TE contained unexpected current [%f, %f] in %v\n", min, max, current)
		return ""
	}
}
