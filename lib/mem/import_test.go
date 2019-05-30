package mem

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// For readability, all strings in this file are encoded with \n, but the code requires \r\n
func toWindows(str string) string {
	return strings.Replace(str, "\n", "\r\n", -1)
}

const headerString = ` File:              	n:\Short Test\FESB70821A.QZD
 Name:              	SHORTY
 Protocol:          	TRONDNF
 Date:              	21/8/17
 Start time:        	12:57:17
 Age:               	30
 Sex:               	M
 Temperature:       	33.5
 S/R sites:         	Median Wr-APB
 NC/disease:        	NC
 Operator:          	MS
 Comments:          	smooth recording

`

const dualHeaderString = ` File:              	n:\Short Test\FESB70821A.QZD
 Name:              	SHORTY
 Protocol:          	TRONDNF
 Date:              	21/8/17 / 28/9/18
 Start time:        	12:57:17 / 14:58:00
 Age:               	30 / 31
 Sex:               	M
 Temperature:       	33.5 / 33.6
 S/R sites:         	Median Wr-APB
 NC/disease:        	NC
 Operator:          	MS
 Comments:          	smooth recording

`

var headerExpected = Header{
	File:          `n:\Short Test\FESB70821A.QZD`,
	Name:          "SHORTY",
	Protocol:      "TRONDNF",
	Date:          time.Date(2017, time.Month(8), 21, 0, 0, 0, 0, time.UTC),   // TODO eventually handle time zones better?
	StartTime:     time.Date(2006, time.Month(1), 2, 12, 57, 17, 0, time.UTC), // TODO eventually handle time zones better?
	Age:           30,
	Sex:           MaleSex,
	Temperature:   33.5,
	SRSites:       "Median Wr-APB",
	NormalControl: true,
	Operator:      "MS",
	Comment:       "smooth recording",
}

func TestImportEmpty(t *testing.T) {
	t.Skip()
	m, err := Import(strings.NewReader(""))
	assert.NoError(t, err)
	assert.Equal(t, m, rawMem{})
}

func TestImportHeader(t *testing.T) {
	header := Header{}
	err := header.Parse(NewStringReader(toWindows(headerString)))
	assert.NoError(t, err)
	assert.Equal(t, headerExpected, header)
}

func TestImportDualHeader(t *testing.T) {
	header := Header{}
	err := header.Parse(NewStringReader(toWindows(dualHeaderString)))
	assert.NoError(t, err)
	assert.Equal(t, headerExpected, header)
}

const sResponseHeaderString = `

 STIMULUS-RESPONSE DATA (2.4-1.9m)
`
const sResponseString = `
Values are those recorded

 Max CMAP  .2 ms =  61.36306 uV
 Max CMAP  1 ms =  1.161296 mV

                    	% Max               	Stimulus(2)
SR.2                	 2                  	 3.915578
SR.4                	 4                  	 4.073214
SR.6                	 6                  	 4.144141
SR.8                	 8                  	 4.20404
SR.10               	 10                 	 4.435846
SR.12               	 12                 	 4.601757
SR.14               	 14                 	 4.824213
SR.16               	 16                 	 4.86682
SR.18               	 18                 	 4.89536
SR.20               	 20                 	 4.9239

`

var sResponsePercentMaxColumn = Column{2., 4., 6., 8., 10., 12., 14., 16., 18., 20.}
var sResponseStimulusColumn = Column{3.915578, 4.073214, 4.144141, 4.20404, 4.435846, 4.601757, 4.824213, 4.86682, 4.89536, 4.9239}
var sResponseExpected = RawSection{
	Header: "STIMULUS-RESPONSE DATA (2.4-1.9m)",
	TableSet: TableSet{
		ColCount: 2,
		Names:    []string{"% Max", "Stimulus(2)"},
		Tables: []Table{Table{
			sResponsePercentMaxColumn,
			sResponseStimulusColumn,
		}},
	},
	ExtraLines: []string{
		"Values are those recorded",
		"Max CMAP  .2 ms =  61.36306 uV",
		"Max CMAP  1 ms =  1.161296 mV",
	},
}
var sResponseParsed = StimResponse{
	MaxCmaps: []MaxCmap{
		MaxCmap{
			Time:  .2,
			Val:   61.36306,
			Units: 'u',
		},
		MaxCmap{
			Time:  1.,
			Val:   1.161296,
			Units: 'm',
		},
	},
	ValueType:  "are those recorded",
	PercentMax: sResponsePercentMaxColumn,
	Stimulus:   sResponseStimulusColumn,
}

