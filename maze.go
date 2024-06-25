package main

import (
	"errors"
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
	leftWall   *canvas.Line
	rightWall  *canvas.Line
	topWall    *canvas.Line
	bottomWall *canvas.Line
	position   fyne.Position
}

func NewCell(position fyne.Position) Cell {
	cell := Cell{
		leftWall:   canvas.NewLine(color.White),
		rightWall:  canvas.NewLine(color.White),
		topWall:    canvas.NewLine(color.White),
		bottomWall: canvas.NewLine(color.White),
		position:   position,
	}
	cell.leftWall.Position1 = fyne.NewPos(position.X, position.Y)
	cell.leftWall.Position2 = fyne.NewPos(position.X, position.Y+cellSize)
	cell.rightWall.Position1 = fyne.NewPos(position.X+cellSize, position.Y)
	cell.rightWall.Position2 = fyne.NewPos(position.X+cellSize, position.Y+cellSize)
	cell.topWall.Position1 = fyne.NewPos(position.X, position.Y)
	cell.topWall.Position2 = fyne.NewPos(position.X+cellSize, position.Y)
	cell.bottomWall.Position1 = fyne.NewPos(position.X, position.Y+cellSize)
	cell.bottomWall.Position2 = fyne.NewPos(position.X+cellSize, position.Y+cellSize)
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
		m.cont.Add(cell.leftWall)
		m.cont.Add(cell.rightWall)
		m.cont.Add(cell.topWall)
		m.cont.Add(cell.bottomWall)
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

func (m *Maze) HideWall(location Location, side Side) error {
	sourceCell := m.cells[location]
	switch side {
	case top:
		neighborX := location.X
		neighborY := location.Y - 1
		if neighborY < 0 {
			return errors.New(fmt.Sprintf("neighborY position is out of range: %d", neighborY))
		}
		sourceCell.topWall.Hide()
		m.cells[Location{X: neighborX, Y: neighborY}].bottomWall.Hide()
	case right:
		neighborX := location.X + 1
		neighborY := location.Y
		if neighborX >= colCount {
			return errors.New(fmt.Sprintf("neighborX position is out of range: %d", neighborX))
		}
		sourceCell.rightWall.Hide()
		m.cells[Location{X: neighborX, Y: neighborY}].leftWall.Hide()
	case bottom:
		neighborX := location.X
		neighborY := location.Y + 1
		if neighborY >= rowCount {
			return errors.New(fmt.Sprintf("neighborY position is out of range: %d", neighborY))
		}
		sourceCell.bottomWall.Hide()
		m.cells[Location{X: neighborX, Y: neighborY}].topWall.Hide()
	case left:
		neighborX := location.X - 1
		neighborY := location.Y
		if neighborX < 0 {
			return errors.New(fmt.Sprintf("neighborX position is out of range: %d", neighborX))
		}
		sourceCell.leftWall.Hide()
		m.cells[Location{X: neighborX, Y: neighborY}].rightWall.Hide()
	default:
		return errors.New(fmt.Sprintf("unknown side: %d", side))
	}
	return nil
}
