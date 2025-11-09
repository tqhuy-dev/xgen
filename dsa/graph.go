package dsa

// NodeGraph represents a node in a directed graph with generic data and relation types
type NodeGraph[T any, K any] struct {
	Data      T
	Relation  K
	NextNodes []*NodeGraph[T, K]
}

// NewNodeGraph creates a new graph node with the given data and relation
func NewNodeGraph[T any, K any](data T, relation K) *NodeGraph[T, K] {
	return &NodeGraph[T, K]{
		Data:      data,
		Relation:  relation,
		NextNodes: make([]*NodeGraph[T, K], 0),
	}
}

// AddNextNode adds a child node to the current node's adjacency list
func (n *NodeGraph[T, K]) AddNextNode(node *NodeGraph[T, K]) {
	if node == nil {
		return
	}
	n.NextNodes = append(n.NextNodes, node)
}

// RemoveNextNode removes a specific child node from the adjacency list
// Returns true if the node was found and removed, false otherwise
func (n *NodeGraph[T, K]) RemoveNextNode(node *NodeGraph[T, K]) bool {
	if node == nil {
		return false
	}
	for i, child := range n.NextNodes {
		if child == node {
			n.NextNodes = append(n.NextNodes[:i], n.NextNodes[i+1:]...)
			return true
		}
	}
	return false
}

// HasNextNode checks if a specific node is in the adjacency list
func (n *NodeGraph[T, K]) HasNextNode(node *NodeGraph[T, K]) bool {
	if node == nil {
		return false
	}
	for _, child := range n.NextNodes {
		if child == node {
			return true
		}
	}
	return false
}

// Degree returns the number of outgoing edges (children) from this node
func (n *NodeGraph[T, K]) Degree() int {
	return len(n.NextNodes)
}

// IsLeaf returns true if the node has no children
func (n *NodeGraph[T, K]) IsLeaf() bool {
	return len(n.NextNodes) == 0
}

// ClearNextNodes removes all child nodes from the adjacency list
func (n *NodeGraph[T, K]) ClearNextNodes() {
	n.NextNodes = make([]*NodeGraph[T, K], 0)
}

// GetNextNodes performs a depth-first traversal of the graph starting from this node
// and calls handleNode for each visited node (including the starting node)
func (n *NodeGraph[T, K]) GetNextNodes(handleNode func(node *NodeGraph[T, K])) {
	if n == nil || handleNode == nil {
		return
	}
	children := n.NextNodes
	handleNode(n)
	if len(children) == 0 {
		return
	}
	for _, child := range children {
		child.GetNextNodes(handleNode)
	}
}

// TraverseDFS performs a depth-first traversal with cycle detection
// Calls handleNode for each visited node exactly once
func (n *NodeGraph[T, K]) TraverseDFS(handleNode func(node *NodeGraph[T, K])) {
	if n == nil || handleNode == nil {
		return
	}
	visited := make(map[*NodeGraph[T, K]]bool)
	n.traverseDFSHelper(handleNode, visited)
}

// traverseDFSHelper is a helper function for DFS traversal with cycle detection
func (n *NodeGraph[T, K]) traverseDFSHelper(handleNode func(node *NodeGraph[T, K]), visited map[*NodeGraph[T, K]]bool) {
	if visited[n] {
		return
	}
	visited[n] = true
	handleNode(n)
	for _, child := range n.NextNodes {
		if child != nil {
			child.traverseDFSHelper(handleNode, visited)
		}
	}
}

// TraverseBFS performs a breadth-first traversal with cycle detection
// Calls handleNode for each visited node exactly once
func (n *NodeGraph[T, K]) TraverseBFS(handleNode func(node *NodeGraph[T, K])) {
	if n == nil || handleNode == nil {
		return
	}
	visited := make(map[*NodeGraph[T, K]]bool)
	queue := []*NodeGraph[T, K]{n}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true
		handleNode(current)

		for _, child := range current.NextNodes {
			if child != nil && !visited[child] {
				queue = append(queue, child)
			}
		}
	}
}

// CountNodes counts the total number of unique nodes reachable from this node (including itself)
// Uses DFS with cycle detection
func (n *NodeGraph[T, K]) CountNodes() int {
	if n == nil {
		return 0
	}
	visited := make(map[*NodeGraph[T, K]]bool)
	return n.countNodesHelper(visited)
}

// countNodesHelper is a helper function for counting nodes with cycle detection
func (n *NodeGraph[T, K]) countNodesHelper(visited map[*NodeGraph[T, K]]bool) int {
	if visited[n] {
		return 0
	}
	visited[n] = true
	count := 1
	for _, child := range n.NextNodes {
		if child != nil {
			count += child.countNodesHelper(visited)
		}
	}
	return count
}

// HasCycle checks if there is a cycle in the graph reachable from this node
func (n *NodeGraph[T, K]) HasCycle() bool {
	if n == nil {
		return false
	}
	visited := make(map[*NodeGraph[T, K]]bool)
	recStack := make(map[*NodeGraph[T, K]]bool)
	return n.hasCycleHelper(visited, recStack)
}

// hasCycleHelper is a helper function for cycle detection using DFS
func (n *NodeGraph[T, K]) hasCycleHelper(visited, recStack map[*NodeGraph[T, K]]bool) bool {
	visited[n] = true
	recStack[n] = true

	for _, child := range n.NextNodes {
		if child == nil {
			continue
		}
		if !visited[child] {
			if child.hasCycleHelper(visited, recStack) {
				return true
			}
		} else if recStack[child] {
			return true
		}
	}

	recStack[n] = false
	return false
}
