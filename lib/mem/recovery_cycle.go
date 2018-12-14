package mem

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type RecoveryCycle struct {
	Values []XY
}

func (section RecoveryCycle) Header() string {
	return "RECOVERY CYCLE DATA"
}

func (rc *RecoveryCycle) Parse(reader *Reader) error {
	err := reader.skipPast("Interval (ms)       	  Threshold change (%)")
	if err != nil {
		return err
	}

	return reader.parseLines(rc)
}

func (rc RecoveryCycle) String() string {
	return fmt.Sprintf("RecoveryCycle{%d values}", len(rc.Values))
}

func (rc RecoveryCycle) ParseRegex() *regexp.Regexp {
	return regexp.MustCompile(`^RC\d+\.\d+\s+(\d*\.?\d+)\s+([-+]?\d*\.?\d+)`)
}

func (rc *RecoveryCycle) ParseLine(result []string) error {
	if len(result) != 3 {
		return errors.New("Incorrect RC line length")
	}

	x, err := strconv.ParseFloat(result[1], 64)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(result[2], 64)
	if err != nil {
		return err
	}

	rc.Values = append(rc.Values, XY{
		X: x,
		Y: y,
	})

	return nil
}
