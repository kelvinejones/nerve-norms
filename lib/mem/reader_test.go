package mem

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const testString = `Hello
Line 2


I am a test
`

func TestBasicReader(t *testing.T) {
	reader := NewStringReader(testString)

	str, err := reader.ReadString('\n')
	assert.Equal(t, "Hello\n", str)
	assert.NoError(t, err)

	str, err = reader.ReadString('\n')
	assert.Equal(t, "Line 2\n", str)
	assert.NoError(t, err)

	str, err = reader.ReadString('\n')
	assert.Equal(t, "\n", str)
	assert.NoError(t, err)

	str, err = reader.ReadString('\n')
	assert.Equal(t, "\n", str)
	assert.NoError(t, err)

	str, err = reader.ReadString('\n')
	assert.Equal(t, "I am a test\n", str)
	assert.NoError(t, err)
}
