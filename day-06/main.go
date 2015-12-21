package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	grid := [1000][1000]int{}
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		var x1, y1, x2, y2 int
		txt := s.Text()
		var n string
		var cmd func(int) int
		switch {
		case strings.HasPrefix(txt, "turn on"):
			n = "turn on"
			cmd = turnon
		case strings.HasPrefix(txt, "turn off"):
			n = "turn off"
			cmd = turnoff
		case strings.HasPrefix(txt, "toggle"):
			n = "toggle"
			cmd = toggle
		}
		txt = txt[len(n)+1:]

		_, err = fmt.Sscanf(txt, "%d,%d through %d,%d",
			&x1, &y1, &x2, &y2,
		)
		if err != nil {
			fmt.Printf("txt=%q\n", txt)
			panic(err)
		}
		for i := x1; i <= x2; i++ {
			for j := y1; j <= y2; j++ {
				grid[i][j] = cmd(grid[i][j])
			}
		}
	}

	total := 0
	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			total += grid[i][j]
		}
	}
	fmt.Printf("total= %d\n", total)
}

func turnon(v int) int {
	return v + 1
}

func turnoff(v int) int {
	if v == 0 {
		return 0
	}
	return v - 1
}

func toggle(v int) int {
	return v + 2

	switch v {
	case 0:
		return 1
	case 1:
		return 0
	}
	panic("unreachable")
}
