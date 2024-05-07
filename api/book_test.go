package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestObject(t *testing.T) {
	var a string = "Hello"
	var b string = "Hello"

	require.Equal(t, a, b, "The two words should be the same.")

	assert := assert.New(t)

	// assert equality
	assert.Equal(123, 123, "they should be equal")

	// assert inequality
	assert.NotEqual(123, 456, "they should not be equal")

	type forTest struct {
		Value string
	}
	object := forTest{"Something"}
	t.Log("object", object)

	// assert for nil (good for errors)
	// assert.Nil(t, object)

	// assert for not nil (good when you expect something)
	if assert.NotNil(object) {

		// now we know that object isn't nil, we are safe to make
		// further assertions without causing any errors
		assert.Equal("Something", object.Value)
	}
}

func TestOrderIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping kickstart integration test")
	}

	var a string = "Hello"
	var b string = "Hello"

	require.Equal(t, a, b, "The two words should be the same.")

	assert := assert.New(t)

	// assert equality
	assert.Equal(123, 123, "they should be equal")

	// assert inequality
	assert.NotEqual(123, 456, "they should not be equal")

	type forTest struct {
		Value string
	}
	object := forTest{"Something"}
	t.Log("object", object)

	// assert for nil (good for errors)
	// assert.Nil(t, object)

	// assert for not nil (good when you expect something)
	if assert.NotNil(object) {

		// now we know that object isn't nil, we are safe to make
		// further assertions without causing any errors
		assert.Equal("Something", object.Value)
	}
}
