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

	str, err := reader.ReadLine()
	assert.Equal(t, "Hello\n", str)
	assert.NoError(t, err)

	str, err = reader.ReadLine()
	assert.Equal(t, "Line 2\n", str)
	assert.NoError(t, err)

	str, err = reader.ReadLine()
	assert.Equal(t, "\n", str)
	assert.NoError(t, err)

	str, err = reader.ReadLine()
	assert.Equal(t, "\n", str)
	assert.NoError(t, err)

	str, err = reader.ReadLine()
	assert.Equal(t, "I am a test\n", str)
	assert.NoError(t, err)
}

func TestReaderUnread(t *testing.T) {
	reader := NewStringReader(testString)

	str, err := reader.ReadLine()
	assert.Equal(t, "Hello\n", str)
	assert.NoError(t, err)

	reader.UnreadString("Another test\n")
	str, err = reader.ReadLine()
	assert.Equal(t, "Another test\n", str)
	assert.NoError(t, err)

	str, err = reader.ReadLine()
	assert.Equal(t, "Line 2\n", str)
	assert.NoError(t, err)
}
