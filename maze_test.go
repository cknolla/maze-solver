package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMaze(t *testing.T) {
	maze := NewMaze(40, 30, 20)
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
			maze := NewMaze(40, 30, 20)
			maze.HideWall(tc.location, tc.side)
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
			maze := NewMaze(40, 30, 20)
			assert.Panics(t, func() { maze.HideWall(tc.location, left) })
		})
	}
}
