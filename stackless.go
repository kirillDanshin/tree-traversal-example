package main

import (
	"runtime"
	"sync"
)

func stackless(poolSize int, f func(*root, *node, func(value))) (func(*root, *node, func(value)) bool, func()) {
	if f == nil {
		panic("BUG: f cannot be nil")
	}

	funcWorkCh := make(chan *funcWork, runtime.GOMAXPROCS(-1)*2048)
	onceInit := func() {
		n := runtime.GOMAXPROCS(-1)
		for i := 0; i < n; i++ {
			go funcWorker(funcWorkCh, f)
		}
	}
	var once sync.Once

	run := func(r *root, n *node, f func(value)) bool {
		once.Do(onceInit)
		fw := getFuncWork()
		fw.node = n
		fw.f = f
		fw.r = r

		select {
		case funcWorkCh <- fw:
		default:
			putFuncWork(fw)
			return false
		}
		return true
	}
	close := func() {
		close(funcWorkCh)
	}

	return run, close
}

func funcWorker(funcWorkCh <-chan *funcWork, f func(*root, *node, func(value))) {
	for fw := range funcWorkCh {
		f(fw.r, fw.node, fw.f)
		putFuncWork(fw)
	}
}

func getFuncWork() *funcWork {
	v := funcWorkPool.Get()
	if v == nil {
		v = &funcWork{}
	}
	return v.(*funcWork)
}

func putFuncWork(fw *funcWork) {
	fw.node = nil
	fw.f = nil
	fw.r = nil
	funcWorkPool.Put(fw)
}

var funcWorkPool sync.Pool

type funcWork struct {
	node *node
	r    *root
	f    func(value)
}
