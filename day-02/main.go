package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

func main() {
	fname := "input.txt"
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	total := 0
	ribbon := 0
	for {
		var b box
		_, err = fmt.Fscanf(f, "%dx%dx%d\n",
			&b[0], &b[1], &b[2],
		)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		total += b.surface()
		ribbon += b.ribbon()
	}
	fmt.Printf("total=%d\n", total)
	fmt.Printf("ribbon= %d\n", ribbon)
}

type box [3]int

func (b box) surface() int {
	slack := b.slack()
	return 2*(b[0]*b[1]+b[1]*b[2]+b[2]*b[0]) + slack
}

func (b box) slack() int {
	sort.Ints(b[:])
	return b[0] * b[1]
}

func (b box) ribbon() int {
	sort.Ints(b[:])
	wrap := 2 * (b[0] + b[1])
	ribbon := b[0] * b[1] * b[2]
	return wrap + ribbon
}
