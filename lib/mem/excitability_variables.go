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

func (exciteVar *ExcitabilityVariables) Parse(reader *Reader) error {
	// Find settings
	err := reader.skipNewlines()
	if err != nil {
		return err
	}

	exciteVar.Program, err = reader.ReadLineExtractingString(`^Program = (.*)`)
	if err != nil {
		return err
	}

	val, err := reader.ReadLineExtractingString(`^Threshold method = (\d+).*`)
	if err != nil {
		return err
	}
	exciteVar.ThresholdMethod, err = strconv.Atoi(val)
	if err != nil {
		return err
	}

	val, err = reader.ReadLineExtractingString(`^SR method = (\d+).*`)
	if err != nil {
		return err
	}
	exciteVar.SRMethod, err = strconv.Atoi(val)
	if err != nil {
		return err
	}

	// Read the main variables
	err = reader.parseLines(exciteVar)
	if err != nil {
		return err
	}

	// Now find any extra variables
	err = reader.skipNewlines()
	if err != nil {
		return err
	}
	err = reader.skipPast(ExtraVariables{}.Header())
	if err != nil {
		return err
	}
	err = reader.skipNewlines()
	if err != nil {
		return err
	}

	err = reader.parseLines(&ExtraVariables{exciteVar})
	return err
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
