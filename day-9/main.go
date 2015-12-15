package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

type graph struct {
	vertices set
	edges    map[edge]int
}

func newGraph() graph {
	return graph{
		vertices: make(set),
		edges:    make(map[edge]int),
	}
}

func (g *graph) add(v1, v2 string, dist int) {
	g.vertices.add(v1)
	g.vertices.add(v2)
	g.edges[edge{v1, v2}] = dist
}

func (g *graph) shortestDist() ([]string, int) {
	cities := make([]string, 0, len(g.vertices))
	for k := range g.vertices {
		cities = append(cities, k)
	}
	sort.Strings(cities)
	perms := make([][]string, 0)
	for perm(sort.StringSlice(cities)) {
		iter := make([]string, len(cities))
		copy(iter, cities)
		perms = append(perms, iter)
	}
	iperm := -1
	dist := -1
	for i, perm := range perms {
		d := g.dist(perm)
		if d < dist || iperm < 0 {
			dist = d
			iperm = i
		}
	}
	return perms[iperm], dist
}

func (g *graph) longuestDist() ([]string, int) {
	cities := make([]string, 0, len(g.vertices))
	for k := range g.vertices {
		cities = append(cities, k)
	}
	sort.Strings(cities)
	perms := make([][]string, 0)
	for perm(sort.StringSlice(cities)) {
		iter := make([]string, len(cities))
		copy(iter, cities)
		perms = append(perms, iter)
	}
	iperm := -1
	dist := 0
	for i, perm := range perms {
		d := g.dist(perm)
		if d > dist || iperm < 0 {
			dist = d
			iperm = i
		}
	}
	return perms[iperm], dist
}

func (g *graph) dist(path []string) int {
	dist := 0
	for i := 1; i < len(path); i++ {
		v1 := path[i-1]
		v2 := path[i]
		e12 := edge{v1, v2}
		e21 := edge{v2, v1}
		if d, ok := g.edges[e12]; ok {
			dist += d
			continue
		}
		if d, ok := g.edges[e21]; ok {
			dist += d
			continue
		}
		log.Panicf("no such edge: %#v\n", e12)
	}
	return dist
}

type edge struct {
	v1 string
	v2 string
}

func main() {
	fname := "input.txt"
	if len(os.Args) > 1 {
		fname = os.Args[1]
	}
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	g := newGraph()
	s := bufio.NewScanner(f)
	for s.Scan() {
		txt := s.Text()
		var v1, v2 string
		var dist = 0
		_, err = fmt.Sscanf(txt, "%s to %s = %d",
			&v1, &v2, &dist,
		)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s -> %s = %d\n", v1, v2, dist)
		g.add(v1, v2, dist)
	}
	err = s.Err()
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		panic(err)
	}

	pmin, dmin := g.shortestDist()
	pmax, dmax := g.longuestDist()

	fmt.Printf("%s => %d\n", strings.Join(pmin, " -> "), dmin)
	fmt.Printf("%s => %d\n", strings.Join(pmax, " -> "), dmax)
}

type set map[string]bool

func (set *set) add(v string) {
	(*set)[v] = true
}
func (set *set) addall(vs []string) {
	for _, v := range vs {
		set.add(v)
	}
}
func (set *set) has(v string) bool {
	_, ok := (*set)[v]
	return ok
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
