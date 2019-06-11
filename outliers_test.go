package jitter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gogs.bellstone.ca/james/jitter/lib/data"
)

func TestOutliersHandlerNoParticipant(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(OutlierScoreHandler).ServeHTTP(rr, req)

	status := rr.Code
	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestOutliersHandlerWithParticipantName(t *testing.T) {
	req, err := http.NewRequest("GET", "?name=CA-CR21S", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(OutlierScoreHandler).ServeHTTP(rr, req)

	status := rr.Code
	assert.Equal(t, http.StatusOK, status)

	// Just make sure something somewhat long is being printed. Great test, isn't it?
	assert.True(t, len(rr.Body.String()) > 100)
}

func TestConvertMemHandlerWithParticipantData(t *testing.T) {
	mefData, err := data.AsMef()
	assert.NoError(t, err)

	memData := mefData.MemWithKey("CA-CR21S")
	jsMem, err := json.Marshal(memData)
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "", bytes.NewReader(jsMem))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(OutlierScoreHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Just make sure something somewhat long is being printed. Great test, isn't it?
	assert.True(t, len(rr.Body.String()) > 100)
}
