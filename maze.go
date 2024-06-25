package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"time"
)

type Side int

const (
	top Side = iota
	right
	bottom
	left
)

type Location struct {
	X int
	Y int
}

type Cell struct {
	walls    map[Side]*canvas.Line
	location Location
	size     float32
	visited  bool
}

func NewCell(location Location, size float32) Cell {
	cell := Cell{
		walls:    make(map[Side]*canvas.Line, 4),
		location: location,
		size:     size,
	}
	position := cell.Position()
	cell.walls[top] = canvas.NewLine(color.White)
	cell.walls[top].Position1 = fyne.NewPos(position.X, position.Y)
	cell.walls[top].Position2 = fyne.NewPos(position.X+size, position.Y)
	cell.walls[right] = canvas.NewLine(color.White)
	cell.walls[right].Position1 = fyne.NewPos(position.X+size, position.Y)
	cell.walls[right].Position2 = fyne.NewPos(position.X+size, position.Y+size)
	cell.walls[bottom] = canvas.NewLine(color.White)
	cell.walls[bottom].Position1 = fyne.NewPos(position.X, position.Y+size)
	cell.walls[bottom].Position2 = fyne.NewPos(position.X+size, position.Y+size)
	cell.walls[left] = canvas.NewLine(color.White)
	cell.walls[left].Position1 = fyne.NewPos(position.X, position.Y)
	cell.walls[left].Position2 = fyne.NewPos(position.X, position.Y+size)
	return cell
}

func (cell Cell) Position() fyne.Position {
	return fyne.NewPos(float32(cell.location.X)*cell.size, float32(cell.location.Y)*cell.size)
}

func (cell Cell) Center() fyne.Position {
	position := cell.Position()
	return fyne.NewPos((position.X+position.X+cell.size)/2, (position.Y+position.Y+cell.size)/2)
}

type Maze struct {
	colCount int
	rowCount int
	cells    map[Location]Cell
	cont     *fyne.Container
	cellSize float32
}

func NewMaze(colCount, rowCount int, cellSize float32) Maze {
	m := Maze{
		colCount: colCount,
		rowCount: rowCount,
		cells:    make(map[Location]Cell, colCount*rowCount),
		cont:     container.NewWithoutLayout(),
		cellSize: cellSize,
	}
	for y := 0; y < rowCount; y++ {
		for x := 0; x < colCount; x++ {
			cell := NewCell(Location{X: x, Y: y}, cellSize)
			m.cells[Location{X: x, Y: y}] = cell
		}
	}
	return m
}

func (m *Maze) DrawCells(w *fyne.Window) {
	m.cont.RemoveAll()
	for _, cell := range m.cells {
		for _, wall := range cell.walls {
			m.cont.Add(wall)
		}
	}
	m.cont.Resize(fyne.NewSize(float32(m.colCount)*m.cellSize, float32(m.rowCount)*m.cellSize))
	(*w).SetContent(m.cont)
	(*w).Resize(fyne.NewSize(float32(m.colCount)*m.cellSize+5, float32(m.rowCount)*m.cellSize+5))
}

func (m *Maze) DrawLine(source, target Location, undo bool) {
	lineColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	if undo {
		lineColor = color.RGBA{R: 100, G: 100, B: 100, A: 255}
	}
	line := canvas.NewLine(lineColor)
	line.StrokeWidth = m.cellSize / 8
	line.Position1 = m.cells[source].Center()
	line.Position2 = m.cells[target].Center()
	m.cont.Add(line)
}

// HideWall will make the specified location's wall invisible
// If the location has a neighbor, its corresponding wall will be hidden as well.
// For example, if the left wall of location 1,1 is hidden,
// the right wall of location 0,1 must also be hidden to not render that line.
// If the location's wall is on an outer edge, only the single wall will be hidden.
func (m *Maze) HideWall(location Location, side Side) {
	sourceCell, ok := m.cells[location]
	if !ok {
		panic(fmt.Sprintf("location %v does not exist", location))
	}
	switch side {
	case top:
		sourceCell.walls[top].Hide()
		if location.Y > 0 {
			m.cells[Location{X: location.X, Y: location.Y - 1}].walls[bottom].Hide()
		}
	case right:
		sourceCell.walls[right].Hide()
		if location.X < m.colCount-1 {
			m.cells[Location{X: location.X + 1, Y: location.Y}].walls[left].Hide()
		}
	case bottom:
		sourceCell.walls[bottom].Hide()
		if location.Y < m.rowCount-1 {
			m.cells[Location{X: location.X, Y: location.Y + 1}].walls[top].Hide()
		}
	case left:
		sourceCell.walls[left].Hide()
		if location.X > 0 {
			m.cells[Location{X: location.X - 1, Y: location.Y}].walls[right].Hide()
		}
	}
}

func (m *Maze) RemoveWalls() {
	time.Sleep(1 * time.Second)
	m.cells[Location{X: 0, Y: 0}].walls[top].Hide()
	m.cells[Location{X: m.colCount - 1, Y: m.rowCount - 1}].walls[bottom].Hide()
}
