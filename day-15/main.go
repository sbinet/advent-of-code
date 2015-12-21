package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

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
	var ins []Ingredient
	s := bufio.NewScanner(f)
	for s.Scan() {
		txt := s.Text()
		in := parse(txt)
		ins = append(ins, in)
	}
	fmt.Printf("ingredients: %v\n", ins)

	calories := false
	rec := &Recipe{Ins: ins}

	ins, amount, score := rec.best(calories)
	fmt.Printf("amount: %v\nscore: %v\n",
		amount,
		score,
	)

	calories = true
	ins, amount, score = rec.best(calories)
	fmt.Printf("amount: %v\nscore: %v\n",
		amount,
		score,
	)
}

func combs(n, m int, emit func([]int)) {
	s := make([]int, m)
	last := m - 1
	var rc func(int, int)
	rc = func(i, next int) {
		for j := next; j < n; j++ {
			s[i] = j
			if i == last {
				emit(s)
			} else {
				rc(i+1, j+1)
			}
		}
		return
	}
	rc(0, 0)
}

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

type Recipe struct {
	Ins []Ingredient
}

func (rec *Recipe) best(calories bool) ([]Ingredient, []int, int) {
	const N = 100
	score := 0
	comb := make([]int, len(rec.Ins))
	combs(N, len(rec.Ins), func(spoons []int) {
		if rec.amount(spoons) != N {
			return
		}
		recipe := make([]int, len(spoons))
		copy(recipe, spoons)
		perms := make([][]int, 0, len(spoons))
		sort.Ints(recipe)
		for perm(sort.IntSlice(recipe)) {
			iter := make([]int, len(recipe))
			copy(iter, recipe)
			perms = append(perms, iter)
		}
		sort.Ints(recipe)
		perms = append(perms, recipe)
		//fmt.Printf("#perms: %d\n", len(perms))
		for _, perm := range perms {
			v, cals := rec.score(perm)
			//fmt.Printf("%v: %v\n", perm, v)
			if v > score {
				if (calories && cals == 500) || !calories {
					score = v
					copy(comb, perm)
					//fmt.Printf("%v: %v <<==\n", perm, v)
				}
			}
		}
	})
	return rec.Ins, comb, score
}

func compr(spoons []int) []int {
	h := make(map[int]int)
	for _, i := range spoons {
		h[i]++
	}
	out := make([]int, len(h))
	for k, v := range h {
		out[k] = v
	}
	return out
}

func (rec *Recipe) amount(spoons []int) int {
	sum := 0
	for _, n := range spoons {
		sum += n
	}
	return sum
}

func (rec *Recipe) score(spoons []int) (int, int) {
	var prop Ingredient
	for i, n := range spoons {
		v := rec.Ins[i]
		prop.Capacity += n * v.Capacity
		prop.Durability += n * v.Durability
		prop.Flavor += n * v.Flavor
		prop.Texture += n * v.Texture
		prop.Calories += n * v.Calories
	}
	if prop.Capacity <= 0 {
		prop.Capacity = 0
	}
	if prop.Durability <= 0 {
		prop.Durability = 0
	}
	if prop.Flavor <= 0 {
		prop.Flavor = 0
	}
	if prop.Texture <= 0 {
		prop.Texture = 0
	}
	return prop.Capacity * prop.Durability * prop.Flavor * prop.Texture, prop.Calories
}

type Ingredient struct {
	Name       string
	Capacity   int
	Durability int
	Flavor     int
	Texture    int
	Calories   int
}

func parse(txt string) Ingredient {
	var in Ingredient
	_, err := fmt.Sscanf(
		txt,
		"%s capacity %d, durability %d, flavor %d, texture %d, calories %d",
		&in.Name, &in.Capacity, &in.Durability, &in.Flavor, &in.Texture,
		&in.Calories,
	)
	if err != nil {
		panic(err)
	}
	in.Name = string(in.Name[:len(in.Name)-1]) // drop ':'
	return in
}
