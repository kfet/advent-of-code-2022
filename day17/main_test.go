package main

import (
	"testing"

	"kfet.org/aoc_common/assert"
)

func TestMove(t *testing.T) {
	w := NewWorld("data/part_one.txt")

	assert.Equals(0, w.rockMove, "")
	w.step()
	assert.Equals(1, w.rockMove, "")

	w.nextRock()
	assert.Equals(0, w.rockMove, "")
	w.step()
	assert.Equals(1, w.rockMove, "")
}

func TestChamberTestMove(t *testing.T) {
	w := NewWorld("data/part_one.txt")
	r := w.rocks[0]

	ch := NewChamber(7)

	assert.False(ch.testMove(r, 10, 10))
	assert.False(ch.testMove(r, -1, 10))
	assert.True(ch.testMove(r, 0, 10))
	assert.True(ch.testMove(r, 3, 10))
	assert.True(ch.testMove(r, 2, 10))
	assert.True(ch.testMove(r, 2, 0))
	assert.False(ch.testMove(r, 2, -1))
}

func TestChamberStampRock(t *testing.T) {
	rStamp := rock{
		{0, 1, 0},
		{1, 1, 1},
		{0, 1, 0},
	}
	rTest := rock{
		{1},
	}
	ch := NewChamber(7)

	ch.stampRock(&rStamp, 0, 2)

	assert.True(ch.testMove(&rTest, 0, 2))
	assert.False(ch.testMove(&rTest, 1, 2))
	assert.True(ch.testMove(&rTest, 2, 2))

	assert.False(ch.testMove(&rTest, 0, 1))
	assert.False(ch.testMove(&rTest, 1, 1))
	assert.False(ch.testMove(&rTest, 2, 1))

	assert.True(ch.testMove(&rTest, 0, 0))
	assert.False(ch.testMove(&rTest, 1, 0))
	assert.True(ch.testMove(&rTest, 2, 0))
}
