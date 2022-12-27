package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRight(t *testing.T) {
	rbr := &rightBlizzardRing{
		m: map[int]struct{}{
			1: {},
		},
		len: 4,
	}

	// t:0
	assert.False(t, rbr.hasBlizzard(0, 0))
	assert.True(t, rbr.hasBlizzard(1, 0))
	assert.False(t, rbr.hasBlizzard(2, 0))
	assert.False(t, rbr.hasBlizzard(3, 0))

	// t:1
	assert.False(t, rbr.hasBlizzard(0, 1))
	assert.False(t, rbr.hasBlizzard(1, 1))
	assert.True(t, rbr.hasBlizzard(2, 1))
	assert.False(t, rbr.hasBlizzard(3, 1))

	// t:4 - should be same as t:0
	assert.False(t, rbr.hasBlizzard(0, 4))
	assert.True(t, rbr.hasBlizzard(1, 4))
	assert.False(t, rbr.hasBlizzard(2, 4))
	assert.False(t, rbr.hasBlizzard(3, 4))

	// t:5 - should be same as t:1
	assert.False(t, rbr.hasBlizzard(0, 5))
	assert.False(t, rbr.hasBlizzard(1, 5))
	assert.True(t, rbr.hasBlizzard(2, 5))
	assert.False(t, rbr.hasBlizzard(3, 5))

	// t:40 - should be same as t:0
	assert.False(t, rbr.hasBlizzard(0, 40))
	assert.True(t, rbr.hasBlizzard(1, 40))
	assert.False(t, rbr.hasBlizzard(2, 40))
	assert.False(t, rbr.hasBlizzard(3, 40))

	// t:41 - should be same as t:1
	assert.False(t, rbr.hasBlizzard(0, 41))
	assert.False(t, rbr.hasBlizzard(1, 41))
	assert.True(t, rbr.hasBlizzard(2, 41))
	assert.False(t, rbr.hasBlizzard(3, 41))
}

func TestLeft(t *testing.T) {
	br := &blizzardRing{
		m: map[int]struct{}{
			2: {},
		},
		len: 4,
	}

	// t:0
	assert.False(t, br.hasBlizzard(0, 0))
	assert.False(t, br.hasBlizzard(1, 0))
	assert.True(t, br.hasBlizzard(2, 0))
	assert.False(t, br.hasBlizzard(3, 0))

	// t:1
	assert.False(t, br.hasBlizzard(0, 1))
	assert.True(t, br.hasBlizzard(1, 1))
	assert.False(t, br.hasBlizzard(2, 1))
	assert.False(t, br.hasBlizzard(3, 1))

	// t:4 - should be same as t:0
	assert.False(t, br.hasBlizzard(0, 0))
	assert.False(t, br.hasBlizzard(1, 0))
	assert.True(t, br.hasBlizzard(2, 0))
	assert.False(t, br.hasBlizzard(3, 0))

	// t:5 - should be same as t:1
	assert.False(t, br.hasBlizzard(0, 1))
	assert.True(t, br.hasBlizzard(1, 1))
	assert.False(t, br.hasBlizzard(2, 1))
	assert.False(t, br.hasBlizzard(3, 1))
}
