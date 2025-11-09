package dsa

import (
	"sync"
	"sync/atomic"
)

type RoundRobinBasic[T any] struct {
	nodes []T
	index uint64
}

func NewRoundRobinBasic[T any](nodes []T) *RoundRobinBasic[T] {
	return &RoundRobinBasic[T]{nodes: nodes, index: 0}
}

func (rrb *RoundRobinBasic[T]) Next() T {
	if len(rrb.nodes) == 0 {
		var zero T
		return zero
	}
	i := atomic.AddUint64(&rrb.index, 1)
	return rrb.nodes[i%uint64(len(rrb.nodes))]
}

// RoundRobinNode represents a node with weight for weighted round-robin
type RoundRobinNode[T any] struct {
	Value         T
	Weight        int
	CurrentWeight int
}

// RoundRobinSmoothWeighted implements smooth weighted round-robin algorithm
type RoundRobinSmoothWeighted[T any] struct {
	nodes       []*RoundRobinNode[T]
	totalWeight int
	sync.Mutex
}

// WeightedNode is a helper structure for creating weighted round-robin
type WeightedNode[T any] struct {
	Value  T
	Weight int
}

// NewRoundRobinSmoothWeighted creates a new smooth weighted round-robin scheduler
// The nodes parameter is a slice of WeightedNode containing value and weight pairs
func NewRoundRobinSmoothWeighted[T any](nodes []WeightedNode[T]) *RoundRobinSmoothWeighted[T] {
	if len(nodes) == 0 {
		return &RoundRobinSmoothWeighted[T]{
			nodes:       []*RoundRobinNode[T]{},
			totalWeight: 0,
		}
	}

	result := &RoundRobinSmoothWeighted[T]{
		nodes:       make([]*RoundRobinNode[T], 0, len(nodes)),
		totalWeight: 0,
	}

	for _, node := range nodes {
		if node.Weight > 0 {
			result.nodes = append(result.nodes, &RoundRobinNode[T]{
				Value:         node.Value,
				Weight:        node.Weight,
				CurrentWeight: 0,
			})
			result.totalWeight += node.Weight
		}
	}

	return result
}

// Next returns the next node according to smooth weighted round-robin algorithm
func (rrs *RoundRobinSmoothWeighted[T]) Next() T {
	if len(rrs.nodes) == 0 {
		var zero T
		return zero
	}

	rrs.Lock()
	defer rrs.Unlock()

	var bestNode *RoundRobinNode[T]

	// Step 1 & 2: Add weight to current weight and find the node with max current weight
	for _, node := range rrs.nodes {
		node.CurrentWeight += node.Weight
		if bestNode == nil || node.CurrentWeight > bestNode.CurrentWeight {
			bestNode = node
		}
	}

	if bestNode == nil {
		return rrs.nodes[0].Value
	}

	// Step 3: Subtract total weight from the selected node's current weight
	bestNode.CurrentWeight -= rrs.totalWeight

	return bestNode.Value
}

// Reset resets all current weights to zero
func (rrs *RoundRobinSmoothWeighted[T]) Reset() {
	rrs.Lock()
	defer rrs.Unlock()

	for _, node := range rrs.nodes {
		node.CurrentWeight = 0
	}
}

// RoundRobinLimitDDLNode represents a node with a concurrency limit
type RoundRobinLimitDDLNode[T any] struct {
	Value        T
	Limit        int // Maximum concurrent uses allowed
	currentLimit int // Current usage count
	inList       bool // Tracks if node is currently in the list
}

// RoundRobinLimitDDL implements a round-robin scheduler with per-node concurrency limits
// Uses a doubly-linked list to efficiently rotate nodes
type RoundRobinLimitDDL[T any] struct {
	list *DLList[*RoundRobinLimitDDLNode[T]]
	sync.Mutex
}

// NewRoundRobinLimitDDL creates a new limited round-robin scheduler
// Each node has a limit on how many times it can be used concurrently
func NewRoundRobinLimitDDL[T any](nodes []*RoundRobinLimitDDLNode[T]) *RoundRobinLimitDDL[T] {
	list := NewDLList[*RoundRobinLimitDDLNode[T]]()
	for _, node := range nodes {
		if node.Limit > 0 {
			node.inList = true
			list.PushBack(node)
		}
	}
	return &RoundRobinLimitDDL[T]{list: list}
}

// Next returns the next available node value
// Increments the node's usage count and moves it to the back of the queue
// If the node reaches its limit, it's removed from rotation until released
func (rrl *RoundRobinLimitDDL[T]) Next() T {
	rrl.Lock()
	defer rrl.Unlock()

	if rrl.list.Len() == 0 {
		var zero T
		return zero
	}

	nodeFirst := rrl.list.Front()
	if nodeFirst == nil {
		var zero T
		return zero
	}

	// Increment usage count
	nodeFirst.Value.currentLimit++

	// Check if node reached its limit
	if nodeFirst.Value.currentLimit >= nodeFirst.Value.Limit {
		// Remove from list if at capacity
		rrl.list.Remove(nodeFirst)
		nodeFirst.Value.inList = false
	} else {
		// Move to back if still under limit
		rrl.list.MoveToBack(nodeFirst)
	}

	return nodeFirst.Value.Value
}

// Release decrements a node's usage count and adds it back to rotation if under limit
func (rrl *RoundRobinLimitDDL[T]) Release(node *RoundRobinLimitDDLNode[T]) {
	rrl.Lock()
	defer rrl.Unlock()

	if node == nil {
		return
	}

	// Decrement usage count
	if node.currentLimit > 0 {
		node.currentLimit--
	}

	// Add back to list if under limit and not already in list
	if node.currentLimit < node.Limit && !node.inList {
		node.inList = true
		rrl.list.PushBack(node)
	}
}

// Len returns the current number of available nodes
func (rrl *RoundRobinLimitDDL[T]) Len() int {
	rrl.Lock()
	defer rrl.Unlock()
	return rrl.list.Len()
}

// Reset resets all nodes' current limits to zero and rebuilds the list
func (rrl *RoundRobinLimitDDL[T]) Reset(nodes []*RoundRobinLimitDDLNode[T]) {
	rrl.Lock()
	defer rrl.Unlock()

	rrl.list = NewDLList[*RoundRobinLimitDDLNode[T]]()
	for _, node := range nodes {
		if node.Limit > 0 {
			node.currentLimit = 0
			node.inList = true
			rrl.list.PushBack(node)
		}
	}
}
