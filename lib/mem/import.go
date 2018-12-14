package mem

import (
	"errors"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	return reader.parseLines(headerRegex, header)
}

var headerRegex = regexp.MustCompile(`^\s+([^:]+):\s+(.*)`)

func (header *Header) Parse(result []string) error {
	if len(result) != 3 {
		return errors.New("Incorrect header line length")
	}

	var err error
	val := strings.TrimSpace(result[2])
	switch result[1] {
	case "NC/disease":
		if val == "NC" {
			header.NormalControl = true
		}
		// TODO update for other options? Currently disease is the default, which I suppose excludes an uncertain ones from the control database
	case "Sex":
		switch val {
		case "M":
			header.Sex = MaleSex
		case "F":
			header.Sex = FemaleSex
		default:
			header.Sex = UnknownSex
		}
	case "Temperature":
		temp, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		header.Temperature = temp
	case "Age":
		header.Age, err = strconv.Atoi(val)
		if err != nil {
			return err
		}
	case "Date":
		layout := "2/1/06"
		header.Date, err = time.Parse(layout, val)
		if err != nil {
			return err
		}
	case "Start time":
		layout := "2/1/06 15:04:05"
		header.StartTime, err = time.Parse(layout, "2/1/06 "+val)
		if err != nil {
			return err
		}
	case "File":
		header.File = val
	case "Name":
		header.Name = val
	case "Protocol":
		header.Protocol = val
	case "S/R sites":
		header.SRSites = val
	case "Operator":
		header.Operator = val
	case "Comments":
		header.Comment = val
	}

	return nil
}

func parseStimResponse(reader *Reader, sr *StimResponse) error {
	var err error
	sr.ValueType, err = reader.ReadLineExtractingString(`^Values (.*)`)
	if err != nil {
		return err
	}

	// Find Max CMAP
	err = reader.parseLines(maxCmapRegex, &sr.MaxCmaps)
	if err != nil {
		return err
	}

	err = reader.skipPast("% Max               	Stimulus")
	if err != nil {
		return err
	}

	// Now parse the actual SR data
	return reader.parseLines(srRegex, sr)
}

var srRegex = regexp.MustCompile(`^SR\.(\d+)\s+(\d+)\s+(\d*\.?\d+)`)

func (sr *StimResponse) Parse(result []string) error {
	if len(result) != 4 {
		return errors.New("Incorrect SR line length")
	}
	if result[1] != result[2] {
		return errors.New("SR fields do not match")
	}

	percentMax, err := strconv.ParseFloat(result[1], 64)
	if err != nil {
		return err
	}
	stim, err := strconv.ParseFloat(result[3], 64)
	if err != nil {
		return err
	}

	sr.Values = append(sr.Values, XY{
		X: percentMax,
		Y: stim,
	})

	return nil
}

var maxCmapRegex = regexp.MustCompile(`^ Max CMAP  (\d*\.?\d+) ms =  (\d*\.?\d+) (.)V`)

func (cmaps *MaxCmaps) Parse(result []string) error {
	if len(result) != 4 {
		return errors.New("Incorrect CMAP line length")
	}

	cmap := MaxCmap{}
	var err error
	cmap.Time, err = strconv.ParseFloat(result[1], 64)
	if err != nil {
		return err
	}
	cmap.Val, err = strconv.ParseFloat(result[2], 64)
	if err != nil {
		return err
	}
	cmap.Units = result[3][0]

	*cmaps = append(*cmaps, cmap)

	return nil
}

func parseChargeDuration(reader *Reader, cd *ChargeDuration) error {
	err := reader.skipPast("Duration (ms)       	 Threshold (mA)     	  Threshold charge (mA.mS)")
	if err != nil {
		return err
	}

	return reader.parseLines(chargeRegex, cd)
}

var chargeRegex = regexp.MustCompile(`^QT\.\d+\s+(\d*\.?\d+)\s+(\d*\.?\d+)\s+(\d*\.?\d+)`)

