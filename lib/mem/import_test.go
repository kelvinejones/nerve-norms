package mem

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const header = ` File:              	n:\Qtrac\Data\Human Normative data\Median nerve raw\FESB70821A.QZD
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

func TestImportEmpty(t *testing.T) {
	t.Skip()
	m, err := Import(strings.NewReader(""))
	assert.NoError(t, err)
	assert.Equal(t, m, Mem{})
}

func TestImportHeader(t *testing.T) {
	m, err := Import(strings.NewReader(header))
	assert.NoError(t, err)
	assert.Equal(t, m, Mem{MemHeader: MemHeader{
		File:          `n:\Qtrac\Data\Human Normative data\Median nerve raw\FESB70821A.QZD`,
		Name:          "CR21S",
		Protocol:      "TRONDNF",
		Date:          time.Date(2017, time.Month(8), 21, 12, 57, 17, 0, time.UTC), // TODO eventually handle time zones better?
		Age:           30,
		Sex:           MaleSex,
		Temperature:   33.5,
		SRSites:       "Median Wr-APB",
		NormalControl: true,
		Operator:      "MS",
		Comment:       "smooth recording",
	}})
}
