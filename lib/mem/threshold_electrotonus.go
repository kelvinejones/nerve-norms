package mem

import (
	"errors"
	"fmt"
)

type TEPair struct {
	Delay           Column
	ThreshReduction Column
	WasImputed      Column
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
		pair := TEPair{Delay: Column([]float64{0, 9, 10, 11, 15, 20, 26, 33, 41, 50, 60, 70, 80, 90, 100, 109, 110, 111, 120, 130, 140, 150, 160, 170, 180, 190, 200, 210})}

		delay, err := sec.columnContainsName("Delay (ms)", i)
		if err != nil {
			return te, errors.New("Could not get threshold electrotonus: " + err.Error())
		}

		pair.ThreshReduction, err = sec.columnContainsName("Thresh redn. (%)", i)
		if err != nil {
			return te, errors.New("Could not get threshold electrotonus: " + err.Error())
		}

		old := pair.ThreshReduction
		pair.WasImputed = pair.ThreshReduction.ImputeWithValue(delay, pair.Delay, 0.01)
		if pair.WasImputed != nil {
			fmt.Println("Imputed TE:", old, pair.ThreshReduction)
		}

		current, err := sec.columnContainsName("Current (%)", i)
		if err != nil {
			return te, errors.New("Could not get threshold electrotonus: " + err.Error())
		}

		// This is a quick and simple way to parse the data we expect to see
		max := current.Maximum()
		min := current.Minimum()
		switch {
		case max > 38 && max < 42 && min == 0:
			te.Hyperpol40 = &pair
		case max > 18 && max < 22 && min == 0:
			te.Hyperpol20 = &pair
		case max == 0 && min < -38 && min > -42:
			te.Depol40 = &pair
		case max == 0 && min < -18 && min > -22:
			te.Depol20 = &pair
		default:
			fmt.Printf("TE contained unexpected current [%f, %f]\n", min, max)
		}
	}

	return te, nil
}
