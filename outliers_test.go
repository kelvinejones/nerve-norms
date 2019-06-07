package jitter

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOutliersHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
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
