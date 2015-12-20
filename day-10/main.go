package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	for _, s := range []string{
		"1",
		"11",
		"21",
		"1211",
		"111221",
	} {
		fmt.Printf("%s -> %s\n", s, next(s))
	}

	s := "1113122113"
	for i := 0; i < 50; i++ {
		//fmt.Printf("=> %s\n", s)
		s = next(s)
	}
	fmt.Printf("=> %d\n", len(s))
}

func next(s string) string {
	var o []string
	for i := 0; i < len(s); i++ {
		v := s[i]
		n := 0
		for j := i; j < len(s) && s[j] == v; j++ {
			n++
		}
		o = append(o, strconv.Itoa(n), string(v))
		i += n - 1
	}
	return strings.Join(o, "")
}
