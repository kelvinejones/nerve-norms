package mem

import (
	"errors"
	"fmt"
	"strconv"
)

type ThresholdElectrotonusSet struct {
	Values []XYZ
}

func (tes ThresholdElectrotonusSet) String() string {
	return fmt.Sprintf("ThresholdElectrotonusSet{%d values}", len(tes.Values))
}

type ThresholdElectrotonusGroup struct {
	Sets []ThresholdElectrotonusSet
}

func (teg ThresholdElectrotonusGroup) String() string {
	str := "ThresholdElectrotonusGroup{"
	for _, tes := range teg.Sets {
		str += tes.String() + ","
	}
	str += "}"
	return str
}

func (section ThresholdElectrotonusGroup) Header() string {
	return "THRESHOLD ELECTROTONUS DATA"
}

func (te *ThresholdElectrotonusGroup) Parse(result []string) error {
	if len(result) != 5 {
		return errors.New("Incorrect TE line length")
	}

	set, err := strconv.Atoi(result[1])
	if err != nil {
		return err
	}
	if set > 100 {
		// Assume this is a parse error
		return errors.New("More than 100 sets of TE data are not supported")
	}

	x, err := strconv.ParseFloat(result[2], 64)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(result[3], 64)
	if err != nil {
		return err
	}
	z, err := strconv.ParseFloat(result[4], 64)
	if err != nil {
		return err
	}

	for len(te.Sets) < set {
		// This would be inefficient for a big difference, but usually this will only run once
		te.Sets = append(te.Sets, ThresholdElectrotonusSet{})
	}

	te.Sets[set-1].Values = append(te.Sets[set-1].Values, XYZ{
		X: x,
		Y: y,
		Z: z,
	})

	return nil
}
