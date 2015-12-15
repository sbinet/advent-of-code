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
	circuit = newCircuit()
)

type Wire struct {
	ID     string
	Source Signal
}

func (w Wire) Name() string {
	return w.ID
}

func (w Wire) Value() uint16 {
	return w.Source.Value()
}

type Signal interface {
	Name() string
	Value() uint16
}

type Gate struct {
	Operands []Signal
	Op       string
}

func (g Gate) Name() string {
	switch g.Op {
	case "AND":
		return g.Operands[0].Name() + " AND " + g.Operands[1].Name()
	case "OR":
		return g.Operands[0].Name() + " OR " + g.Operands[1].Name()
	case "LSHIFT":
		return g.Operands[0].Name() + " << " + g.Operands[1].Name()
	case "RSHIFT":
		return g.Operands[0].Name() + " >> " + g.Operands[1].Name()
	case "NOT":
		return "NOT " + g.Operands[0].Name()
	}
	panic("unreachable")
}

func (g Gate) Value() uint16 {
	vals := make([]uint16, len(g.Operands))
	for i, op := range g.Operands {
		vals[i] = op.Value()
	}
	return g.Eval(vals...)
}

func (g Gate) Eval(vals ...uint16) uint16 {
	switch g.Op {
	case "AND":
		return vals[0] & vals[1]
	case "OR":
		return vals[0] | vals[1]
	case "LSHIFT":
		return vals[0] << vals[1]
	case "RSHIFT":
		return vals[0] >> vals[1]
	case "NOT":
		return ^vals[0]
	}
	panic("unreachable")
}

func newGate(op string, operands []Signal) Gate {
	return Gate{
		Op:       op,
		Operands: operands,
	}
}

type Value struct {
	Source uint16
}

func (v Value) Name() string {
	return strconv.Itoa(int(v.Source))
}

func (v Value) Value() uint16 {
	return v.Source
}

type Circuit struct {
	wires map[string]*Wire
	cache map[string]uint16
}

func newCircuit() Circuit {
	return Circuit{
		wires: make(map[string]*Wire),
		cache: make(map[string]uint16),
	}
}

func (c *Circuit) add(w *Wire) {
	if _, dup := c.wires[w.ID]; dup {
		log.Panicf("duplicate wire [%s]", w.ID)
	}
	c.wires[w.ID] = w
}

func (c *Circuit) get(id string) *Wire {
	if w, ok := c.wires[id]; ok {
		return w
	}
	w := &Wire{ID: id}
	c.wires[id] = w
	return w
}

func (c *Circuit) signal(s string) Signal {
	switch isnumber(s) {
	case true:
		return Value{Source: atoi(s)}
	case false:
		return c.get(s)
	}
	panic("unreachable")
}

func (c *Circuit) process(id string, toks []string) {
	w := c.get(id)
	switch len(toks) {
	case 1:
		tok := toks[0]
		signal := c.signal(tok)
		w.Source = signal

	case 2:
		if toks[0] != "NOT" {
			log.Panicf("invalid command. got %q. want %q (%s -> %s)\n",
				toks[0], "NOT",
				strings.Join(toks, " "), id,
			)
		}
		signal := c.signal(toks[1])
		gate := newGate(toks[0], []Signal{signal})
		w.Source = gate
	case 3:
		lhs := c.signal(toks[0])
		op := toks[1]
		rhs := c.signal(toks[2])
		gate := newGate(op, []Signal{lhs, rhs})
		w.Source = gate
	default:
		log.Panicf("invalid tokens (len=%d tokens=%v)\n", len(toks), toks)
	}
	c.wires[id] = w
}

func (c *Circuit) Eval(id string) uint16 {
	if v, ok := c.cache[id]; ok {
		return v
	}
	indent := 0
	w, ok := c.wires[id]
	if !ok {
		panic("no such wire [" + id + "]")
	}
	return c.eval(w, indent)
}

func (c *Circuit) eval(signal Signal, indent int) uint16 {
	id := signal.Name()
	if v, ok := c.cache[id]; ok {
		return v
	}

	var v uint16
	switch sig := signal.(type) {
	case *Wire:
		v = c.eval(sig.Source, indent+1)
		c.cache[id] = v
	case Gate:
		vals := make([]uint16, len(sig.Operands))
		for i, op := range sig.Operands {
			vals[i] = c.eval(op, indent+1)
		}
		v = sig.Eval(vals...)
	case Value:
		v = sig.Value()
	}
	return v
}

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
	s := bufio.NewScanner(f)
	for s.Scan() {
		txt := s.Text()
		//fmt.Printf("%q\n", txt)
		toks := strings.Split(txt, " ")
		idx := len(toks) - 2
		if toks[idx] != "->" {
			log.Panicf("err: %v\n", txt)
		}
		circuit.process(toks[idx+1], toks[:idx])
	}
	err = s.Err()
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		panic(err)
	}

	fmt.Printf("--- a ---\n")
	v := circuit.Eval("a")
	fmt.Printf("a=%v\n", v)

	// override wire b to signal of wire a
	b := circuit.get("b")
	b.Source = Value{v}
	circuit.cache = make(map[string]uint16, len(circuit.cache))
	vv := circuit.Eval("a")
	fmt.Printf("a=%v\n", vv)
}

func atoi(s string) uint16 {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return uint16(v)
}

func isnumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
