package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
)

const colCount = 40
const rowCount = 20
const cellSize = 50

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
	position fyne.Position
}

func NewCell(position fyne.Position) Cell {
	cell := Cell{
		walls:    make(map[Side]*canvas.Line, 4),
		position: position,
	}
	cell.walls[top] = canvas.NewLine(color.White)
	cell.walls[top].Position1 = fyne.NewPos(position.X, position.Y)
	cell.walls[top].Position2 = fyne.NewPos(position.X+cellSize, position.Y)
	cell.walls[right] = canvas.NewLine(color.White)
	cell.walls[right].Position1 = fyne.NewPos(position.X+cellSize, position.Y)
	cell.walls[right].Position2 = fyne.NewPos(position.X+cellSize, position.Y+cellSize)
	cell.walls[bottom] = canvas.NewLine(color.White)
	cell.walls[bottom].Position1 = fyne.NewPos(position.X, position.Y+cellSize)
	cell.walls[bottom].Position2 = fyne.NewPos(position.X+cellSize, position.Y+cellSize)
	cell.walls[left] = canvas.NewLine(color.White)
	cell.walls[left].Position1 = fyne.NewPos(position.X, position.Y)
	cell.walls[left].Position2 = fyne.NewPos(position.X, position.Y+cellSize)
	return cell
}

func (cell Cell) Center() fyne.Position {
	return fyne.NewPos((cell.position.X+cell.position.X+cellSize)/2, (cell.position.Y+cell.position.Y+cellSize)/2)
}

type Maze struct {
	cells map[Location]Cell
	cont  *fyne.Container
}

func NewMaze() Maze {
	m := Maze{
		cells: make(map[Location]Cell, colCount*rowCount),
		cont:  container.NewWithoutLayout(),
	}
	for y := 0; y < rowCount; y++ {
		for x := 0; x < colCount; x++ {
			position := fyne.NewPos(float32(x*cellSize), float32(y*cellSize))
			cell := NewCell(position)
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
	m.cont.Resize(fyne.NewSize(colCount*cellSize, rowCount*cellSize))
	(*w).SetContent(m.cont)
	(*w).Resize(fyne.NewSize(colCount*cellSize+5, rowCount*cellSize+5))
}

func (m *Maze) DrawLine(source, target Location, undo bool) {
	lineColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	if undo {
		lineColor = color.RGBA{R: 100, G: 100, B: 100, A: 255}
	}
	line := canvas.NewLine(lineColor)
	line.StrokeWidth = cellSize / 8
	line.Position1 = m.cells[source].Center()
	line.Position2 = m.cells[target].Center()
	m.cont.Add(line)
}

func (m *Maze) HideWall(location Location, side Side) {
	sourceCell := m.cells[location]
	switch side {
	case top:
		neighborX := location.X
		neighborY := location.Y - 1
		if neighborY < 0 {
			panic(fmt.Sprintf("neighborY position is out of range: %d", neighborY))
		}
		sourceCell.walls[top].Hide()
		m.cells[Location{X: neighborX, Y: neighborY}].walls[bottom].Hide()
	case right:
		neighborX := location.X + 1
		neighborY := location.Y
		if neighborX >= colCount {
			panic(fmt.Sprintf("neighborX position is out of range: %d", neighborX))
		}
		sourceCell.walls[right].Hide()
		m.cells[Location{X: neighborX, Y: neighborY}].walls[left].Hide()
	case bottom:
		neighborX := location.X
		neighborY := location.Y + 1
		if neighborY >= rowCount {
			panic(fmt.Sprintf("neighborY position is out of range: %d", neighborY))
		}
		sourceCell.walls[bottom].Hide()
		m.cells[Location{X: neighborX, Y: neighborY}].walls[top].Hide()
	case left:
		neighborX := location.X - 1
		neighborY := location.Y
		if neighborX < 0 {
			panic(fmt.Sprintf("neighborX position is out of range: %d", neighborX))
		}
		sourceCell.walls[left].Hide()
		m.cells[Location{X: neighborX, Y: neighborY}].walls[right].Hide()
	default:
		panic(fmt.Sprintf("unknown side: %d", side))
	}
}
