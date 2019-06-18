package main

import (
	"unsafe"

	"github.com/brentp/intintmap"
)

func (n *node) offsetID() int64 {
	return int64(uintptr(unsafe.Pointer(n)))
}

func (r *root) VisitAllParallel(visitor func(value)) {
	visitAllParallel(r, &r.node, visitor)
}

func visitAllParallel(r *root, n *node, visitor func(value)) {
	for _, v := range n.children {
		if v.children != nil {
			if !r.visitorPool(r, v, visitor) {
				visitAllParallel(r, v, visitor)
			}
		} else {
			visitor(v.value)
		}
	}
}

func (r *root) VisitAll(visitor func(value)) {
	curr := &r.node
	offsets := intintmap.New(r.values/2, .96)
outerLoop:
	for {
		offset, ok := offsets.Get(curr.offsetID())
		if ok && offset == int64(len(curr.children))-1 {
			curr = curr.parent
			offsets.Del(curr.offsetID())
			continue
		}
		for _, v := range curr.children[offset:] {
			offset++
			offsets.Put(curr.offsetID(), offset)
			if v.children != nil {
				curr = v
				continue outerLoop
			} else {
				visitor(v.value)
			}
		}
		if curr.parent == nil {
			break
		}
		curr = curr.parent
	}
}

func (n *node) VisitAllRecursive(visitor func(value)) {
	curr := n
	for _, v := range curr.children {
		if v.children != nil {
			v.VisitAllRecursive(visitor)
		} else {
			visitor(v.value)
		}
	}
}
