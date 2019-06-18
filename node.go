package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/valyala/fastjson"
)

type node struct {
	children []*node
	value    value
	parent   *node
}

func (n *node) MarshalJSON() ([]byte, error) {
	b := bytes.NewBuffer(nil)
	n.marshalJSON(b)
	return b.Bytes(), nil
}

func (n *node) marshalJSON(b *bytes.Buffer) {
	b.WriteRune('[')
	for i, v := range n.children {
		if i != 0 {
			b.WriteRune(',')
		}
		if v.children != nil {
			v.marshalJSON(b)
		} else {
			fmt.Fprint(b, v.value)
		}
	}
	b.WriteRune(']')
}

func (n *node) fromJSONArray(array []*fastjson.Value) int {
	values := 0
	for _, v := range array {
		curr := &node{}
		switch v.Type() {
		case fastjson.TypeArray:
			array, err := v.Array()
			must(err)
			values += curr.fromJSONArray(array)
		case fastjson.TypeNumber:
			values++
			num, err := v.Int64()
			must(err)
			curr.value = num
		default:
			log.Fatalf("Unexpected type %s", v.Type().String())
		}

		curr.parent = n
		n.children = append(n.children, curr)
	}
	return values
}
