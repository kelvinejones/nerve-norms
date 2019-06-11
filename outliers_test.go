package jitter

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
	file, err := os.Open("res/data/FESB70821B.MEM")
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "", bufio.NewReader(file))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(OutlierScoreHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Just make sure something somewhat long is being printed. Great test, isn't it?
	assert.True(t, len(rr.Body.String()) > 100)
}
