package jitter

import (
	"bufio"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertMemHandlerNoParticipant(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(ConvertMemHandler).ServeHTTP(rr, req)

	status := rr.Code
	assert.Equal(t, http.StatusInternalServerError, status)
}

func TestConvertMemHandlerWithParticipant(t *testing.T) {
	file, err := os.Open("res/data/FESB70821B.MEM")
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "", bufio.NewReader(file))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(ConvertMemHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	// Just make sure something somewhat long is being printed. Great test, isn't it?
	assert.True(t, len(rr.Body.String()) > 100)
}

func TestConvertMemHandlerWithInvalidParticipant(t *testing.T) {
	file, err := os.Open("res/data/participants.json")
	assert.NoError(t, err)

	req, err := http.NewRequest("GET", "", bufio.NewReader(file))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(ConvertMemHandler).ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
