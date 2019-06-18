// Copyright 2019-present Kirill Danshin
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0

package main

import (
	"log"

	"github.com/valyala/fastjson"
)

func New() *root {
	n := &root{}

	n.visitorPool, n.closePool = stackless(n.values, visitAllParallel)

	return n
}

func CreateGraphFromJSON(b []byte) *root {
	n := New()

	p := fastjson.Parser{}
	v, err := p.ParseBytes(b)
	must(err)
	n.fromJSONRoot(v)
	n.closePool()
	n.visitorPool, n.closePool = stackless(n.values, visitAllParallel)
	return n
}

type root struct {
	node
	visitorPool func(*root, *node, func(value)) bool
	closePool   func()
	values      int
}

func (r *root) fromJSONRoot(root *fastjson.Value) {
	list, err := root.Array()
	must(err)
	for _, v := range list {
		curr := &node{}
		switch v.Type() {
		case fastjson.TypeArray:
			array, err := v.Array()
			must(err)
			r.values += curr.fromJSONArray(array)

		case fastjson.TypeNumber:
			num, err := v.Int64()
			must(err)
			curr.value = num
			r.values++
		default:
			log.Fatalf("Unexpected type %s", v.Type().String())
		}
		curr.parent = &r.node
		r.children = append(r.children, curr)
	}
}
