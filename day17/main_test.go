package main

import (
	"testing"

	"kfet.org/aoc_common/assert"
)

func TestMove(t *testing.T) {
	w := NewWorld("data/part_one.txt")

	assert.Equals(0, w.r.moveType, "")
	w.step()
	assert.Equals(1, w.r.moveType, "")

	w.nextRock()
	assert.Equals(0, w.r.moveType, "")
	w.step()
	assert.Equals(1, w.r.moveType, "")
}

func TestChamberTestMove(t *testing.T) {
	w := NewWorld("data/part_one.txt")
	r := NewRock(w.rockSprites[0], 0, 0)

	ch := NewChamber(7)

	assert.False(ch.testMove(r, move{10, 10}))
	assert.False(ch.testMove(r, move{-1, 10}))
	assert.True(ch.testMove(r, move{0, 10}))
	assert.True(ch.testMove(r, move{3, 10}))
	assert.True(ch.testMove(r, move{2, 10}))
	assert.True(ch.testMove(r, move{2, 0}))
	assert.False(ch.testMove(r, move{2, -1}))
}

func TestChamberStampRock(t *testing.T) {
	rStamp := NewRock(&mask{
		{0, 1, 0},
		{1, 1, 1},
		{0, 1, 0}}, 0, 2)
	rTest := NewRock(&mask{
		{1}}, 0, 0)

	ch := NewChamber(7)

	ch.stampRock(rStamp)

	assert.True(ch.testMove(rTest, move{0, 2}))
	assert.False(ch.testMove(rTest, move{1, 2}))
	assert.True(ch.testMove(rTest, move{2, 2}))

	assert.False(ch.testMove(rTest, move{0, 1}))
	assert.False(ch.testMove(rTest, move{1, 1}))
	assert.False(ch.testMove(rTest, move{2, 1}))

	assert.True(ch.testMove(rTest, move{0, 0}))
	assert.False(ch.testMove(rTest, move{1, 0}))
	assert.True(ch.testMove(rTest, move{2, 0}))
}
