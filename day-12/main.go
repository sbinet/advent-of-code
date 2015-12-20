package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type State struct {
	sum int
	obj bool
	red bool
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	stack := make([]State, 1)
	ptr := &stack[0]

	var tok json.Token
	dec := json.NewDecoder(f)
	for {
		tok, err = dec.Token()
		if err != nil {
			break
		}
		switch v := tok.(type) {
		case json.Delim:
			switch v {
			case '[':
				stack = append(stack,
					State{
						sum: 0,
						obj: false,
						red: false,
					},
				)
			case '{':
				stack = append(stack,
					State{
						sum: 0,
						obj: true,
						red: false,
					},
				)
			case ']', '}':
				old := ptr
				stack = stack[:len(stack)-1]
				ptr = &stack[len(stack)-1]
				if (old.obj && !old.red) || !old.obj {
					ptr.sum += old.sum
				}
			}
			ptr = &stack[len(stack)-1]
		case float64:
			ptr.sum += int(v)
		case json.Number:
			var ii int64
			ii, err = v.Int64()
			if err != nil {
				panic(err)
			}
			ptr.sum += int(ii)
		case string:
			if v == "red" {
				ptr.red = true
			}
		}
	}
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("sum= %d\n", ptr.sum)
}
