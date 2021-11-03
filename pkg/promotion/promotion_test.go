package promotion

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetDays(t *testing.T) {
	daysString := "monday,sunday"
	result := SetDays(daysString)

	assert.Equal(t, []string{"monday", "sunday"}, result)
}

func TestSetDaysString(t *testing.T) {
	days := []string{"monday", "sunday"}
	result := SetDaysString(days)

	assert.Equal(t, "monday,sunday", result)
}
