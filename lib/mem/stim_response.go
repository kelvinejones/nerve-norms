package mem

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type MaxCmap struct {
	Val   float64
	Time  float64
	Units byte
}

type MaxCmaps []MaxCmap

func (cmaps MaxCmaps) ParseRegex() *regexp.Regexp {
	return regexp.MustCompile(`^ Max CMAP  (\d*\.?\d+) ms =  (\d*\.?\d+) (.)V`)
}

func (cmaps *MaxCmaps) ParseLine(result []string) error {
	if len(result) != 4 {
		return errors.New("Incorrect CMAP line length")
	}

	cmap := MaxCmap{}
	var err error
	cmap.Time, err = strconv.ParseFloat(result[1], 64)
	if err != nil {
		return err
	}
	cmap.Val, err = strconv.ParseFloat(result[2], 64)
	if err != nil {
		return err
	}
	cmap.Units = result[3][0]

	*cmaps = append(*cmaps, cmap)

	return nil
}

type StimResponse struct {
	MaxCmaps
	Values    []XY
	ValueType string
}

func (section StimResponse) Header() string {
	return "STIMULUS-RESPONSE DATA"
}

func (sr StimResponse) String() string {
	return fmt.Sprintf("StimResponse{%d MaxCmaps, %d values}", len(sr.MaxCmaps), len(sr.Values))
}

func (sr StimResponse) ParseRegex() *regexp.Regexp {
	return regexp.MustCompile(`^SR\.(\d+)\s+(\d+)\s+(\d*\.?\d+)`)
}

func (sr *StimResponse) ParseLine(result []string) error {
	if len(result) != 4 {
		return errors.New("Incorrect SR line length")
	}
	if result[1] != result[2] {
		return errors.New("SR fields do not match")
	}

	percentMax, err := strconv.ParseFloat(result[1], 64)
	if err != nil {
		return err
	}
	stim, err := strconv.ParseFloat(result[3], 64)
	if err != nil {
		return err
	}

	sr.Values = append(sr.Values, XY{
		X: percentMax,
		Y: stim,
	})

	return nil
}
