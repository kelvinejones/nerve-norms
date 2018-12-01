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
`

var sResponseExpected = StimResponse{
	MaxCmap: 1.161296,
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
