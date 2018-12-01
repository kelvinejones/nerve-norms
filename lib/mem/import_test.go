package mem

import (
	"bufio"
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
	err := parseHeader(bufio.NewReader(strings.NewReader(headerString)), &header)
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
	err := parseStimResponse(bufio.NewReader(strings.NewReader(sResponseString)), &sResp)
	assert.NoError(t, err)
	assert.Equal(t, sResponseExpected, sResp)
}

func TestImportAll(t *testing.T) {
	memString := headerString + sResponseString
	memExpected := Mem{
		Header:       headerExpected,
		StimResponse: sResponseExpected,
	}

	mem, err := Import(strings.NewReader(memString))

	assert.NoError(t, err)
	assert.Equal(t, memExpected, mem)
}
