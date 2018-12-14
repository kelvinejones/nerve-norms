package mem

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type ThresholdIV struct {
	Values []XY
}

func (section ThresholdIV) Header() []string {
	return []string{"THRESHOLD I/V DATA"}
}

func (tiv *ThresholdIV) Parse(reader *Reader) error {
	err := reader.skipPast("Current (%)         	  Threshold redn. (%)")
	if err != nil {
		return err
	}

	return reader.parseLines(tiv)
}

func (tiv ThresholdIV) String() string {
	return fmt.Sprintf("ThresholdIV{%d values}", len(tiv.Values))
}

func (tiv ThresholdIV) LinePrefix() string {
	return "IV"
}

func (tiv ThresholdIV) ParseRegex() *regexp.Regexp {
	return regexp.MustCompile(`^IV\d+\.\d+\s+([-+]?\d*\.?\d+)\s+([-+]?\d*\.?\d+)`)
}

func (tiv *ThresholdIV) ParseLine(result []string) error {
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
