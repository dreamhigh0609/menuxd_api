package dish

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetString(t *testing.T) {
	s := []string{"a", "b", "c"}
	result := SetString(s)
	fmt.Println(result)

	assert.Equal(t, "a,b,c", result)
}

func TestSetSlice(t *testing.T) {
	s := "a,b,c"
	result := SetSlice(s)
	fmt.Println(result)

	assert.Equal(t, []string{"a", "b", "c"}, result)
}
