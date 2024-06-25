package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMaze(t *testing.T) {
	maze := NewMaze()
	assert.Equal(t, len(maze.cells), colCount*rowCount)
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
			maze := NewMaze()
			maze.HideWall(tc.location, tc.side)
			assert.False(t, maze.cells[tc.location].walls[tc.side].Visible())
			assert.False(t, maze.cells[tc.neighborLocation].walls[tc.neighborSide].Visible())
		})
	}
}

func TestHideWallPanics(t *testing.T) {
	tests := map[string]struct {
		location Location
		side     Side
	}{
		"top":    {Location{X: 10, Y: 0}, top},
		"right":  {Location{X: colCount - 1, Y: 2}, right},
		"bottom": {Location{X: 10, Y: rowCount - 1}, bottom},
		"left":   {Location{X: 0, Y: 5}, left},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			maze := NewMaze()
			assert.Panics(t, func() { maze.HideWall(tc.location, tc.side) })
		})
	}
}
