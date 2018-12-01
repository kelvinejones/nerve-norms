package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportEmpty(t *testing.T) {
	m, err := Import("")
	assert.NoError(t, err)
	assert.Equal(t, m, Mem{})
}
