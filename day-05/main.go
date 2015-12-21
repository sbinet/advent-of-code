package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	total := 0
	s := bufio.NewScanner(f)
	for s.Scan() {
		txt := s.Bytes()
		if nice2(txt) {
			total++
		}
	}
	err = s.Err()
	if err != nil {
		if err != io.EOF {
			panic(err)
		}
	}

	fmt.Printf("nice= %d\n", total)
}

func nice1(data []byte) bool {
	if bytes.Contains(data, []byte("ab")) ||
		bytes.Contains(data, []byte("cd")) ||
		bytes.Contains(data, []byte("pq")) ||
		bytes.Contains(data, []byte("xy")) {
		return false
	}
	vowels := 0
	pairs := 0
	for i, v := range data {
		switch v {
		case 'a', 'e', 'i', 'o', 'u':
			vowels++
		}
		if i > 0 && data[i-1] == v {
			pairs++
		}
	}

	return vowels >= 3 && pairs > 0
}

func nice2(data []byte) bool {
	pairs := 0
	trios := 0
	for i := range data {
		if i < len(data)-1 {
			p := data[i : i+2]
			if bytes.Contains(data[i+2:], p) {
				pairs++
			}
			if i < len(data)-2 && data[i] == data[i+2] {
				trios++
			}
		}
		if pairs > 0 && trios > 0 {
			return true
		}
	}
	return pairs > 0 && trios > 0
}
