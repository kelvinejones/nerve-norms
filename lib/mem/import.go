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

	err := parseHeader(reader, &mem.Header)
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
		err = parseStimResponse(reader, &mem.StimResponse)
		if err != nil {
			return err
		}
	case strings.Contains(str, mem.ChargeDuration.Header()):
		err = parseChargeDuration(reader, &mem.ChargeDuration)
		if err != nil {
			return err
		}
	case strings.Contains(str, mem.ThresholdElectrotonusGroup.Header()):
		err = parseThresholdElectrotonus(reader, &mem.ThresholdElectrotonusGroup)
		if err != nil {
			return err
		}
	case strings.Contains(str, mem.RecoveryCycle.Header()):
		err = parseRecoveryCycle(reader, &mem.RecoveryCycle)
		if err != nil {
			return err
		}
	case strings.Contains(str, mem.ThresholdIV.Header()):
		err = parseThresholdIV(reader, &mem.ThresholdIV)
		if err != nil {
			return err
		}
	case strings.Contains(str, mem.ExcitabilityVariables.Header()):
		err = parseExcitabilityVariables(reader, &mem.ExcitabilityVariables)
		if err != nil {
			return err
		}
	case strings.Contains(str, mem.StrengthDuration.Header()):
		err = parseStrengthDuration(reader, &mem.StrengthDuration)
		if err != nil {
			return err
		}
	default:
		log.Println("WARNING: Line could not be parsed: " + str)
	}

	return err
}

func parseHeader(reader *Reader, header *Header) error {
	return reader.parseLines(header)
}

func parseStimResponse(reader *Reader, sr *StimResponse) error {
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

func parseChargeDuration(reader *Reader, cd *ChargeDuration) error {
	err := reader.skipPast("Duration (ms)       	 Threshold (mA)     	  Threshold charge (mA.mS)")
	if err != nil {
		return err
	}

	return reader.parseLines(cd)
}

func parseThresholdElectrotonus(reader *Reader, te *ThresholdElectrotonusGroup) error {
	err := reader.skipPast("Delay (ms)          	Current (%)         	Thresh redn. (%)")
	if err != nil {
		return err
	}

	return reader.parseLines(te)
}

func parseRecoveryCycle(reader *Reader, rc *RecoveryCycle) error {
	err := reader.skipPast("Interval (ms)       	  Threshold change (%)")
	if err != nil {
		return err
	}

	return reader.parseLines(rc)
}

func parseThresholdIV(reader *Reader, tiv *ThresholdIV) error {
	err := reader.skipPast("Current (%)         	  Threshold redn. (%)")
	if err != nil {
		return err
	}

	return reader.parseLines(tiv)
}

func parseExcitabilityVariables(reader *Reader, exciteVar *ExcitabilityVariables) error {
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

func parseStrengthDuration(reader *Reader, section *StrengthDuration) error {
	err := reader.skipPast(`%CMAP              	Threshold`)
	if err != nil {
		return err
	}

	return reader.parseLines(section)
}
