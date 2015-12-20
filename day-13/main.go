package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Edge struct {
	N1, N2 string
}

type Graph struct {
	vertices map[string]bool
	edges    map[Edge]int
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	gr := Graph{
		vertices: make(map[string]bool),
		edges:    make(map[Edge]int),
	}
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		var n1, n2, sign string
		var n int
		txt := scan.Text()
		txt = txt[:len(txt)-1] // drop the .
		toks := strings.Split(txt, " ")
		n1 = toks[0]
		sign = toks[2]
		n, err = strconv.Atoi(toks[3])
		if err != nil {
			panic(err)
		}
		n2 = toks[10]
		switch sign {
		case "gain":
			n = +n
		case "lose":
			n = -n
		}
		gr.add(n1, n2, n)
	}
	err = scan.Err()
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		panic(err)
	}

	table, hap := gr.best()
	fmt.Printf("table: %v\n", table)
	fmt.Printf("happiness: %d\n", hap)

	for _, v := range table {
		gr.add(v, "me", 0)
	}

	table, hap = gr.best()
	fmt.Printf("table: %v\n", table)
	fmt.Printf("happiness: %d\n", hap)

}

func (gr *Graph) add(n1, n2 string, n int) {
	edge := Edge{
		N1: n1,
		N2: n2,
	}
	//fmt.Printf("--> %v\n", edge)
	_, dup := gr.edges[edge]
	if dup {
		panic("duplicate")
	}
	gr.edges[edge] = n
	gr.vertices[n1] = true
	gr.vertices[n2] = true
}

func (g *Graph) best() ([]string, int) {
	table := make([]string, 0, len(g.vertices))
	for k := range g.vertices {
		table = append(table, k)
	}
	sort.Strings(table)
	perms := make([][]string, 0)
	for perm(sort.StringSlice(table)) {
		iter := make([]string, len(table))
		copy(iter, table)
		perms = append(perms, iter)
	}
	iperm := -1
	weight := 0
	for i, perm := range perms {
		w := g.weight(perm)
		// fmt.Printf("table: %v => %d\n", perm, w)
		if w > weight || iperm < 0 {
			weight = w
			iperm = i
		}
	}
	return perms[iperm], weight
}

func (g *Graph) weight(table []string) int {
	w := 0
	for i := range table {
		var nl, n0, nr string
		if i == 0 {
			nl = table[len(table)-1]
		} else {
			nl = table[i-1]
		}
		n0 = table[i]
		if i == len(table)-1 {
			nr = table[0]
		} else {
			nr = table[i+1]
		}
		w += g.get(n0, nl)
		w += g.get(n0, nr)
	}
	return w
}

func (g *Graph) get(n1, n2 string) int {
	if v, ok := g.edges[Edge{n1, n2}]; ok {
		return v
	}
	if v, ok := g.edges[Edge{n2, n1}]; ok {
		return v
	}
	panic("unreachable")
}

// Generate the next permutation of data if possible and return true.
// Return false if there is no more permutation left.
// Based on the algorithm described here:
// http://en.wikipedia.org/wiki/Permutation#Generation_in_lexicographic_order
func perm(data sort.Interface) bool {
	var k, l int
	for k = data.Len() - 2; ; k-- { // 1.
		if k < 0 {
			return false
		}

		if data.Less(k, k+1) {
			break
		}
	}
	for l = data.Len() - 1; !data.Less(k, l); l-- { // 2.
	}
	data.Swap(k, l)                             // 3.
	for i, j := k+1, data.Len()-1; i < j; i++ { // 4.
		data.Swap(i, j)
		j--
	}
	return true
}
