package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoard_CountNeighbours(t *testing.T) {
	b := &board{}
	b.createGrid(3, 3)
	b.currentGenCells[0][0] = true
	b.currentGenCells[0][2] = true
	b.currentGenCells[2][0] = true
	b.currentGenCells[2][2] = true

	assert.Equal(t, 4, b.countNeighbours(1, 1))
}

func TestBoard_CountNeighbours_Corner(t *testing.T) {
	b := &board{}
	b.createGrid(2, 2)
	b.currentGenCells[0][1] = true
	b.currentGenCells[1][0] = true
	b.currentGenCells[1][1] = true

	assert.Equal(t, 3, b.countNeighbours(0, 0))
}

func TestBoard_CountNeighbours_IgnoresMiddle(t *testing.T) {
	b := &board{}
	b.createGrid(3, 3)
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			b.currentGenCells[y][x] = true
		}
	}

	assert.Equal(t, 8, b.countNeighbours(1, 1))
}

func TestBoard_CreateGrid(t *testing.T) {
	b := &board{}
	b.createGrid(5, 3)

	assert.Equal(t, 3, len(b.currentGenCells))
	assert.Equal(t, 5, len(b.currentGenCells[0]))
}

func TestBoard_EnsureGridSize(t *testing.T) {
	b := &board{}
	b.createGrid(2, 2)
	b.currentGenCells[1][1] = true

	b.ensureGridSize(4, 4)
	assert.Equal(t, 4, len(b.currentGenCells))
	assert.Equal(t, 4, len(b.currentGenCells[0]))
	assert.True(t, b.currentGenCells[1][1])
}

func TestBoard_Generatoin(t *testing.T) {
	b := &board{}
	b.createGrid(5, 5)
	assert.Equal(t, 0, b.generation)

	b.nextGen()
	assert.Equal(t, 1, b.generation)
}

func BenchmarkBoard_NextGeneration(b *testing.B) {
	board := &board{}
	board.createGrid(15, 16)

	// 0x0
	board.currentGenCells[6][7] = true

	// x0x
	// x0x
	board.currentGenCells[7][6] = true
	board.currentGenCells[7][8] = true
	board.currentGenCells[8][6] = true
	board.currentGenCells[8][8] = true

	// 0x0
	board.currentGenCells[9][7] = true

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		board.nextGen()
	}
}