func TestImportSRResponse(t *testing.T) {
	sec := RawSection{Header: sResponseExpected.Header}
	err := sec.parse(NewStringReader(toWindows(sResponseString)))
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, sResponseExpected, sec)
}

const chargeDurationHeaderString = `

  CHARGE DURATION DATA (2.4-3.5m)
`
const chargeDurationString = `
                    	Duration (ms)       	 Threshold (mA)     	  Threshold charge (mA.mS)
QT.1                	 .2                 	 9.790961           	 1.958192
QT.2                	 .4                 	 6.905862           	 2.762345
QT.3                	 .6                 	 5.978864           	 3.587318
QT.4                	 .8                 	 5.44341            	 4.354728
QT.5                	 1                  	 5.187509           	 5.187509

`

var chargeDurationDurationColumn = Column{.2, .4, .6, .8, 1.}
var chargeDurationThreshChargeColumn = Column{1.958192, 2.762345, 3.587318, 4.354728, 5.187509}
var chargeDurationExpected = RawSection{
	Header: "CHARGE DURATION DATA (2.4-3.5m)",
	TableSet: TableSet{
		ColCount: 3,
		Names:    []string{"Duration (ms)", "Threshold (mA)", "Threshold charge (mA.mS)"},
		Tables: []Table{Table{
			chargeDurationDurationColumn,
			Column{9.790961, 6.905862, 5.978864, 5.44341, 5.187509},
			chargeDurationThreshChargeColumn,
		}},
	},
}
var chargeDurationParsed = ChargeDuration{
	Duration:     chargeDurationDurationColumn,
	ThreshCharge: chargeDurationThreshChargeColumn,
}

func TestImportChargeDuration(t *testing.T) {
	sec := RawSection{Header: chargeDurationExpected.Header}
	err := sec.parse(NewStringReader(toWindows(chargeDurationString)))
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, chargeDurationExpected, sec)
}

const thresholdElectrotonusHeaderString = `

  THRESHOLD ELECTROTONUS DATA (3.5-8.8m)
`
const thresholdElectrotonusString = `
                    	Delay (ms)          	Current (%)         	Thresh redn. (%)
TE1.1               	 0                  	0                   	0.00
TE1.2               	 9                  	0                   	0.00
TE1.3               	 10                 	40                  	40.02
TE1.4               	 11                 	40                  	42.71
TE1.5               	 11                 	40                  	-42.71

TE2.1               	 0                  	0                   	0.00
TE2.2               	 9                  	0                   	0.00
TE2.3               	 10                 	-40                 	-39.34
TE2.4               	 11                 	-40                 	-40.92

`

var thresholdElectrotonusHyp40DelayColumn = Column{0., 9., 10., 11., 11.}
var thresholdElectrotonusHyp40ThreshReductionColumn = Column{0.00, 0.00, 40.02, 42.71, -42.71}
var thresholdElectrotonusDep40DelayColumn = Column{0., 9., 10., 11.}
var thresholdElectrotonusDep40ThreshReductionColumn = Column{0.00, 0.00, -39.34, -40.92}
var thresholdElectrotonusExpected = RawSection{
	Header: "THRESHOLD ELECTROTONUS DATA (3.5-8.8m)",
	TableSet: TableSet{
		ColCount: 3,
		Names:    []string{"Delay (ms)", "Current (%)", "Thresh redn. (%)"},
		Tables: []Table{
			Table{
				thresholdElectrotonusHyp40DelayColumn,
				Column{0., 0., 40., 40., 40.},
				thresholdElectrotonusHyp40ThreshReductionColumn,
			},
			Table{
				thresholdElectrotonusDep40DelayColumn,
				Column{0., 0., -40., -40.},
				thresholdElectrotonusDep40ThreshReductionColumn,
			},
		},
	},
}
var thresholdElectrotonusParsed = ThresholdElectrotonus{
	Hyperpol40: TEPair{
		Delay:           thresholdElectrotonusHyp40DelayColumn,
		ThreshReduction: thresholdElectrotonusHyp40ThreshReductionColumn,
	},
	Depol40: TEPair{
		Delay:           thresholdElectrotonusDep40DelayColumn,
		ThreshReduction: thresholdElectrotonusDep40ThreshReductionColumn,
	},
}

