package mem

import (
	"errors"
	"fmt"
	"strconv"
)

type ChargeDuration struct {
	Values []XYZ
}

func (section ChargeDuration) Header() string {
	return "CHARGE DURATION DATA"
}

func (cd ChargeDuration) String() string {
	return fmt.Sprintf("ChargeDuration{%d values}", len(cd.Values))
}

func (cd *ChargeDuration) Parse(result []string) error {
	if len(result) != 4 {
		return errors.New("Incorrect CD line length")
	}

	x, err := strconv.ParseFloat(result[1], 64)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(result[2], 64)
	if err != nil {
		return err
	}
	z, err := strconv.ParseFloat(result[3], 64)
	if err != nil {
		return err
	}

	cd.Values = append(cd.Values, XYZ{
		X: x,
		Y: y,
		Z: z,
	})

	return nil
}
