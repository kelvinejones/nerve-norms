package mem

import (
	"encoding/json"
	"regexp"
	"strconv"
)

type StimResponse struct {
	MC        MaxCmaps `json:"maxCmaps"`
	ValueType string   `json:"valueType"`
	LT        LabTab   `json:"data"`
}

var SRPercentMax = Column([]float64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58, 60, 62, 64, 66, 68, 70, 72, 74, 76, 78, 80, 82, 84, 86, 88, 90, 92, 94, 96, 98})

func newSR() *StimResponse {
	sr := &StimResponse{
		LT: LabTab{
			section:   "STIMULUS RESPONSE",
			xname:     "% Max",
			yname:     "Stimulus",
			xcol:      SRPercentMax,
			precision: 0.1,
		},
	}
	sr.LT.extraImport = func(sec RawSection) {
		sr.ValueType = parseValueType(sec.ExtraLines)
		sr.MC.parseMaxCmap(sec.ExtraLines)
	}
	return sr
}

func (sr *StimResponse) LoadFromMem(mem *rawMem) error {
	return sr.LT.LoadFromMem(mem)
}

func (sr *StimResponse) LabelledTable(subsec string) LabelledTable {
	switch subsec {
	case "CMAP":
		return sr.MC.AsLabelledTable()
	case "":
		return sr.LT
	default:
		return nil
	}
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

type MaxCmaps struct {
	// all contains all CMAP measurements
	all []MaxCmap

	// standard is the max CMAP at 1ms
	standard float64
}

func (mc MaxCmaps) MarshalJSON() ([]byte, error) {
	return json.Marshal(&mc.all)
}

func (mc *MaxCmaps) UnmarshalJSON(value []byte) error {
	err := json.Unmarshal(value, &mc.all)
	mc.parseStandardCmap()
	return err
}

func (mc MaxCmaps) AsLabelledTable() LabelledTable {
	lt := LabTab{
		xname: "Time (ms)",
		yname: "CMAP",
		xcol:  Column{1},
		ycol:  Column{mc.standard},
	}
	if mc.standard == 0.0 {
		lt.wasimp = Column{1}
	}
	return &lt
}

func (mc *MaxCmaps) parseMaxCmap(strs []string) {
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

		mc.all = append(mc.all, cmap)
	}

	mc.parseStandardCmap()
}

func (mc *MaxCmaps) parseStandardCmap() {
	for _, val := range mc.all {
		if val.Time > 0.99 && val.Time < 1.01 {
			switch val.Units {
			case 'u':
				mc.standard = val.Val / 1000
			case 'm':
				// Do nothing; default is mV
				mc.standard = val.Val
			}
		}
	}
}
