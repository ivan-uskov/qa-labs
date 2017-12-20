package main

import (
	"os"
	"fmt"
	"math"
	"strconv"
)

type Triangle struct {
	a float64
	b float64
	c float64
}

const floatEpsilon = 0.00000001

func floatEquals(a, b float64) bool {
	return math.Abs(a-b) < floatEpsilon
}

const (
	_ = iota
	regular
	isosceles
	equilateral
	invalid
)

type TriangleType int

func checkTriangleInvalid(tr *Triangle) bool {
	checkInvalidSides := func(sum float64, another float64) bool {
		return sum < another || floatEquals(sum, another)
	}

	return checkInvalidSides(tr.a+tr.b, tr.c) ||
		checkInvalidSides(tr.a+tr.c, tr.b) ||
		checkInvalidSides(tr.b+tr.c, tr.a)
}

func checkEquilateral(tr *Triangle) bool {
	return floatEquals(tr.a, tr.b) && floatEquals(tr.b, tr.c)
}

func checkIsosceles(tr *Triangle) bool {
	return floatEquals(tr.a, tr.b) || floatEquals(tr.b, tr.c) || floatEquals(tr.a, tr.c)
}

func detectTriangleType(tr *Triangle) TriangleType {
	if checkTriangleInvalid(tr) {
		return invalid
	}
	if checkEquilateral(tr) {
		return equilateral
	}
	if checkIsosceles(tr) {
		return isosceles
	}

	return regular
}

func parseTriangleSide(side string) (float64, error) {
	num, err := strconv.ParseFloat(side, 64)
	if err != nil {
		return 0, fmt.Errorf(side + " не число")
	}

	if num < 0 || floatEquals(num, 0) {
		return 0, fmt.Errorf("значение должно быть больше 0")
	}

	return num, nil
}

func parseTriangle(args []string) (*Triangle, error) {
	a, err := parseTriangleSide(args[0])
	if err != nil {
		return nil, fmt.Errorf("неправильный параметр <a>: %s", err.Error())
	}

	b, err := parseTriangleSide(args[1])
	if err != nil {
		return nil, fmt.Errorf("неправильный параметр <b>: %s", err.Error())
	}

	c, err := parseTriangleSide(args[2])
	if err != nil {
		return nil, fmt.Errorf("неправильный параметр <c>: %s", err.Error())
	}

	return &Triangle{a, b, c}, nil
}

func (t TriangleType) ToString() string {
	switch t {
	case regular:
		return "Обычный"
	case isosceles:
		return "Равнобедренный"
	case equilateral:
		return "Равносторонний"
	default:
		return "Не треугольник"
	}
}

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Println("Ожидается 3 аргумента <a> <b> <c>")
		os.Exit(1)
	}

	tr, err := parseTriangle(args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	trType := detectTriangleType(tr)
	fmt.Println(trType.ToString())
}
