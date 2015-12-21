package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fname := "input.txt"
	fmt.Printf("processing [%s]...\n", fname)
	r, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	data, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	floor := 0
	for i, v := range data {
		switch v {
		case '(':
			floor++
		case ')':
			floor--
		}
		if floor == -1 {
			fmt.Printf(">>> %d\n", i+1)
			break
		}
	}
	inc := bytes.Count(data, []byte("("))
	dec := bytes.Count(data, []byte(")"))
	fmt.Printf("floor=+%d -%d => %d\n", inc, dec, inc-dec)
}
