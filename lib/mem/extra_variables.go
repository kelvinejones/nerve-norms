package mem

import (
	"errors"
	"strconv"
)

type ExtraVariables struct {
	*ExcitabilityVariables
}

func (section ExtraVariables) Header() string {
	return "EXTRA VARIABLES"
}

func (extraVar *ExtraVariables) Parse(result []string) error {
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
