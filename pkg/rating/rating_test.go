package rating

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	ratings := []Rating{
		{
			Score: 1,
		},
		{
			Score: 2,
		},
		{
			Score: 3,
		},
		{
			Score: 4,
		},
		{
			Score: 5,
		},
	}

	assert := assert.New(t)

	for _, r := range ratings {
		assert.True(r.IsValid())
	}
}

func TestIsValidFail(t *testing.T) {
	ratings := []Rating{
		{
			Score: 0,
		},
		{
			Score: 6,
		},
		{
			Score: 8,
		},
	}

	assert := assert.New(t)

	for _, r := range ratings {
		assert.False(r.IsValid())
	}
}
