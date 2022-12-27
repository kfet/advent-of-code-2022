package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPartOneSmall(t *testing.T) {
	res, err := processFile("data/part_one_small.txt", true)
	assert.Nil(t, err)
	assert.Equal(t, 25, res)
}
