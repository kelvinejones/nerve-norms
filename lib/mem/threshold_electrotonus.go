package mem

import (
	"errors"
	"fmt"
)

type TEPair struct {
	Delay           Column
	ThreshReduction Column
}

type ThresholdElectrotonus struct {
	Hyperpol40 *TEPair
	Hyperpol20 *TEPair
	Depol40    *TEPair
	Depol20    *TEPair
}

func (mem *Mem) ThresholdElectrotonus() (ThresholdElectrotonus, error) {
	te := ThresholdElectrotonus{}

	sec, err := mem.sectionContainingHeader("THRESHOLD ELECTROTONUS")
	if err != nil {
		return te, errors.New("Could not get threshold electrotonus: " + err.Error())
	}

	for i := range sec.Tables {
		pair := TEPair{}

		pair.Delay, err = sec.columnContainsName("Delay (ms)", i)
		if err != nil {
			return te, errors.New("Could not get threshold electrotonus: " + err.Error())
		}

		pair.ThreshReduction, err = sec.columnContainsName("Thresh redn. (%)", i)
		if err != nil {
			return te, errors.New("Could not get threshold electrotonus: " + err.Error())
		}

		current, err := sec.columnContainsName("Current (%)", i)
		if err != nil {
			return te, errors.New("Could not get threshold electrotonus: " + err.Error())
		}

		// This is a quick and simple way to parse the data we expect to see
		max := current.Maximum()
		min := current.Minimum()
		switch {
		case max == 40 && min == 0:
			te.Hyperpol40 = &pair
		case max == 20 && min == 0:
			te.Hyperpol20 = &pair
		case max == 0 && min == -40:
			te.Depol40 = &pair
		case max == 0 && min == -20:
			te.Depol20 = &pair
		default:
			fmt.Printf("TE contained unexpected current [%f, %f]\n", min, max)
		}
	}

	return te, nil
}
