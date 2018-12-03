package mem

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
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

	return mem, nil
}

func skipNewlines(reader *Reader) error {
	s := "\n"
	var err error
	for s == "\n" && err == nil {
		s, err = reader.ReadString('\n')
	}
	reader.UnreadString(s)

	return err
}

func skipPast(reader *Reader, search string) error {
	err := skipNewlines(reader)
	if err != nil {
		return err
	}

	s, err := reader.ReadString('\n')
	if err == nil && !strings.Contains(s, search) {
		err = errors.New("Could not find '" + search + "'")
	}
	return err
}

func parseHeader(reader *Reader, header *Header) error {
	return parseLines(reader, headerRegex, header)
}

func parseStimResponse(reader *Reader, sr *StimResponse) error {
	// Find section header
	err := skipPast(reader, "STIMULUS-RESPONSE DATA")
	if err != nil {
		return err
	}

	// Find some random string that's there
	err = skipPast(reader, "Values are those recorded")
	if err != nil {
		return err
	}

	// Find Max CMAP
	err = skipNewlines(reader)
	if err != nil {
		return err
	}
	s, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	n, err := fmt.Sscanf(s, " Max CMAP  1 ms =  %f mV\n", &sr.MaxCmap)
	if n != 1 || err != nil {
		return errors.New("Could not find Max CMAP: " + s)
	}

	err = skipPast(reader, "% Max               	Stimulus")
	if err != nil {
		return err
	}

	// Now parse the actual SR data
	return parseLines(reader, srRegex, sr)
}

var srRegex = regexp.MustCompile(`^SR\.(\d+)\s+(\d+)\s+(\d*\.?\d+)`)

func (sr *StimResponse) Parse(result []string) error {
	if len(result) != 4 {
		return errors.New("Incorrect SR line length")
	}
	if result[1] != result[2] {
		return errors.New("SR fields do not match")
	}

	percentMax, err := strconv.Atoi(result[1])
	if err != nil {
		return err
	}
	stim, err := strconv.ParseFloat(result[3], 32)
	if err != nil {
		return err
	}

	sr.Values = append(sr.Values, XY{
		X: percentMax,
		Y: float32(stim),
	})

	return nil
}

func parseChargeDuration(reader *Reader, cd *ChargeDuration) error {
	// Find section header
	err := skipPast(reader, "CHARGE DURATION DATA")
	if err != nil {
		return err
	}

	err = skipPast(reader, "Duration (ms)       	 Threshold (mA)     	  Threshold charge (mA.mS)")
	if err != nil {
		return err
	}

	return parseLines(reader, chargeRegex, cd)
}

var chargeRegex = regexp.MustCompile(`^QT\.\d+\s+(\d*\.?\d+)\s+(\d*\.?\d+)\s+(\d*\.?\d+)`)

func (cd *ChargeDuration) Parse(result []string) error {
	if len(result) != 4 {
		return errors.New("Incorrect CD line length")
	}

	x, err := strconv.ParseFloat(result[1], 32)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(result[2], 32)
	if err != nil {
		return err
	}
	z, err := strconv.ParseFloat(result[3], 32)
	if err != nil {
		return err
	}

	cd.Values = append(cd.Values, XYZ{
		X: float32(x),
		Y: float32(y),
		Z: float32(z),
	})

	return nil
}

func parseThresholdElectrotonus(reader *Reader, te *ThresholdElectrotonusGroup) error {
	// Find section header
	err := skipPast(reader, "THRESHOLD ELECTROTONUS DATA")
	if err != nil {
		return err
	}

	err = skipPast(reader, "Delay (ms)          	Current (%)         	Thresh redn. (%)")
	if err != nil {
		return err
	}

	return parseLines(reader, teRegex, te)
}

var teRegex = regexp.MustCompile(`^TE(\d+)\.\d+\s+(\d*\.?\d+)\s+(\d*\.?\d+)\s+(\d*\.?\d+)`)

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

	x, err := strconv.ParseFloat(result[2], 32)
	if err != nil {
		return err
	}
	y, err := strconv.ParseFloat(result[3], 32)
	if err != nil {
		return err
	}
	z, err := strconv.ParseFloat(result[4], 32)
	if err != nil {
		return err
	}

	for len(te.Sets) < set {
		// This would be inefficient for a big difference, but usually this will only run once
		te.Sets = append(te.Sets, ThresholdElectrotonusSet{})
	}

	te.Sets[set-1].Values = append(te.Sets[set-1].Values, XYZ{
		X: float32(x),
		Y: float32(y),
		Z: float32(z),
	})

	return nil
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
		temp, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return err
		}
		header.Temperature = float32(temp)
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

type Parser interface {
	Parse([]string) error
}

func parseLines(reader *Reader, regex *regexp.Regexp, parser Parser) error {
	err := skipNewlines(reader)
	if err != nil {
		return err
	}

	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		if len(s) == 1 {
			// Done with section; break!
			break
		}
		result := regex.FindStringSubmatch(s)

		err = parser.Parse(result)
		if err != nil {
			return errors.New(err.Error() + ": '" + s + "'")
		}
	}

	return nil
}
