package main

import (
	"fmt"

	"github.com/tqhuy-dev/xgen/dsa"
)

func main() {
	nodeA := &dsa.NodeGraph[int, string]{
		Data:     1,
		Relation: "relation1",
	}
	nodeB := &dsa.NodeGraph[int, string]{
		Data:     2,
		Relation: "relation2",
	}
	nodeC := &dsa.NodeGraph[int, string]{
		Data:     3,
		Relation: "relation3",
	}
	nodeD := &dsa.NodeGraph[int, string]{
		Data:     4,
		Relation: "relation3",
	}
	nodeE := &dsa.NodeGraph[int, string]{
		Data:     5,
		Relation: "relation3",
	}
	nodeF := &dsa.NodeGraph[int, string]{
		Data:     6,
		Relation: "relation3",
	}
	nodeG := &dsa.NodeGraph[int, string]{
		Data: 7,
	}
	nodeA.AddNextNode(nodeB)
	nodeA.AddNextNode(nodeC)
	nodeC.AddNextNode(nodeD)
	nodeC.AddNextNode(nodeE)
	nodeE.AddNextNode(nodeF)
	nodeB.AddNextNode(nodeG)
	nodeA.TraverseDFS(func(node *dsa.NodeGraph[int, string]) {
		fmt.Println(node.Data)
	})
}
