package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"math/rand/v2"
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
	position fyne.Position
	size     float32
	visited  bool
}

func NewCell(location Location, size float32) Cell {
	cell := Cell{
		walls:    make(map[Side]*canvas.Line, 4),
		location: location,
		position: fyne.NewPos(float32(location.X)*size, float32(location.Y)*size),
		size:     size,
	}
	cell.walls[top] = canvas.NewLine(color.White)
	cell.walls[top].Position1 = fyne.NewPos(cell.position.X, cell.position.Y)
	cell.walls[top].Position2 = fyne.NewPos(cell.position.X+size, cell.position.Y)
	cell.walls[right] = canvas.NewLine(color.White)
	cell.walls[right].Position1 = fyne.NewPos(cell.position.X+size, cell.position.Y)
	cell.walls[right].Position2 = fyne.NewPos(cell.position.X+size, cell.position.Y+size)
	cell.walls[bottom] = canvas.NewLine(color.White)
	cell.walls[bottom].Position1 = fyne.NewPos(cell.position.X, cell.position.Y+size)
	cell.walls[bottom].Position2 = fyne.NewPos(cell.position.X+size, cell.position.Y+size)
	cell.walls[left] = canvas.NewLine(color.White)
	cell.walls[left].Position1 = fyne.NewPos(cell.position.X, cell.position.Y)
	cell.walls[left].Position2 = fyne.NewPos(cell.position.X, cell.position.Y+size)
	return cell
}

type Maze struct {
	colCount       int
	rowCount       int
	cells          map[Location]*Cell
	cont           *fyne.Container
	cellSize       float32
	generator      *rand.Rand
	animationDelay float64
}

func NewMaze(colCount, rowCount int, cellSize float32, seed rand.Source) Maze {
	m := Maze{
		colCount:  colCount,
		rowCount:  rowCount,
		cells:     make(map[Location]*Cell, colCount*rowCount),
		cont:      container.NewWithoutLayout(),
		cellSize:  cellSize,
		generator: rand.New(seed),
	}
	for y := 0; y < rowCount; y++ {
		for x := 0; x < colCount; x++ {
			cell := NewCell(Location{X: x, Y: y}, cellSize)
			m.cells[Location{X: x, Y: y}] = &cell
		}
	}
	// create entrance and exit
	m.hideWall(Location{X: 0, Y: 0}, top)
	m.hideWall(Location{X: colCount - 1, Y: rowCount - 1}, bottom)
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

func (m *Maze) drawLine(source, target Location, undo bool) {
	lineColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}
	if undo {
		lineColor = color.RGBA{R: 100, G: 100, B: 100, A: 255}
	}
	line := canvas.NewLine(lineColor)
	line.StrokeWidth = m.cellSize / 8
	line.Position1 = fyne.NewPos((float32(source.X)*m.cellSize+float32(source.X)*m.cellSize+m.cellSize)/2, (float32(source.Y)*m.cellSize+float32(source.Y)*m.cellSize+m.cellSize)/2)
	line.Position2 = fyne.NewPos((float32(target.X)*m.cellSize+float32(target.X)*m.cellSize+m.cellSize)/2, (float32(target.Y)*m.cellSize+float32(target.Y)*m.cellSize+m.cellSize)/2)
	m.cont.Add(line)
	m.cont.Refresh()
}

