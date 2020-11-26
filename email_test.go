package jitter

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendContactMail(t *testing.T) {
	/*req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}*/

	pack := EmailPackage{
		Name:       "Go Test",
		Sender:     "GoTest@fakemail.com",
		Subject:    "Test Email from Golang test framework",
		Message:    "This is a test email from golang testing framework",
		CarbonCopy: false,
	}

	jsonBytes, err := json.Marshal(pack)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "", bytes.NewReader(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	http.HandlerFunc(ContactEmailHandler).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	//SendContactMail(pack)
}
