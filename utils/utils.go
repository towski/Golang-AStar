package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"math"
	"github.com/towski/Golang-AStar/term"
)

type Point struct {
	X, Y    int
	H, G, F int
	Parent  *Point
}

func (p Point) String() string {
	return "[" + strconv.Itoa(p.X) + ", " + strconv.Itoa(p.Y) + ", " + strconv.Itoa(p.F) + "]"
}

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func Clear() {
	fmt.Printf("\033[100B")
	for i := 0; i < 100; i++ {
		fmt.Printf("\033[1A")
		fmt.Printf("\033[K")
	}
}

func GetRandInt(limit int) int {
	return r.Intn(limit)
}

var origin, dest Point
var openList, closeList, path []Point
var FinalPoint Point
var Result int

// Set the origin point
func setOrig(s *Scene) {
	origin = Point{GetRandInt(s.rows-2) + 1, GetRandInt(s.cols-2) + 1, 0, 0, 0, nil}
	if s.Data[origin.X][origin.Y] == ' ' {
		s.Data[origin.X][origin.Y] = 'A'
	} else {
		setOrig(s)
	}
}

// Set the destination point
func setDest(s *Scene) {
	dest = Point{GetRandInt(s.rows-2) + 1, GetRandInt(s.cols-2) + 1, 0, 0, 0, nil}

	if s.Data[dest.X][dest.Y] == ' ' {
		s.Data[dest.X][dest.Y] = 'B'
	} else {
		setDest(s)
	}
}

// Init origin, destination. Put the origin point into the openlist by the way
func InitAstar(s *Scene) {
    Result = 10
	//setOrig(s)
	//setDest(s)
	openList = append(openList, origin)
}

func FindPath(s *Scene) {
	current := getFMin()
	addToCloseList(current, s)
	walkable := getWalkable(current, s)
	for _, p := range walkable {
		addToOpenList(p)
	}
}

func getFMin() Point {
	if len(openList) == 0 {
		fmt.Println("No way!!!")
        Result = -1
	}
	index := 0
	for i, p := range openList {
		if (i > 0) && (p.F <= openList[index].F) {
			index = i
		}
	}
	return openList[index]
}

func getWalkable(p Point, s *Scene) []Point {
	var around []Point
	row, col := p.X, p.Y
	left := s.Data[row][col-1]
	up := s.Data[row-1][col]
	right := s.Data[row][col+1]
	down := s.Data[row+1][col]
	leftup := s.Data[row-1][col-1]
	rightup := s.Data[row-1][col+1]
	leftdown := s.Data[row+1][col-1]
	rightdown := s.Data[row+1][col+1]
	if (left == ' ') || (left == 'B') {
		around = append(around, Point{row, col - 1, 0, 0, 0, &p})
	}
	if (leftup == ' ') || (leftup == 'B') {
		around = append(around, Point{row - 1, col - 1, 0, 0, 0, &p})
	}
	if (up == ' ') || (up == 'B') {
		around = append(around, Point{row - 1, col, 0, 0, 0, &p})
	}
	if (rightup == ' ') || (rightup == 'B') {
		around = append(around, Point{row - 1, col + 1, 0, 0, 0, &p})
	}
	if (right == ' ') || (right == 'B') {
		around = append(around, Point{row, col + 1, 0, 0, 0, &p})
	}
	if (rightdown == ' ') || (rightdown == 'B') {
		around = append(around, Point{row + 1, col + 1, 0, 0, 0, &p})
	}
	if (down == ' ') || (down == 'B') {
		around = append(around, Point{row + 1, col, 0, 0, 0, &p})
	}
	if (leftdown == ' ') || (leftdown == 'B') {
		around = append(around, Point{row + 1, col - 1, 0, 0, 0, &p})
	}
	return around
}

func addToOpenList(p Point) {
	updateWeight(&p)
	if checkExist(p, closeList) {
		return
	}
	if !checkExist(p, openList) {
		openList = append(openList, p)
	} else {
		if openList[findPoint(p, openList)].F > p.F { //New path found
			openList[findPoint(p, openList)].Parent = p.Parent
		}
	}
}

