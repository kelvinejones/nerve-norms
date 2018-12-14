package mem

import (
	"errors"
	"fmt"
	"io"
	"strings"
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
	case strings.Contains(str, mem.StimResponse.Header()):
		err = mem.StimResponse.Parse(reader)
		if err != nil {
			return err
		}
	case strings.Contains(str, mem.ChargeDuration.Header()):
		err = mem.ChargeDuration.Parse(reader)
		if err != nil {
			return err
		}
	case strings.Contains(str, mem.ThresholdElectrotonusGroup.Header()):
		err = mem.ThresholdElectrotonusGroup.Parse(reader)
		if err != nil {
			return err
		}
	case strings.Contains(str, mem.RecoveryCycle.Header()):
		err = mem.RecoveryCycle.Parse(reader)
		if err != nil {
			return err
		}
	case strings.Contains(str, mem.ThresholdIV.Header()):
		err = mem.ThresholdIV.Parse(reader)
		if err != nil {
			return err
		}
	case strings.Contains(str, mem.ExcitabilityVariables.Header()):
		err = mem.ExcitabilityVariables.Parse(reader)
		if err != nil {
			return err
		}
	case strings.Contains(str, mem.StrengthDuration.Header()):
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
