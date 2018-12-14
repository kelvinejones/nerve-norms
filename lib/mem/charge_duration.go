package mem

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type ChargeDuration struct {
	Values []XYZ
}

func (section ChargeDuration) Header() string {
	return "CHARGE DURATION DATA"
}

func (cd *ChargeDuration) Parse(reader *Reader) error {
	err := reader.skipPast("Duration (ms)       	 Threshold (mA)     	  Threshold charge (mA.mS)")
	if err != nil {
		return err
	}

	return reader.parseLines(cd)
}

func (cd ChargeDuration) String() string {
	return fmt.Sprintf("ChargeDuration{%d values}", len(cd.Values))
}

func (cd ChargeDuration) LinePrefix() string {
	return "QT"
}

func (cd ChargeDuration) ParseRegex() *regexp.Regexp {
	return regexp.MustCompile(`^QT\.\d+\s+(\d*\.?\d+)\s+(\d*\.?\d+)\s+(\d*\.?\d+)`)
}

func (cd *ChargeDuration) ParseLine(result []string) error {
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
