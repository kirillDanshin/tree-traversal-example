// Copyright 2019-present Kirill Danshin
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0

package main

func GenerateNodes(n int) *root {
	root := New()
	curr := &root.node
	root.closePool()
	root.visitorPool, root.closePool = stackless(n, visitAllParallel)

	for i := 0; i < n; i++ {
		v := &node{
			value:  value(i),
			parent: curr,
		}

		curr.children = append(curr.children, v)
		root.values++

		if i > 0 {
			if i%100 == 0 {
				newNode := &node{
					parent: curr,
				}
				curr.children = append(curr.children, newNode)
				curr = newNode
				continue
			}
			if i%200 == 0 {
				curr = curr.parent
			}
		}
	}
	return root
}
