package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
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
	nstr := 0
	nmem := 0
	nenc := 0
	s := bufio.NewScanner(f)
	for s.Scan() {
		raw := s.Bytes()
		mem := raw[1 : len(raw)-1]

		n := 0
		for i := 0; i < len(mem); i++ {
			v := mem[i]
			switch {
			case 'a' <= v && v <= 'z':
				n++
			case 'A' <= v && v <= 'Z':
				n++
			case v == 92: // a \
				switch mem[i+1] {
				case 92, 34:
					n++
					i++
					continue
				case 120: // \x
					n++
					i += 3
				}
			default:
				log.Panicf("raw=%s mem[%d]=%d (mem=%v)\n", string(raw), i, v,
					mem)

			}
		}
		nstr += len(raw)
		nmem += n
		nenc += len(fmt.Sprintf("%q", string(raw)))
	}
	err = s.Err()
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d - %d = %d\n", nstr, nmem, nstr-nmem)
	fmt.Printf("%d - %d = %d\n", nenc, nstr, nenc-nstr)
}
