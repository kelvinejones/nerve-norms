package mem

import (
	"errors"
	"fmt"
	"strconv"
)

type ThresholdIV struct {
	Values []XY
}

func (section ThresholdIV) Header() string {
	return "THRESHOLD I/V DATA"
}

func (tiv ThresholdIV) String() string {
	return fmt.Sprintf("ThresholdIV{%d values}", len(tiv.Values))
}

func (tiv *ThresholdIV) Parse(result []string) error {
	if len(result) != 3 {
		return errors.New("Incorrect TIV line length")
	}

	x, err := strconv.ParseFloat(result[1], 64)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(result[2], 64)
	if err != nil {
		return err
	}

	tiv.Values = append(tiv.Values, XY{
		X: x,
		Y: y,
	})

	return nil
}
