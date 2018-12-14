package mem

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type ExcitabilityVariables struct {
	Values          map[string]float64
	Program         string
	ThresholdMethod int
	SRMethod        int
}

func (section ExcitabilityVariables) Header() string {
	return "DERIVED EXCITABILITY VARIABLES"
}

func (ev ExcitabilityVariables) String() string {
	return fmt.Sprintf("ExcitabilityVariables{%d values}", len(ev.Values))
}

func (exciteVar ExcitabilityVariables) ParseRegex() *regexp.Regexp {
	return regexp.MustCompile(`^ \d+\.\s+([-+]?\d*\.?\d+)\s+(.+)`)
}

func (exciteVar *ExcitabilityVariables) ParseLine(result []string) error {
	if len(result) != 3 {
		return errors.New("Incorrect ExVar line length")
	}

	val, err := strconv.ParseFloat(result[1], 64)
	if err != nil {
		return err
	}

	exciteVar.Values[result[2]] = val

	return nil
}
