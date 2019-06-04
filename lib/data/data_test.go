package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalParticipants(t *testing.T) {
	_, err := AsMef()
	assert.NoError(t, err)
}