func TestImportThresholdElectrotonus(t *testing.T) {
	sec := RawSection{Header: thresholdElectrotonusExpected.Header}
	err := sec.parse(NewStringReader(toWindows(thresholdElectrotonusString)))
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, thresholdElectrotonusExpected, sec)
}

const recoveryCycleHeaderString = `

  RECOVERY CYCLE DATA (11.1-15.3m)
`
const recoveryCycleString = `
                    	Interval (ms)       	  Threshold change (%)
RC1.1               	 3.2                	4.99
RC1.2               	 4                  	-12.75
RC1.3               	 5                  	-22.24
RC1.4               	 6.3                	-24.45
RC1.5               	 7.9                	-24.05

`

var recoveryCycleIntervalColumn = Column{3.2, 4., 5., 6.3, 7.9}
var recoveryCycleThreshChangeColumn = Column{4.99, -12.75, -22.24, -24.45, -24.05}
var recoveryCycleExpected = RawSection{
	Header: "RECOVERY CYCLE DATA (11.1-15.3m)",
	TableSet: TableSet{
		ColCount: 2,
		Names:    []string{"Interval (ms)", "Threshold change (%)"},
		Tables: []Table{Table{
			recoveryCycleIntervalColumn,
			recoveryCycleThreshChangeColumn,
		}},
	},
}
var recoveryCycleParsed = RecoveryCycle{
	Interval:     recoveryCycleIntervalColumn,
	ThreshChange: recoveryCycleThreshChangeColumn,
}

func TestImportRecoveryCycle(t *testing.T) {
	sec := RawSection{Header: recoveryCycleExpected.Header}
	err := sec.parse(NewStringReader(toWindows(recoveryCycleString)))
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, recoveryCycleExpected, sec)
}

const thresholdIVHeaderString = `

  THRESHOLD I/V DATA (8.9-11m)
`
const thresholdIVString = `
                    	Current (%)         	  Threshold redn. (%)
IV1.1               	 50                 	49.28
IV1.2               	 40                 	39.01
IV1.3               	 30                 	31.59
IV1.4               	 20                 	22.58
IV1.5               	 10                 	13.06
IV1.6               	 0                  	-0.78
IV1.7               	-10                 	-17.58
IV1.8               	-20                 	-39.31

`

var thresholdIVCurrentColumn = Column{50., 40., 30., 20., 10., 0., -10., -20.}
var thresholdIVThreshReductionColumn = Column{49.28, 39.01, 31.59, 22.58, 13.06, -0.78, -17.58, -39.31}
var thresholdIVExpected = RawSection{
	Header: "THRESHOLD I/V DATA (8.9-11m)",
	TableSet: TableSet{
		ColCount: 2,
		Names:    []string{"Current (%)", "Threshold redn. (%)"},
		Tables: []Table{Table{
			thresholdIVCurrentColumn,
			thresholdIVThreshReductionColumn,
		}},
	},
}
var thresholdIVParsed = ThresholdIV{
	Current:         thresholdIVCurrentColumn,
	ThreshReduction: thresholdIVThreshReductionColumn,
}

func TestImportThresholdIV(t *testing.T) {
	sec := RawSection{Header: thresholdIVExpected.Header}
	err := sec.parse(NewStringReader(toWindows(thresholdIVString)))
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, thresholdIVExpected, sec)
}

const excitabilityVariablesHeaderString = `

  DERIVED EXCITABILITY VARIABLES
`
const excitabilityVariablesString = `
Program = QTracP 9/12/2016
Threshold method = 6 (optimised for CAP, using data for present condition only)
SR method = 1 (using actual data values)

 1.                 	5.307               	Stimulus (mA) for 50% max response
 3.                 	0.287               	Strength-duration\time constant (ms)
 16.                	0.000               	Polarizing current\(mA)
 18.                	1                   	Sex (M=1, F=2)
 22.                	-64.192             	TEh(10-20ms)

  EXTRA VARIABLES (add here as required, e.g. Potassium = 4.5)

TEd40(Accom) = 19.6
TEd20(10-20ms) = 30.8
TEh20(10-20ms) = -32.2
TESTINGxTABx                  	=  23.7

`