// Update G, H, F of the point
func updateWeight(p *Point) {
	if checkRelativePos(*p) == 1 {
		p.G = p.Parent.G + 10
	} else {
		p.G = p.Parent.G + 14
	}
	absx := (int)(math.Abs((float64)(dest.X - p.X)))
	absy := (int)(math.Abs((float64)(dest.Y - p.Y)))
	p.H = (absx + absy) * 30
	p.F = p.G + p.H
}

func removeFromOpenList(p Point) {
	index := findPoint(p, openList)
	if index == -1 {
        Result = 0
	}
	openList = append(openList[:index], openList[index+1:]...)
}

func addToCloseList(p Point, s *Scene) {
	removeFromOpenList(p)
	if (p.X == dest.X) && (p.Y == dest.Y) {
		//generatePath(p, s)
        FinalPoint = p
        // don't draw at the end
		//s.Draw()
        Result = 1
	}
	// if (p.Parent != nil) && (checkRelativePos(p) == 2) {
	// 	parent := p.Parent
	// 	//rdblck := s.Data[p.X][parent.Y] | s.Data[parent.X][p.Y]
	// 	//fmt.Printf("%c\n", rdblck)
	// 	if (s.Data[p.X][parent.Y] == '#') || (s.Data[parent.X][p.Y] == '#') {
	// 		return
	// 	}
	// }
	if s.Data[p.X][p.Y] != 'A' {
		s.Data[p.X][p.Y] = '·'
	}
	closeList = append(closeList, p)
}

func checkExist(p Point, arr []Point) bool {
	for _, point := range arr {
		if p.X == point.X && p.Y == point.Y {
			return true
		}
	}
	return false
}

func findPoint(p Point, arr []Point) int {
	for index, point := range arr {
		if p.X == point.X && p.Y == point.Y {
			return index
		}
	}

	return -1
}

func checkRelativePos(p Point) int {
	parent := p.Parent
	hor := (int)(math.Abs((float64)(p.X - parent.X)))
	ver := (int)(math.Abs((float64)(p.Y - parent.Y)))
	return hor + ver
}

func generatePath(p Point, s *Scene) {
	if (s.Data[p.X][p.Y] != 'A') && (s.Data[p.X][p.Y] != 'B') {
		s.Data[p.X][p.Y] = '*'
	}
	if p.Parent != nil {
		generatePath(*(p.Parent), s)
	}
}

type Scene struct {
	rows, cols int
	Data      [][]byte
}

func (s *Scene) InitScene(rows int, cols int) {
	s.rows = rows
	s.cols = cols

	s.Data = make([][]byte, s.rows)
	for i := 0; i < s.rows; i++ {
		s.Data[i] = make([]byte, s.cols)
		for j := 0; j < s.cols; j++ {
			if i == 0 || i == s.rows-1 || j == 0 || j == s.cols-1 {
				s.Data[i][j] = '#'
			} else {
				s.Data[i][j] = ' '
			}
		}
	}
}

func (s *Scene) Draw() {
	for i := 0; i < s.rows; i++ {
		for j := 0; j < s.cols; j++ {
			var color string
			switch s.Data[i][j] {
			case '#':
				color = term.FgCyan
			case 'A':
				color = term.FgRed
			case 'B':
				color = term.FgBlue
			case '*':
				color = term.FgYellow
				// case ' ':
				// 	if checkExist(Point{i, j, 0, 0, 0, nil}, closeList) {
				// 		fmt.Printf("·")
				// 		continue
				// 	}
			}
			fmt.Printf("%s%c%s", color, s.Data[i][j], term.Reset)
		}
		fmt.Printf("\n")
	}
}

func (s *Scene) AddWalls(num int) {
	for i := 0; i < num; i++ {
		ori := GetRandInt(2)
		length := GetRandInt(16) + 1
		row := GetRandInt(s.rows)
		col := GetRandInt(s.cols)
		switch ori {
		case 0:
			for i := 0; i < length; i++ {
				if col+i >= s.cols {
					break
				}
				s.Data[row][col+i] = '#'
			}

		case 1:
			for i := 0; i < length; i++ {
				if row+i >= s.rows {
					break
				}
				s.Data[row+i][col] = '#'
			}
		}
	}
}
