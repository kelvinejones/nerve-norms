package mem

import (
	"errors"
	"fmt"
	"io"
)

type Mem struct {
	Header
	StimResponse
	ChargeDuration
	ThresholdElectrotonusGroup
	RecoveryCycle
	ThresholdIV
	ExcitabilityVariables
	StrengthDuration
}

func Import(data io.Reader) (Mem, error) {
	reader := NewReader(data)
	mem := Mem{}
	mem.ExcitabilityVariables.Values = make(map[string]float64)

	err := mem.Header.Parse(reader)
	if err != nil {
		return mem, err
	}

	for err == nil {
		err = mem.importSection(reader)
	}

	if err != io.EOF {
		return mem, errors.New("Error encountered before EOF: " + err.Error())
	}

	return mem, nil
}

func (mem *Mem) importSection(reader *Reader) error {
	str, err := reader.ReadLine()
	if err != nil {
		return err
	}

	switch {
	case sectionHeaderMatches(&mem.StimResponse, str):
		err = mem.StimResponse.Parse(reader)
		if err != nil {
			return err
		}
	case sectionHeaderMatches(&mem.ChargeDuration, str):
		err = mem.ChargeDuration.Parse(reader)
		if err != nil {
			return err
		}
	case sectionHeaderMatches(&mem.ThresholdElectrotonusGroup, str):
		err = mem.ThresholdElectrotonusGroup.Parse(reader)
		if err != nil {
			return err
		}
	case sectionHeaderMatches(&mem.RecoveryCycle, str):
		err = mem.RecoveryCycle.Parse(reader)
		if err != nil {
			return err
		}
	case sectionHeaderMatches(&mem.ThresholdIV, str):
		err = mem.ThresholdIV.Parse(reader)
		if err != nil {
			return err
		}
	case sectionHeaderMatches(&mem.ExcitabilityVariables, str):
		err = mem.ExcitabilityVariables.Parse(reader)
		if err != nil {
			return err
		}
	case sectionHeaderMatches(&mem.StrengthDuration, str):
		err = mem.StrengthDuration.Parse(reader)
		if err != nil {
			return err
		}
	default:
		fmt.Println("WARNING: Line could not be parsed: " + str)
	}

	return err
}

func (mem Mem) String() string {
	str := "Mem{\n"
	str += "\t" + mem.Header.String() + ",\n"
	str += "\t" + mem.StimResponse.String() + ",\n"
	str += "\t" + mem.ChargeDuration.String() + ",\n"
	str += "\t" + mem.ThresholdElectrotonusGroup.String() + ",\n"
	str += "\t" + mem.RecoveryCycle.String() + ",\n"
	str += "\t" + mem.ThresholdIV.String() + ",\n"
	str += "\t" + mem.ExcitabilityVariables.String() + ",\n"
	str += "\t" + mem.StrengthDuration.String() + ",\n"
	str += "}"
	return str
}
