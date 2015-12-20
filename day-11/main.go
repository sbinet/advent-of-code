package main

import (
	"bytes"
	"fmt"
)

func main() {
	for _, v := range []string{
		"abcdefgh",
		"ghijklmn",
		"vzbxkghb",
		"vzbxxyzz",
	} {
		vv := next([]byte(v))
		fmt.Printf("%s => %s\n", v, string(vv))
	}
}

func inc(pwd []byte) {
	for i := len(pwd) - 1; i >= 0; i-- {
		pwd[i]++
		if pwd[i] > 'z' {
			pwd[i] = 'a'
			continue
		}
		return
	}
}

func next(old []byte) []byte {
	pwd := make([]byte, len(old))
	copy(pwd, old)
	for {
		inc(pwd)

		if !req1(pwd) ||
			!req2(pwd) ||
			!req3(pwd) {
			continue
		}
		return pwd
	}
	fmt.Printf("**error**\n")
	return pwd
}

func req1(pwd []byte) bool {
	n := 0
	for i := range pwd[:len(pwd)-2] {
		u := pwd[i]
		v := pwd[i+1]
		w := pwd[i+2]
		if u == v-1 && v == w-1 {
			n++
		}
	}
	return n > 0
}

func req2(pwd []byte) bool {
	return !(bytes.Contains(pwd, []byte("i")) ||
		bytes.Contains(pwd, []byte("o")) ||
		bytes.Contains(pwd, []byte("l")))
}

func req3(pwd []byte) bool {
	n := 0
	pairs := make(map[byte]bool)
	for i := 1; i < len(pwd); i++ {
		v0 := pwd[i-1]
		v1 := pwd[i]
		if v1 == v0 {
			_, dup := pairs[v0]
			if !dup {
				n++
				i++
				pairs[v0] = true
			}
		}
	}
	return n >= 2
}
