package category

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidVoid(t *testing.T) {
	c := Category{}

	assert.True(t, c.IsValid())
}

func TestIsValid(t *testing.T) {
	var suggested1 uint = 1
	var suggested2 uint = 2
	var suggested3 uint = 3

	c := Category{}
	c.Suggested1 = &suggested1
	c.Suggested2 = &suggested2
	c.Suggested3 = &suggested3

	assert.True(t, c.IsValid())
}

func TestIsValidOne(t *testing.T) {
	var suggested1 uint = 1

	c := Category{}
	c.Suggested1 = &suggested1

	assert.True(t, c.IsValid())
}

func TestIsValidTwo(t *testing.T) {
	var suggested2 uint = 2

	c := Category{}
	c.Suggested2 = &suggested2

	assert.True(t, c.IsValid())
}

func TestIsValidFail(t *testing.T) {
	var suggested1 uint = 1
	var suggested2 uint = 1
	var suggested3 uint = 1

	c := Category{}
	c.Suggested1 = &suggested1
	c.Suggested2 = &suggested2
	c.Suggested3 = &suggested3

	assert.False(t, c.IsValid())
}
