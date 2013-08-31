package main

import "fmt"

type List struct {
	data int
	nextNode *List
}

func main() {
	//head := &List{1, &List{2, &List{3, nil} } }
	head := new(List)
	/**
	second := new(List)
	third := new(List)

	head.data = 100
	second.data = 200
	third.data = 300
	head.nextNode = second
	second.nextNode = third
	**/

	recurse(1, head)

	for node := head; node != nil; node = node.nextNode {
		fmt.Println(node.data)
	}
}

func recurse(data int, node *List) {
	node.data = data
	if data >= 10 { return }

	node.nextNode = new(List)
	recurse(data + 1, node.nextNode)
}
