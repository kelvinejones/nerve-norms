package mem

import (
	"errors"
	"io"
	"log"
	"strings"
)

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
		log.Println("WARNING: Line could not be parsed: " + str)
	}

	return err
}
