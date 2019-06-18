package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"sync"
	"time"

	"github.com/spf13/pflag"
)

var (
	filepath  = pflag.StringP("filepath", "f", "", "json file path")
	generateN = pflag.IntP("generateNodes", "g", 0, "generate N nodes and write to a given filepath")
	parallel  = pflag.BoolP("parallel", "p", false, "enable parallel method. uses waitgroup in each method")
)

type value = int64

func main() {
	pflag.Parse()

	var r *root
	if *generateN == 0 {
		rawData := mustRead(*filepath)
		r = CreateGraphFromJSON(rawData)
	} else {
		r = GenerateNodes(*generateN)
	}
	if *parallel {
		visitWithSync("recursive approach", r, (*root).VisitAllRecursive)
		visitWithSync("iterative approach", r, (*root).VisitAll)
		visitWithSync("parallel approach", r, (*root).VisitAllParallel)
	} else {
		visitWith("recursive approach", r, (*root).VisitAllRecursive)
		visitWith("iterative approach", r, (*root).VisitAll)
	}

	b, err := r.MarshalJSON()
	must(err)
	if len(*filepath) > 0 && *generateN > 0 {
		must(ioutil.WriteFile(*filepath, b, 0644))
	}
}

func visitWith(name string, r *root, visitMethod func(*root, func(value))) {
	var deoptimize value // this will deny any visitor optimizations

	start := time.Now().UnixNano()
	visitMethod(r, func(v value) {
		deoptimize = v
	})
	end := time.Now().UnixNano()

	fmt.Printf("Visited every node using %s in %s\n", name, time.Duration(end-start))
	fmt.Fprint(ioutil.Discard, "last value was", deoptimize)
}

func visitWithSync(name string, r *root, visitMethod func(*root, func(value))) {
	var deoptimize value // this will deny any visitor optimizations

	wg := sync.WaitGroup{}
	wg.Add(r.values)
	start := time.Now().UnixNano()
	visitMethod(r, func(v value) {
		deoptimize = v
		wg.Done()
	})
	wg.Wait()
	end := time.Now().UnixNano()

	fmt.Printf("Visited every node using %s in %s\n", name, time.Duration(end-start))
	fmt.Fprint(ioutil.Discard, "last value was", deoptimize)
}

func mustRead(filePath string) []byte {
	b, err := ioutil.ReadFile(path.Clean(filePath))
	must(err)
	return b
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
