package mem

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"
)

func Import(data io.Reader) (Mem, error) {
	reader := NewReader(data)
	mem := Mem{}
	var err error

	err = parseHeader(reader, &mem.Header)
	if err != nil {
		return mem, err
	}

	err = parseStimResponse(reader, &mem.StimResponse)
	if err != nil {
		return mem, err
	}

	err = parseChargeDuration(reader, &mem.ChargeDuration)
	if err != nil {
		return mem, err
	}

	err = parseThresholdElectrotonus(reader, &mem.ThresholdElectrotonusGroup)
	if err != nil {
		return mem, err
	}

	err = parseRecoveryCycle(reader, &mem.RecoveryCycle)
	if err != nil {
		return mem, err
	}

	err = parseThresholdIV(reader, &mem.ThresholdIV)
	if err != nil {
		return mem, err
	}

	mem.ExcitabilityVariables.Values = make(map[string]float64)
	err = parseExcitabilityVariables(reader, &mem.ExcitabilityVariables)
	if err != nil {
		return mem, err
	}

	return mem, nil
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
	val := result[2]
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
	// Find section header
	err := reader.skipPast("STIMULUS-RESPONSE DATA")
	if err != nil {
		return err
	}

	// Find some random string that's there
	err = reader.skipPast("Values are those recorded")
	if err != nil {
		return err
	}

	// Find Max CMAP
	err = reader.skipNewlines()
	if err != nil {
		return err
	}
	s, err := reader.ReadLine()
	if err != nil {
		return err
	}
	n, err := fmt.Sscanf(s, " Max CMAP  1 ms =  %f mV", &sr.MaxCmap)
	if n != 1 || err != nil {
		return errors.New("Could not find Max CMAP: " + s)
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

func parseChargeDuration(reader *Reader, cd *ChargeDuration) error {
	// Find section header
	err := reader.skipPast("CHARGE DURATION DATA")
	if err != nil {
		return err
	}

	err = reader.skipPast("Duration (ms)       	 Threshold (mA)     	  Threshold charge (mA.mS)")
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
	// Find section header
	err := reader.skipPast("THRESHOLD ELECTROTONUS DATA")
	if err != nil {
		return err
	}

	err = reader.skipPast("Delay (ms)          	Current (%)         	Thresh redn. (%)")
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
	// Find section header
	err := reader.skipPast("RECOVERY CYCLE DATA")
	if err != nil {
		return err
	}

	err = reader.skipPast("Interval (ms)       	  Threshold change (%)")
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
	// Find section header
	err := reader.skipPast("THRESHOLD I/V DATA")
	if err != nil {
		return err
	}

	err = reader.skipPast("Current (%)         	  Threshold redn. (%)")
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
	// Find section header
	err := reader.skipPast("DERIVED EXCITABILITY VARIABLES")
	if err != nil {
		return err
	}

	// Find settings
	err = reader.skipNewlines()
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
	err = reader.skipPast("EXTRA VARIABLES")
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

type ExtraVariables struct {
	*ExcitabilityVariables
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
