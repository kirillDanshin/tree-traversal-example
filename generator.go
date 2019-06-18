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
