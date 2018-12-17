package mem

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type StrengthDuration struct {
	Values []XY
}

func (section StrengthDuration) Header() []string {
	return []string{
		"STRENGTH-DURATION DATA",
		"STRENGTH DURATION DATA",
	}
}

func (section *StrengthDuration) Parse(reader *Reader) error {
	s, err := reader.skipNewlines()
	if err != nil {
		return err
	}

	if !strings.Contains(s, "%CMAP") || !strings.Contains(s, "Threshold") {
		return errors.New("Could not find '%CMAP     Threshold' header line")
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
	return regexp.MustCompile(`^SD\.\d+\s+(\d*\.?\d+)\s+(\d*\.?\d+)`) // The line might be longer, but we don't care
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
