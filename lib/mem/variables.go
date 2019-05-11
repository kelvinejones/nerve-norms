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
	Values map[int]float64
	ExcitabilitySettings
}

func (exciteVar *ExcitabilityVariables) Parse(reader *Reader) error {
	if exciteVar.Values == nil {
		exciteVar.Values = map[int]float64{}
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
	str, err := reader.skipEmptyLines()
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
	return regexp.MustCompile(`^ (\d+)\.\s+([-+]?\d*\.?\d+)\s+.+`)
}

func (exciteVar *ExcitabilityVariables) ParseLine(result []string) error {
	if len(result) != 3 {
		return errors.New("Incorrect ExVar line length")
	}

	val, err := strconv.ParseFloat(result[2], 64)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(result[1])
	if err != nil {
		return err
	}

	// Since we have an ID, just store that as the name.
	exciteVar.Values[id] = val

	return nil
}

type ExtraVariables struct {
	*ExcitabilityVariables
}

func (extraVar ExtraVariables) ParseRegex() *regexp.Regexp {
	return regexp.MustCompile(`^(.+)\s+=\s+([-+]?\d*\.?\d+)`)
}

func (extraVar *ExtraVariables) ParseLine(result []string) error {
	if len(result) != 3 {
		return errors.New("Incorrect ExtraVar line length")
	}

	val, err := strconv.ParseFloat(result[2], 64)
	if err != nil {
		return err
	}

	id := idForName(strings.TrimSpace(result[1]))
	if id < 1000 {
		return errors.New("Invalid name '" + result[1] + "'")
	}
	extraVar.Values[id] = val

	return nil
}

func idForName(name string) int {
	switch name {
	case "TEd40(Accom)":
		return 1001
	case "TEd20(10-20ms)":
		return 1002
	case "TEh20(10-20ms)":
		return 1003
	case "TEh20(90-100ms)":
		return 1004
	case "MScPeak(mV)":
		return 1005
	case "MScS50(mA)":
		return 1006
	case "MSFNUnits":
		return 1007
	case "MSFMeanUnitAmp(uV)":
		return 1008
	case "MRRP":
		return 1009
	case "MSuperN(<15)%":
		return 1010
	case "MSuperN@(ms)":
		return 1011
	case "CRRP(ms)":
		return 1012
	case "CSuperN(%)":
		return 1013
	case "CSuperN@(ms)":
		return 1014
	case "RMT200":
		return 1015
	case "T-SICI(70%)2.5ms":
		return 1016
	case "T-ICF(70%)15ms":
		return 1017
	case "Potassium":
		return 1018
	case "pH":
		return 1019
	case "MRCsumscore":
		return 1020
	default:
		return -1
	}

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
