package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalParticipants(t *testing.T) {
	md, err := AsMef()
	assert.NotEqual(t, 0, len(md))
	assert.NoError(t, err)
}
