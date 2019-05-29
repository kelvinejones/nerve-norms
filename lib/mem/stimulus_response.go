package mem

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
)

type StimResponse struct {
	MaxCmaps   []MaxCmap
	ValueType  string
	PercentMax Column
	Stimulus   Column
	WasImputed Column
}

func (mem *Mem) StimulusResponse() (StimResponse, error) {
	sr := StimResponse{PercentMax: Column([]float64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98})}

	sec, err := mem.sectionContainingHeader("STIMULUS RESPONSE")
	if err != nil {
		return sr, errors.New("Could not get stimulus response: " + err.Error())
	}

	perMax, err := sec.columnContainsName("% Max", 0)
	if err != nil {
		return sr, errors.New("Could not get stimulus response: " + err.Error())
	}

	sr.Stimulus, err = sec.columnContainsName("Stimulus", 0)
	if err != nil {
		return sr, errors.New("Could not get stimulus response: " + err.Error())
	}

	sr.WasImputed = sr.Stimulus.ImputeWithValue(perMax, sr.PercentMax, 0.1)

	sr.ValueType = parseValueType(sec.ExtraLines)
	sr.MaxCmaps = parseMaxCmap(sec.ExtraLines)

	return sr, nil
}

func (dat *StimResponse) MarshalJSON() ([]byte, error) {
	str := &struct {
		Columns   []string  `json:"columns"`
		Data      Table     `json:"data"`
		MaxCmaps  []MaxCmap `json:"maxCmaps"`
		ValueType string    `json:"valueType"`
	}{
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
