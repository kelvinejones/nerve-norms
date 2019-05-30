package mem

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
)

type StimResponse struct {
	MaxCmaps
	ValueType  string
	PercentMax Column
	Stimulus   Column
	WasImputed Column
}

var SRPercentMax = Column([]float64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98})

func (sr *StimResponse) LoadFromMem(mem *rawMem) error {
	sr.PercentMax = SRPercentMax

	sec, err := mem.sectionContainingHeader("STIMULUS RESPONSE")
	if err != nil {
		return errors.New("Could not get stimulus response: " + err.Error())
	}

	perMax, err := sec.columnContainsName("% Max", 0)
	if err != nil {
		return errors.New("Could not get stimulus response: " + err.Error())
	}

	sr.Stimulus, err = sec.columnContainsName("Stimulus", 0)
	if err != nil {
		return errors.New("Could not get stimulus response: " + err.Error())
	}

	sr.WasImputed = sr.Stimulus.ImputeWithValue(perMax, sr.PercentMax, 0.1, false)

	sr.ValueType = parseValueType(sec.ExtraLines)
	sr.MaxCmaps = parseMaxCmap(sec.ExtraLines)

	return nil
}

type jsonStimResponse struct {
	Columns   []string `json:"columns"`
	Data      Table    `json:"data"`
	MaxCmaps  `json:"maxCmaps"`
	ValueType string `json:"valueType"`
}

func (dat *StimResponse) MarshalJSON() ([]byte, error) {
	str := &jsonStimResponse{
		Columns:   []string{"% Max", "Stimulus"},
		Data:      []Column{dat.PercentMax, dat.Stimulus},
		MaxCmaps:  dat.MaxCmaps,
		ValueType: dat.ValueType,
	}

	if dat.WasImputed != nil {
		str.Columns = append(str.Columns, "Was Imputed")
		str.Data = append(str.Data, dat.WasImputed)
	}

	return json.Marshal(&str)
}

func (dat *StimResponse) UnmarshalJSON(value []byte) error {
	jsDat := jsonStimResponse{}
	err := json.Unmarshal(value, &jsDat)
	if err != nil {
		return err
	}
	numCol := len(jsDat.Columns)

	if numCol < 2 || numCol > 3 {
		return errors.New("Incorrect number of StimResponse columns in JSON")
	}
	if jsDat.Columns[0] != "% Max" || jsDat.Columns[1] != "Stimulus" || (numCol == 3 && jsDat.Columns[2] != "Was Imputed") {
		return errors.New("Incorrect StimResponse column names in JSON")
	}

	dat.PercentMax = jsDat.Data[0]
	dat.Stimulus = jsDat.Data[1]
	if numCol == 3 {
		dat.WasImputed = jsDat.Data[2]
	}

	dat.MaxCmaps = jsDat.MaxCmaps
	dat.ValueType = jsDat.ValueType

	return nil
}

func parseValueType(strs []string) string {
	reg := regexp.MustCompile(`^Values (.*)`)
	for _, str := range strs {
		result := reg.FindStringSubmatch(str)

		if len(result) != 2 {
			// The string couldn't be parsed. This isn't an error;
			// it just means this line wasn't a value type.
			continue
		}

		return result[1]
	}

	return ""
}

type MaxCmap struct {
	Time  float64 `json:"time"`
	Val   float64 `json:"value"`
	Units byte    `json:"units"`
}

type MaxCmaps []MaxCmap

func (mc MaxCmaps) AsFloat() (float64, error) {
	for _, val := range mc {
		if val.Time > 0.99 && val.Time < 1.01 {
			scale := 1.0
			switch val.Units {
			case 'u':
				scale = 1.0 / 1000
			case 'm':
				// Do nothing; default is mV
			default:
				return 0.0, errors.New("Could not parse CMAP unit '" + string(val.Units) + "'")
			}
			return val.Val * scale, nil
		}
	}
	return 0.0, errors.New("Could not find CMAP at 1ms")
}

func parseMaxCmap(strs []string) []MaxCmap {
	cmaps := []MaxCmap{}
	reg := regexp.MustCompile(`Max CMAP  (\d*\.?\d+) ms =  (\d*\.?\d+) (.)V`)

	for _, str := range strs {
		cmap := MaxCmap{}

		result := reg.FindStringSubmatch(str)

		if len(result) != 4 {
			// The string couldn't be parsed. This isn't an error;
			// it just means this line wasn't a MaxCmap.
			continue
		}

		var err error
		cmap.Time, err = strconv.ParseFloat(result[1], 64)
		if err != nil {
			// The string couldn't be parsed. This isn't an error;
			// it just means this line wasn't a MaxCmap.
			continue
		}
		cmap.Val, err = strconv.ParseFloat(result[2], 64)
		if err != nil {
			// The string couldn't be parsed. This isn't an error;
			// it just means this line wasn't a MaxCmap.
			continue
		}
		cmap.Units = result[3][0]

		cmaps = append(cmaps, cmap)
	}

	return cmaps
}
