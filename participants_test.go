package jitter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gogs.bellstone.ca/james/jitter/lib/data"
)

func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(ParticipantHandler).ServeHTTP(rr, req)

	status := rr.Code
	assert.Equal(t, http.StatusOK, status)

	assert.Equal(t, data.Participants+"\n", rr.Body.String())
}
