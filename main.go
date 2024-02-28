package main

import (
	"bytes"
	"fmt"
	rand "math/rand/v2"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Direction int

// set direction enum
const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

// set wall width and height
const (
	HEIGHT = 20
	WIDTH  = 50
)

type Model struct {
	Map   [HEIGHT][WIDTH]string
	Food  food
	Snake snake
}

// 蛇
type snake struct {
	head  [2]int
	body  [][2]int
	speed int
}

// 食物
type food struct {
	position [2]int
	eated    bool
}

// initFood initialization food.
func (m *Model) initFood() {
	m.Food = food{
		position: [2]int{7, 8},
		eated:    false,
	}
}

// addFood add food to map.
func (m *Model) addFood() {
	m.Map[m.Food.position[0]][m.Food.position[1]] = "i"
}

// randomFood Randomly generate a food.
func (m *Model) randomFood() {
	if !m.Food.eated {
		return
	}
	m.Food.position[1], m.Food.position[0] = rand.IntN(WIDTH-1)+1, rand.IntN(HEIGHT-1)+1
	if m.Snake.head == m.Food.position {
		m.randomFood()
		return
	}
	for _, val := range m.Snake.body {
		if val == m.Food.position {
			m.randomFood()
			return
		}
	}
}

// initSnake initialization snake.
func (m *Model) initSnake() {
	m.Snake = snake{
		head:  [2]int{1, 1},
		body:  [][2]int{{2, 1}, {3, 1}}, // 蛇身
		speed: 600,
	}
}
func (m *Model) addSnake() {
	m.Map[m.Snake.head[0]][m.Snake.head[1]] = "O"
	for _, val := range m.Snake.body {
		m.Map[val[0]][val[1]] = "o"
	}

}

func initialModel() *Model {
	m := &Model{}
	for i := range HEIGHT {
		for j := range WIDTH {
			if i == 0 || i == HEIGHT-1 || j == 0 || j == WIDTH-1 {
				m.Map[i][j] = "#"
			} else {
				m.Map[i][j] = " "
			}
		}
	}
	m.initFood()
	m.initSnake()
	m.addFood()
	m.addSnake()
	return m
}
func (m *Model) moveSnake(dir Direction) {
	var eated = false
	beforHead := m.Snake.head
	switch dir {
	case UP:
		m.Snake.head[0]--
	case DOWN:
		m.Snake.head[0]++
	case LEFT:
		m.Snake.head[1]--
	case RIGHT:
		m.Snake.head[1]++
	}
	eated = (m.Snake.head == m.Food.position)
	m.Snake.body = append([][2]int{beforHead}, m.Snake.body...)
	m.Map[m.Snake.head[0]][m.Snake.head[1]] = "O"
	m.Map[beforHead[0]][beforHead[1]] = "o"
	if !eated {
		tail := m.Snake.body[len(m.Snake.body)-1]
		m.Map[tail[0]][tail[1]] = " "
		m.Snake.body = m.Snake.body[:len(m.Snake.body)-1]
		return
	}
	m.Food.eated = eated
	m.randomFood()
	m.addFood()
}

var currentDir = RIGHT

func (m *Model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var dir Direction = currentDir
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "w":
			dir = UP
		case "s":
			dir = DOWN
		case "a":
			dir = LEFT
		case "d":
			dir = RIGHT
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	}
	currentDir = dir
	m.moveSnake(dir)
	return m, nil
}
func (m *Model) View() string {
	var buf = bytes.NewBufferString("")
	for _, row := range m.Map {
		buf.WriteString(fmt.Sprintf("%s\n", strings.Join(row[:], "")))
	}
	return buf.String()
}
func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		os.Exit(1)
	}
}
