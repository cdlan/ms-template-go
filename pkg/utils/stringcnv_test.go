package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStringToBool(t *testing.T) {

	str := "true"

	res, err := StringToBool(str)

	assert.Nil(t, err)
	assert.True(t, res)
}

func TestStringToInt(t *testing.T) {

	num := 999
	str := fmt.Sprintf("%d", num)

	res, err := StringToInt(str)

	assert.Nil(t, err)
	assert.Equal(t, num, res)
}

func TestStringToUrl(t *testing.T) {

	str := "http://localhost:8080"

	add, err := StringToUrl(str)

	assert.Nil(t, err)
	assert.Equal(t, str, add.String())
}