var excitabilityVariablesExpected = ExcitabilityVariables{
	Values: map[int]float64{
		1:    5.307,   // 		`Stimulus (mA) for 50% max response`
		3:    0.287,   // 		`Strength-duration\time constant (ms)`
		16:   0.000,   // 		`Polarizing current\(mA)`
		18:   1,       // 		`Sex (M=1, F=2)`
		22:   -64.192, // 		`TEh(10-20ms)`
		1001: 19.6,    // 		`TEd40(Accom)`
		1002: 30.8,    // 		`TEd20(10-20ms)`
		1003: -32.2,   // 		`TEh20(10-20ms)`
		1020: 23.7,    // 		`MRCsumscore`
	},
	ExcitabilitySettings: map[string]string{
		"Program":          "QTracP 9/12/2016",
		"Threshold method": "6 (optimised for CAP, using data for present condition only)",
		"SR method":        "1 (using actual data values)",
	},
}

func TestImportExcitabilityVariables(t *testing.T) {
	actual := ExcitabilityVariables{Values: make(map[int]float64)}
	err := actual.Parse(NewStringReader(toWindows(excitabilityVariablesString)))
	assert.NoError(t, err)
	assert.Equal(t, excitabilityVariablesExpected, actual)
}

var completeExpectedRawMem = rawMem{
	Header: headerExpected,
	Sections: []RawSection{
		sResponseExpected,
		chargeDurationExpected,
		thresholdElectrotonusExpected,
		recoveryCycleExpected,
		thresholdIVExpected,
	},
	ExcitabilityVariables: excitabilityVariablesExpected,
}

var memString = headerString + sResponseHeaderString + sResponseString + chargeDurationHeaderString + chargeDurationString +
	thresholdElectrotonusHeaderString + thresholdElectrotonusString + recoveryCycleHeaderString + recoveryCycleString +
	thresholdIVHeaderString + thresholdIVString + excitabilityVariablesHeaderString + excitabilityVariablesString

func TestImportAll(t *testing.T) {
	// This common setup should fail in one place for all of these tests.
	mem, err := Import(strings.NewReader(toWindows(memString)))
	assert.NoError(t, err)

	t.Run("rawMemStruct", func(t *testing.T) {
		assert.Equal(t, completeExpectedRawMem, mem)
	})
	t.Run("StimulusResponse", func(t *testing.T) {
		actualParsed := StimResponse{}
		err := actualParsed.LoadFromMem(&mem)
		assert.NoError(t, err)
		assert.Equal(t, sResponseParsed, actualParsed)
	})
	t.Run("ThresholdElectrotonus", func(t *testing.T) {
		actualParsed := ThresholdElectrotonus{}
		err := actualParsed.LoadFromMem(&mem)
		assert.NoError(t, err)
		assert.Equal(t, thresholdElectrotonusParsed, actualParsed)
	})
	t.Run("ThresholdIV", func(t *testing.T) {
		actualParsed := ThresholdIV{}
		err := actualParsed.LoadFromMem(&mem)
		assert.NoError(t, err)
		assert.Equal(t, thresholdIVParsed, actualParsed)
	})
	t.Run("RecoveryCycle", func(t *testing.T) {
		actualParsed := RecoveryCycle{}
		err := actualParsed.LoadFromMem(&mem)
		assert.NoError(t, err)
		assert.Equal(t, recoveryCycleParsed, actualParsed)
	})
	t.Run("ChargeDuration", func(t *testing.T) {
		actualParsed := ChargeDuration{}
		err := actualParsed.LoadFromMem(&mem)
		assert.NoError(t, err)
		assert.Equal(t, chargeDurationParsed, actualParsed)
	})
}

func TestImportFile(t *testing.T) {
	file, err := os.Open("../../res/data/short_test.MEM")
	assert.NoError(t, err)

	mem, err := Import(bufio.NewReader(file))

	assert.NoError(t, err)
	assert.Equal(t, completeExpectedRawMem, mem)
}