// hideWall will make the specified location's wall invisible
// If the location has a neighbor, its corresponding wall will be hidden as well.
// For example, if the left wall of location 1,1 is hidden,
// the right wall of location 0,1 must also be hidden to not render that line.
// If the location's wall is on an outer edge, only the single wall will be hidden.
func (m *Maze) hideWall(location Location, side Side) {
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

func (m *Maze) removeWallsR(location Location, from Side) {
	if m.animationDelay != 0 && location.X == 0 && location.Y == 0 {
		time.Sleep(10 * time.Second)
	}
	if cell, ok := m.cells[location]; !ok {
		panic(fmt.Sprintf("location %v does not exist", location))
	} else {
		time.Sleep(time.Duration(m.animationDelay * float64(time.Second)))
		m.hideWall(location, from)
		cell.visited = true
	}
	for {
		var nextLocations []Location
		if location.X > 0 && !m.cells[Location{X: location.X - 1, Y: location.Y}].visited {
			nextLocations = append(nextLocations, Location{X: location.X - 1, Y: location.Y})
		}
		if location.X < m.colCount-1 && !m.cells[Location{X: location.X + 1, Y: location.Y}].visited {
			nextLocations = append(nextLocations, Location{X: location.X + 1, Y: location.Y})
		}
		if location.Y > 0 && !m.cells[Location{X: location.X, Y: location.Y - 1}].visited {
			nextLocations = append(nextLocations, Location{X: location.X, Y: location.Y - 1})
		}
		if location.Y < m.rowCount-1 && !m.cells[Location{X: location.X, Y: location.Y + 1}].visited {
			nextLocations = append(nextLocations, Location{X: location.X, Y: location.Y + 1})
		}
		if len(nextLocations) == 0 {
			return
		}
		locationIndex := m.generator.Int() % len(nextLocations)
		nextLocation := nextLocations[locationIndex]
		if nextLocation.X > location.X {
			m.removeWallsR(nextLocation, left)
		}
		if nextLocation.X < location.X {
			m.removeWallsR(nextLocation, right)
		}
		if nextLocation.Y > location.Y {
			m.removeWallsR(nextLocation, top)
		}
		if nextLocation.Y < location.Y {
			m.removeWallsR(nextLocation, bottom)
		}
	}
}

func (m *Maze) resetVisited() {
	for _, cell := range m.cells {
		cell.visited = false
	}
}

func (m *Maze) solveR(location Location) bool {
	cell := m.cells[location]
	cell.visited = true
	time.Sleep(time.Duration(m.animationDelay * float64(time.Second)))
	if location.X == m.colCount-1 && location.Y == m.rowCount-1 {
		return true
	}
	targetLocation := Location{X: location.X, Y: location.Y - 1}
	if location.Y > 0 && !m.cells[targetLocation].visited && !m.cells[location].walls[top].Visible() {
		m.drawLine(location, targetLocation, false)
		winner := m.solveR(targetLocation)
		if winner {
			return true
		}
		m.drawLine(location, targetLocation, true)
	}
	targetLocation = Location{X: location.X + 1, Y: location.Y}
	if location.X < m.colCount-1 && !m.cells[targetLocation].visited && !m.cells[location].walls[right].Visible() {
		m.drawLine(location, targetLocation, false)
		winner := m.solveR(targetLocation)
		if winner {
			return true
		}
		m.drawLine(location, targetLocation, true)
	}
	targetLocation = Location{X: location.X, Y: location.Y + 1}
	if location.Y < m.rowCount-1 && !m.cells[targetLocation].visited && !m.cells[location].walls[bottom].Visible() {
		m.drawLine(location, targetLocation, false)
		winner := m.solveR(targetLocation)
		if winner {
			return true
		}
		m.drawLine(location, targetLocation, true)
	}
	targetLocation = Location{X: location.X - 1, Y: location.Y}
	if location.X > 0 && !m.cells[targetLocation].visited && !m.cells[location].walls[left].Visible() {
		m.drawLine(location, targetLocation, false)
		winner := m.solveR(targetLocation)
		if winner {
			return true
		}
		m.drawLine(location, targetLocation, true)
	}
	return false
}

func (m *Maze) solve() bool {
	return m.solveR(Location{X: 0, Y: 0})
}

func (m *Maze) Animate() {
	m.removeWallsR(Location{X: 0, Y: 0}, top)
	m.resetVisited()
	m.drawLine(Location{X: 0, Y: -1}, Location{X: 0, Y: 0}, false)
	m.solve()
	m.drawLine(Location{X: m.colCount - 1, Y: m.rowCount - 1}, Location{X: m.colCount - 1, Y: m.rowCount}, false)
}
