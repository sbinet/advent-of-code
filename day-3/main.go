package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	data, err := ioutil.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}

	var p1, p2 pos
	houses := make(map[pos]int)
	houses[p1]++
	for i, v := range data {
		var p *pos
		if i%2 == 0 {
			p = &p1
		} else {
			p = &p2
		}
		switch v {
		case '>':
			p.x++
		case '<':
			p.x--
		case '^':
			p.y++
		case 'v':
			p.y--
		}
		houses[*p]++
	}
	fmt.Printf("houses: %d\n", len(houses))
}

type pos struct {
	x, y int
}