func (cd *ChargeDuration) Parse(result []string) error {
	if len(result) != 4 {
		return errors.New("Incorrect CD line length")
	}

	x, err := strconv.ParseFloat(result[1], 64)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(result[2], 64)
	if err != nil {
		return err
	}
	z, err := strconv.ParseFloat(result[3], 64)
	if err != nil {
		return err
	}

	cd.Values = append(cd.Values, XYZ{
		X: x,
		Y: y,
		Z: z,
	})

	return nil
}

func parseThresholdElectrotonus(reader *Reader, te *ThresholdElectrotonusGroup) error {
	err := reader.skipPast("Delay (ms)          	Current (%)         	Thresh redn. (%)")
	if err != nil {
		return err
	}

	return reader.parseLines(teRegex, te)
}

var teRegex = regexp.MustCompile(`^TE(\d+)\.\d+\s+(\d*\.?\d+)\s+([-+]?\d*\.?\d+)\s+([-+]?\d*\.?\d+)`)

func (te *ThresholdElectrotonusGroup) Parse(result []string) error {
	if len(result) != 5 {
		return errors.New("Incorrect TE line length")
	}

	set, err := strconv.Atoi(result[1])
	if err != nil {
		return err
	}
	if set > 100 {
		// Assume this is a parse error
		return errors.New("More than 100 sets of TE data are not supported")
	}

	x, err := strconv.ParseFloat(result[2], 64)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(result[3], 64)
	if err != nil {
		return err
	}
	z, err := strconv.ParseFloat(result[4], 64)
	if err != nil {
		return err
	}

	for len(te.Sets) < set {
		// This would be inefficient for a big difference, but usually this will only run once
		te.Sets = append(te.Sets, ThresholdElectrotonusSet{})
	}

	te.Sets[set-1].Values = append(te.Sets[set-1].Values, XYZ{
		X: x,
		Y: y,
		Z: z,
	})

	return nil
}

func parseRecoveryCycle(reader *Reader, rc *RecoveryCycle) error {
	err := reader.skipPast("Interval (ms)       	  Threshold change (%)")
	if err != nil {
		return err
	}

	return reader.parseLines(rcRegex, rc)
}

var rcRegex = regexp.MustCompile(`^RC\d+\.\d+\s+(\d*\.?\d+)\s+([-+]?\d*\.?\d+)`)

func (rc *RecoveryCycle) Parse(result []string) error {
	if len(result) != 3 {
		return errors.New("Incorrect RC line length")
	}

	x, err := strconv.ParseFloat(result[1], 64)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(result[2], 64)
	if err != nil {
		return err
	}

	rc.Values = append(rc.Values, XY{
		X: x,
		Y: y,
	})

	return nil
}

func parseThresholdIV(reader *Reader, tiv *ThresholdIV) error {
	err := reader.skipPast("Current (%)         	  Threshold redn. (%)")
	if err != nil {
		return err
	}

	return reader.parseLines(tivRegex, tiv)
}

var tivRegex = regexp.MustCompile(`^IV\d+\.\d+\s+([-+]?\d*\.?\d+)\s+([-+]?\d*\.?\d+)`)

func (tiv *ThresholdIV) Parse(result []string) error {
	if len(result) != 3 {
		return errors.New("Incorrect TIV line length")
	}

	x, err := strconv.ParseFloat(result[1], 64)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(result[2], 64)
	if err != nil {
		return err
	}

	tiv.Values = append(tiv.Values, XY{
		X: x,
		Y: y,
	})

	return nil
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
	err = reader.parseLines(exciteVarRegex, exciteVar)
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

	err = reader.parseLines(extraVarRegex, &ExtraVariables{exciteVar})
	return err
}

var exciteVarRegex = regexp.MustCompile(`^ \d+\.\s+([-+]?\d*\.?\d+)\s+(.+)`)

func (exciteVar *ExcitabilityVariables) Parse(result []string) error {
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

var extraVarRegex = regexp.MustCompile(`^(.+) = ([-+]?\d*\.?\d+)`)

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

func parseStrengthDuration(reader *Reader, section *StrengthDuration) error {
	err := reader.skipPast(`%CMAP              	Threshold`)
	if err != nil {
		return err
	}

	return reader.parseLines(strengthDurationRegex, section)
}

var strengthDurationRegex = regexp.MustCompile(`^SD\.\d+.*`)

func (section *StrengthDuration) Parse(result []string) error {
	return nil
}
