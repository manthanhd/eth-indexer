package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_StripZeroPadding(t *testing.T) {
	hexInput := "000000000000000000000000000000000000000000000000000000174876fde4"
	strippedInput := StripZeroPadding(hexInput)
	assert.Equal(t, "174876fde4", strippedInput)
}
