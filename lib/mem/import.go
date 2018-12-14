package mem

import (
	"errors"
	"io"
	"log"
	"strconv"
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

func (header *Header) Parse(reader *Reader) error {
	return reader.parseLines(header)
}

func (sr *StimResponse) Parse(reader *Reader) error {
	var err error
	sr.ValueType, err = reader.ReadLineExtractingString(`^Values (.*)`)
	if err != nil {
		return err
	}

	// Find Max CMAP
	err = reader.parseLines(&sr.MaxCmaps)
	if err != nil {
		return err
	}

	err = reader.skipPast("% Max               	Stimulus")
	if err != nil {
		return err
	}

	// Now parse the actual SR data
	return reader.parseLines(sr)
}

func (cd *ChargeDuration) Parse(reader *Reader) error {
	err := reader.skipPast("Duration (ms)       	 Threshold (mA)     	  Threshold charge (mA.mS)")
	if err != nil {
		return err
	}

	return reader.parseLines(cd)
}

func (te *ThresholdElectrotonusGroup) Parse(reader *Reader) error {
	err := reader.skipPast("Delay (ms)          	Current (%)         	Thresh redn. (%)")
	if err != nil {
		return err
	}

	return reader.parseLines(te)
}

func (rc *RecoveryCycle) Parse(reader *Reader) error {
	err := reader.skipPast("Interval (ms)       	  Threshold change (%)")
	if err != nil {
		return err
	}

	return reader.parseLines(rc)
}

func (tiv *ThresholdIV) Parse(reader *Reader) error {
	err := reader.skipPast("Current (%)         	  Threshold redn. (%)")
	if err != nil {
		return err
	}

	return reader.parseLines(tiv)
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

func (section *StrengthDuration) Parse(reader *Reader) error {
	err := reader.skipPast(`%CMAP              	Threshold`)
	if err != nil {
		return err
	}

	return reader.parseLines(section)
}
