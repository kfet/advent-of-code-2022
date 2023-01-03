package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	res, err := processFile("data/part_one.txt", geode, 24, true)
	assert.Nil(t, err)
	assert.Equal(t, 33, res)
}
