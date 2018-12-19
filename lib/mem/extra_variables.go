package mem

import (
	"errors"
	"regexp"
	"strconv"
)

type ExtraVariables struct {
	*ExcitabilityVariables
}

func (extraVar ExtraVariables) LinePrefix() string {
	return ""
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
