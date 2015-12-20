package main

import (
	"bufio"
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
	var deers []Reindeer
	s := bufio.NewScanner(f)
	for s.Scan() {
		var name string
		var speed int
		var t, rest int
		txt := s.Text()
		_, err = fmt.Sscanf(
			txt,
			"%s can fly %d km/s for %d seconds, but then must rest for %d seconds.",
			&name, &speed, &t, &rest,
		)
		if err != nil {
			panic(err)
		}
		deers = append(deers, Reindeer{
			Name:  name,
			Speed: speed,
			Time:  t,
			Rest:  rest,
			gauge: t,
		})
	}
	err = s.Err()
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		panic(err)
	}

	for i := 0; i < 2503; i++ {
		for j := range deers {
			deer := &deers[j]
			deer.tick()
		}
		points(deers)
	}

	for _, d := range deers {
		fmt.Printf("%s: %d %d\n", d.Name, d.len, d.points)
	}
}

type Reindeer struct {
	Name  string
	Speed int
	Time  int
	Rest  int

	state  State
	gauge  int
	len    int
	points int
}

func (d *Reindeer) tick() {
	switch d.state {
	case running:
		d.len += d.Speed
		d.gauge--
	case resting:
		d.gauge--
	}
	if d.gauge == 0 {
		switch d.state {
		case running:
			d.state = resting
			d.gauge = d.Rest
		case resting:
			d.state = running
			d.gauge = d.Time
		}
	}
}

func points(deers []Reindeer) {
	max := 0
	for _, deer := range deers {
		if deer.len >= max {
			max = deer.len
		}
	}
	for i := range deers {
		deer := &deers[i]
		if deer.len == max {
			deer.points++
		}
	}
}

type State int

const (
	running State = 0
	resting State = 1
)
