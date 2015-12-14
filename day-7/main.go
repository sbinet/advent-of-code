package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	vars = make(map[string]uint16)
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		txt := s.Text()
		fmt.Printf("%q\n", txt)
		toks := strings.Split(txt, " ")
		idx := len(toks) - 2
		if toks[idx] != "->" {
			log.Panicf("err: %v\n", txt)
		}
		switch len(toks[:idx]) {
		case 1:
			// assignment
			v, ok := vars[toks[0]]
			if ok {
				vars[toks[idx+1]] = v
			} else {
				v := atoi(toks[0])
				vars[toks[idx+1]] = v
			}

		case 2:
			if toks[0] != "NOT" {
				log.Panicf("invalid command. got %q. want %q (txt=%q)\n",
					toks[0], "NOT",
					txt,
				)
			}
			v, ok := vars[toks[1]]
			if !ok {
				log.Panicf("no such variable %q\n", toks[1])
			}
			vars[toks[idx+1]] = ^v

		case 3:
			lhs := get(toks, 0)
			op := toks[1]
			rhs := get(toks, 2)
			v := toks[idx+1]
			switch op {
			case "AND":
				vars[v] = lhs & rhs
			case "OR":
				vars[v] = lhs | rhs
			case "LSHIFT":
				vars[v] = lhs << rhs
			case "RSHIFT":
				vars[v] = lhs >> rhs
			default:
				log.Panicf("unknown command %q (txt=%q)\n", op, txt)
			}
		}
	}
	err = s.Err()
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		panic(err)
	}
	for k, v := range vars {
		fmt.Printf("%v: %d\n", k, v)
	}
	fmt.Printf("a=%v\n", vars["a"])
}

func get(toks []string, idx int) uint16 {
	if v, ok := vars[toks[idx]]; ok {
		return v
	}
	return atoi(toks[idx])
}

func atoi(s string) uint16 {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return uint16(v)
}
