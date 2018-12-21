package mem

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ExcitabilitySettings map[string]string

type ExcitabilityVariables struct {
	Values map[string]float64
	ExcitabilitySettings
}

func (exciteVar *ExcitabilityVariables) Parse(reader *Reader) error {
	if exciteVar.Values == nil {
		exciteVar.Values = map[string]float64{}
	}
	if exciteVar.ExcitabilitySettings == nil {
		exciteVar.ExcitabilitySettings = map[string]string{}
	}
	// Until a line matches the regex, allow parsing of other things
	err := reader.parseLines(&exciteVar.ExcitabilitySettings)
	if err != nil {
		return err
	}

	// Read the main variables
	err = reader.parseLines(exciteVar)
	if err != nil {
		return err
	}

	// Now find any extra variables
	str, err := reader.skipNewlines()
	if err != nil {
		return err
	}

	if strings.Contains(str, "EXTRA VARIABLES") {
		err = reader.parseLines(&ExtraVariables{exciteVar})
	} else {
		// It looks like this header doesn't belong to us, so give it back
		reader.UnreadString(str)
	}

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

type ExtraVariables struct {
	*ExcitabilityVariables
}

func (extraVar ExtraVariables) ParseRegex() *regexp.Regexp {
	return regexp.MustCompile(`^(.+) = ([-+]?\d*\.?\d+)`)
}

func (extraVar *ExtraVariables) ParseLine(result []string) error {
	if len(result) != 3 {
		return errors.New("Incorrect ExtraVar line length")
	}

	val, err := strconv.ParseFloat(result[2], 64)
	if err != nil {
		return err
	}

	extraVar.Values[result[1]] = val

	return nil
}

func (es ExcitabilitySettings) ParseRegex() *regexp.Regexp {
	return regexp.MustCompile(`^(.+) = (.+)`)
}

func (es *ExcitabilitySettings) ParseLine(result []string) error {
	if len(result) != 3 {
		return errors.New("Incorrect ExcitabilitySettings line length")
	}

	(*es)[result[1]] = result[2]

	return nil
}
