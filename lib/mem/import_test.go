package mem

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportEmpty(t *testing.T) {
	m, err := Import(strings.NewReader(""))
	assert.NoError(t, err)
	assert.Equal(t, m, Mem{})
}
