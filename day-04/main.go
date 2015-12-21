package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

const secret = "iwrupvqb"

func main() {
	data := make(chan int)
	out := make(chan int)
	n := 100
	for i := 0; i < n; i++ {
		go process(data, out)
	}
	go generate(data)

	res := <-out
	fmt.Printf("res: %s%d\n", secret, res)
}

func generate(in chan int) {
	i := 0
	for {
		in <- i
		i++
	}
}

func process(in, out chan int) {
	for i := range in {
		txt := fmt.Sprintf("%s%d", secret, i)
		hash := md5sum(txt)
		if hash[:6] == "000000" {
			out <- i
			return
		}
	}
}

func md5sum(txt string) string {
	h := md5.New()
	h.Write([]byte(txt))
	return hex.EncodeToString(h.Sum(nil))
}
