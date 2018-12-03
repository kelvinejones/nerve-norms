package mem

import (
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

func TestImportAll(t *testing.T) {
	memString := headerString + sResponseString + chargeDurationString + thresholdElectrotonusString
	memExpected := Mem{
		Header:                     headerExpected,
		StimResponse:               sResponseExpected,
		ChargeDuration:             chargeDurationExpected,
		ThresholdElectrotonusGroup: thresholdElectrotonusExpected,
	}

	mem, err := Import(strings.NewReader(memString))

	assert.NoError(t, err)
	assert.Equal(t, memExpected, mem)
}
