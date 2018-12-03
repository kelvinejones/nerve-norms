package mem

import (
	"bufio"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const headerString = ` File:              	n:\Qtrac\Data\Human Normative data\Median nerve raw\FESB70821A.QZD
 Name:              	CR21S
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

var headerExpected = Header{
	File:          `n:\Qtrac\Data\Human Normative data\Median nerve raw\FESB70821A.QZD`,
	Name:          "CR21S",
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
	assert.Equal(t, m, Mem{})
}

func TestImportHeader(t *testing.T) {
	header := Header{}
	err := parseHeader(NewStringReader(headerString), &header)
	assert.NoError(t, err)
	assert.Equal(t, headerExpected, header)
}

const sResponseString = `

 STIMULUS-RESPONSE DATA (2.4-1.9m)

Values are those recorded

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

var sResponseExpected = StimResponse{
	MaxCmap: 1.161296,
	Values: []XY{
		XY{X: 2, Y: 3.915578},
		XY{X: 4, Y: 4.073214},
		XY{X: 6, Y: 4.144141},
		XY{X: 8, Y: 4.20404},
		XY{X: 10, Y: 4.435846},
		XY{X: 12, Y: 4.601757},
		XY{X: 14, Y: 4.824213},
		XY{X: 16, Y: 4.86682},
		XY{X: 18, Y: 4.89536},
		XY{X: 20, Y: 4.9239},
	},
}

func TestImportSRResponse(t *testing.T) {
	sResp := StimResponse{}
	err := parseStimResponse(NewStringReader(sResponseString), &sResp)
	assert.NoError(t, err)
	assert.Equal(t, sResponseExpected, sResp)
}

const chargeDurationString = `

  CHARGE DURATION DATA (2.4-3.5m)

                    	Duration (ms)       	 Threshold (mA)     	  Threshold charge (mA.mS)
QT.1                	 .2                 	 9.790961           	 1.958192
QT.2                	 .4                 	 6.905862           	 2.762345
QT.3                	 .6                 	 5.978864           	 3.587318
QT.4                	 .8                 	 5.44341            	 4.354728
QT.5                	 1                  	 5.187509           	 5.187509


`

var chargeDurationExpected = ChargeDuration{
	Values: []XYZ{
		XYZ{X: .2, Y: 9.790961, Z: 1.958192},
		XYZ{X: .4, Y: 6.905862, Z: 2.762345},
		XYZ{X: .6, Y: 5.978864, Z: 3.587318},
		XYZ{X: .8, Y: 5.44341, Z: 4.354728},
		XYZ{X: 1., Y: 5.187509, Z: 5.187509},
	},
}

func TestImportChargeDuration(t *testing.T) {
	sResp := ChargeDuration{}
	err := parseChargeDuration(NewStringReader(chargeDurationString), &sResp)
	assert.NoError(t, err)
	assert.Equal(t, chargeDurationExpected, sResp)
}

const thresholdElectrotonusString = `

  THRESHOLD ELECTROTONUS DATA (3.5-8.8m)

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

var thresholdElectrotonusExpected = ThresholdElectrotonusGroup{
	Sets: []ThresholdElectrotonusSet{
		ThresholdElectrotonusSet{Values: []XYZ{
			XYZ{X: 0, Y: 0, Z: 0.00},
			XYZ{X: 9, Y: 0, Z: 0.00},
			XYZ{X: 10, Y: 40, Z: 40.02},
			XYZ{X: 11, Y: 40, Z: 42.71},
			XYZ{X: 11, Y: 40, Z: -42.71},
		}},
		ThresholdElectrotonusSet{Values: []XYZ{
			XYZ{X: 0, Y: 0, Z: 0.00},
			XYZ{X: 9, Y: 0, Z: 0.00},
			XYZ{X: 10, Y: -40, Z: -39.34},
			XYZ{X: 11, Y: -40, Z: -40.92},
		}},
	},
}

func TestImportThresholdElectrotonus(t *testing.T) {
	actual := ThresholdElectrotonusGroup{}
	err := parseThresholdElectrotonus(NewStringReader(thresholdElectrotonusString), &actual)
	assert.NoError(t, err)
	assert.Equal(t, thresholdElectrotonusExpected, actual)
}

const recoveryCycleString = `

  RECOVERY CYCLE DATA (11.1-15.3m)

                    	Interval (ms)       	  Threshold change (%)
RC1.1               	 3.2                	4.99
RC1.2               	 4                  	-12.75
RC1.3               	 5                  	-22.24
RC1.4               	 6.3                	-24.45
RC1.5               	 7.9                	-24.05

`

var recoveryCycleExpected = RecoveryCycle{
	Values: []XY{
		XY{X: 3.2, Y: 4.99},
		XY{X: 4, Y: -12.75},
		XY{X: 5, Y: -22.24},
		XY{X: 6.3, Y: -24.45},
		XY{X: 7.9, Y: -24.05},
	},
}

func TestImportRecoveryCycle(t *testing.T) {
	actual := RecoveryCycle{}
	err := parseRecoveryCycle(NewStringReader(recoveryCycleString), &actual)
	assert.NoError(t, err)
	assert.Equal(t, recoveryCycleExpected, actual)
}

const thresholdIVString = `

  THRESHOLD I/V DATA (8.9-11m)

                    	Current (%)         	  Threshold redn. (%)
IV1.1               	 50                 	49.28
IV1.2               	 40                 	39.01
IV1.3               	 30                 	31.59
IV1.4               	 20                 	22.58
IV1.5               	 10                 	13.06
IV1.6               	 0                  	-0.78

`

var thresholdIVExpected = ThresholdIV{
	Values: []XY{
		XY{X: 50, Y: 49.28},
		XY{X: 40, Y: 39.01},
		XY{X: 30, Y: 31.59},
		XY{X: 20, Y: 22.58},
		XY{X: 10, Y: 13.06},
		XY{X: 0, Y: -0.78},
	},
}

func TestImportThresholdIV(t *testing.T) {
	actual := ThresholdIV{}
	err := parseThresholdIV(NewStringReader(thresholdIVString), &actual)
	assert.NoError(t, err)
	assert.Equal(t, thresholdIVExpected, actual)
}

const excitabilityVariablesString = `

  DERIVED EXCITABILITY VARIABLES

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

`

var excitabilityVariablesExpected = ExcitabilityVariables{
	Values: map[string]float64{
		`Stimulus (mA) for 50% max response`:   5.307,
		`Strength-duration\time constant (ms)`: 0.287,
		`Polarizing current\(mA)`:              0.000,
		`Sex (M=1, F=2)`:                       1,
		`TEh(10-20ms)`:                         -64.192,
		`TEd40(Accom)`:                         19.6,
		`TEd20(10-20ms)`:                       30.8,
		`TEh20(10-20ms)`:                       -32.2,
	},
	Program:         "QTracP 9/12/2016",
	ThresholdMethod: 6,
	SRMethod:        1,
}

func TestImportExcitabilityVariables(t *testing.T) {
	actual := ExcitabilityVariables{Values: make(map[string]float64)}
	err := parseExcitabilityVariables(NewStringReader(excitabilityVariablesString), &actual)
	assert.NoError(t, err)
	assert.Equal(t, excitabilityVariablesExpected, actual)
}

var completeExpectedMem = Mem{
	Header:                     headerExpected,
	StimResponse:               sResponseExpected,
	ChargeDuration:             chargeDurationExpected,
	ThresholdElectrotonusGroup: thresholdElectrotonusExpected,
	RecoveryCycle:              recoveryCycleExpected,
	ThresholdIV:                thresholdIVExpected,
	ExcitabilityVariables:      excitabilityVariablesExpected,
}

func TestImportAll(t *testing.T) {
	memString := headerString + sResponseString + chargeDurationString + thresholdElectrotonusString + recoveryCycleString + thresholdIVString + excitabilityVariablesString
	mem, err := Import(strings.NewReader(memString))

	assert.NoError(t, err)
	assert.Equal(t, completeExpectedMem, mem)
}

func TestImportFile(t *testing.T) {
	file, err := os.Open("../../res/data/short_test.MEM")
	assert.NoError(t, err)

	mem, err := Import(bufio.NewReader(file))

	assert.NoError(t, err)
	assert.Equal(t, completeExpectedMem, mem)
}
