package mem

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type StrengthDuration struct {
	Values []XY
}

func (section StrengthDuration) Header() string {
	return "STRENGTH-DURATION DATA"
}

func (section *StrengthDuration) Parse(reader *Reader) error {
	err := reader.skipPast(`%CMAP              	Threshold`)
	if err != nil {
		return err
	}

	return reader.parseLines(section)
}

func (sd StrengthDuration) String() string {
	return fmt.Sprintf("StrengthDuration{%d values partially imported}", len(sd.Values))
}

func (sd StrengthDuration) LinePrefix() string {
	return "SD"
}

func (sd StrengthDuration) ParseRegex() *regexp.Regexp {
	return regexp.MustCompile(`^SD\.\d+\s+(\d*\.?\d+)\s+(\d*\.?\d+)\s+(\d*\.?\d+)\s+(\d*\.?\d+)\s+(\d*\.?\d+)`)
}

func (sd *StrengthDuration) ParseLine(result []string) error {
	if len(result) != 6 {
		return errors.New("Incorrect SD line length")
	}

	// Only import the first two columns

	x, err := strconv.ParseFloat(result[1], 64)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(result[2], 64)
	if err != nil {
		return err
	}

	sd.Values = append(sd.Values, XY{
		X: x,
		Y: y,
	})

	return nil
}
