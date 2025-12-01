package main

import (
	"fmt"

	"github.com/tqhuy-dev/xgen/dsa"
)

func main() {
	newLinkList := dsa.NewLinkList[int]()
	newLinkList.AddNode(&dsa.LinkListNode[int]{Data: 1})
	newLinkList.AddNode(&dsa.LinkListNode[int]{Data: 3})
	newLinkList.AddNode(&dsa.LinkListNode[int]{Data: 5})
	newLinkList.AddNode(&dsa.LinkListNode[int]{Data: 7})

	newLinkList.Reverse()
	newLinkList.Iteration(func(node *dsa.LinkListNode[int]) (*dsa.LinkListNode[int], bool, error) {
		if node != nil {
			fmt.Println(node.Data)
		}
		return node, false, nil
	})
}
