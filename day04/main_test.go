package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullOverlap(t *testing.T) {
	assert.True(t, fullOverlap(secRange{5, 5}, secRange{5, 5}))
	assert.True(t, fullOverlap(secRange{5, 5}, secRange{5, 100}))
	assert.True(t, fullOverlap(secRange{5, 100}, secRange{5, 5}))
	assert.True(t, fullOverlap(secRange{5, 100}, secRange{100, 100}))
	assert.True(t, fullOverlap(secRange{5, 10}, secRange{6, 9}))
	assert.True(t, fullOverlap(secRange{5, 10}, secRange{4, 11}))

	assert.False(t, fullOverlap(secRange{5, 50}, secRange{50, 100}))
	assert.False(t, fullOverlap(secRange{5, 50}, secRange{1, 5}))
	assert.False(t, fullOverlap(secRange{5, 50}, secRange{10, 200}))
	assert.False(t, fullOverlap(secRange{5, 50}, secRange{1, 10}))
}

func TestAnyOverlap(t *testing.T) {
	assert.True(t, anyOverlap(secRange{5, 5}, secRange{5, 5}))
	assert.True(t, anyOverlap(secRange{5, 100}, secRange{100, 200}))
	assert.True(t, anyOverlap(secRange{5, 100}, secRange{1, 5}))
	assert.False(t, anyOverlap(secRange{5, 100}, secRange{101, 200}))
}
