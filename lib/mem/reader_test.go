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
	assert.Equal(t, "Hello", str)
	assert.NoError(t, err)

	str, err = reader.ReadLine()
	assert.Equal(t, "Line 2", str)
	assert.NoError(t, err)

	str, err = reader.ReadLine()
	assert.Equal(t, "", str)
	assert.NoError(t, err)

	str, err = reader.ReadLine()
	assert.Equal(t, "", str)
	assert.NoError(t, err)

	str, err = reader.ReadLine()
	assert.Equal(t, "I am a test", str)
	assert.NoError(t, err)
}

func TestReaderUnread(t *testing.T) {
	reader := NewStringReader(testString)

	str, err := reader.ReadLine()
	assert.Equal(t, "Hello", str)
	assert.NoError(t, err)

	reader.UnreadString("Another test")
	str, err = reader.ReadLine()
	assert.Equal(t, "Another test", str)
	assert.NoError(t, err)

	str, err = reader.ReadLine()
	assert.Equal(t, "Line 2", str)
	assert.NoError(t, err)
}
