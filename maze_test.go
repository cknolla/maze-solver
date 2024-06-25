package main

import (
	"github.com/stretchr/testify/assert"
	"math/rand/v2"
	"testing"
)

func TestNewMaze(t *testing.T) {
	maze := NewMaze(40, 30, 20, rand.NewPCG(0, 0))
	assert.Equal(t, len(maze.cells), maze.colCount*maze.rowCount)
}

func TestHideWall(t *testing.T) {
	tests := map[string]struct {
		location         Location
		side             Side
		neighborLocation Location
		neighborSide     Side
	}{
		"top":    {Location{X: 10, Y: 5}, top, Location{X: 10, Y: 4}, bottom},
		"right":  {Location{X: 10, Y: 5}, right, Location{X: 11, Y: 5}, left},
		"bottom": {Location{X: 10, Y: 5}, bottom, Location{X: 10, Y: 6}, top},
		"left":   {Location{X: 10, Y: 5}, left, Location{X: 9, Y: 5}, right},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			maze := NewMaze(40, 30, 20, rand.NewPCG(0, 0))
			maze.hideWall(tc.location, tc.side)
			assert.False(t, maze.cells[tc.location].walls[tc.side].Visible())
			assert.False(t, maze.cells[tc.neighborLocation].walls[tc.neighborSide].Visible())
		})
	}
}

func TestHideWallPanics(t *testing.T) {
	tests := map[string]struct {
		location Location
	}{
		"top":    {Location{X: 0, Y: -1}},
		"right":  {Location{X: 40, Y: 0}},
		"bottom": {Location{X: 0, Y: 30}},
		"left":   {Location{X: -1, Y: 0}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			maze := NewMaze(40, 30, 20, rand.NewPCG(0, 0))
			assert.Panics(t, func() { maze.hideWall(tc.location, left) })
		})
	}
}

func TestSetVisited(t *testing.T) {
	maze := NewMaze(5, 5, 5, rand.NewPCG(0, 0))
	cell := maze.cells[Location{X: 0, Y: 0}]
	cell.visited = true
	assert.True(t, maze.cells[Location{X: 0, Y: 0}].visited)
}

func TestMaze_removeWallsR(t *testing.T) {
	m := NewMaze(10, 10, 10, rand.NewPCG(0, 0))
	m.removeWallsR(Location{X: 0, Y: 0}, top)
	for _, cell := range m.cells {
		assert.True(t, cell.visited)
	}
	m.resetVisited()
	for _, cell := range m.cells {
		assert.False(t, cell.visited)
	}
}

func TestMaze_solveR(t *testing.T) {
	m := NewMaze(10, 10, 10, rand.NewPCG(0, 0))
	m.removeWallsR(Location{X: 0, Y: 0}, top)
	m.resetVisited()
	m.solveR(Location{X: 0, Y: 0})
	m.cells[Location{X: 0, Y: 0}].visited = true
	m.cells[Location{X: m.colCount - 1, Y: m.rowCount - 1}].visited = true
}
