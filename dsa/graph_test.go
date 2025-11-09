package dsa

import (
	"reflect"
	"sort"
	"testing"
)

// TestNewNodeGraph tests the NewNodeGraph constructor
// which creates a new graph node with given data and relation
func TestNewNodeGraph(t *testing.T) {
	tests := []struct {
		name     string
		data     int
		relation string
	}{
		{
			name:     "create node with integer data and string relation",
			data:     42,
			relation: "parent",
		},
		{
			name:     "create node with zero value data",
			data:     0,
			relation: "child",
		},
		{
			name:     "create node with empty relation",
			data:     100,
			relation: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			node := NewNodeGraph(tt.data, tt.relation)

			// Assert
			if node == nil {
				t.Fatal("NewNodeGraph() returned nil")
			}
			if node.Data != tt.data {
				t.Errorf("NewNodeGraph().Data = %v; expected %v", node.Data, tt.data)
			}
			if node.Relation != tt.relation {
				t.Errorf("NewNodeGraph().Relation = %v; expected %v", node.Relation, tt.relation)
			}
			if node.NextNodes == nil {
				t.Error("NewNodeGraph().NextNodes is nil; expected empty slice")
			}
			if len(node.NextNodes) != 0 {
				t.Errorf("NewNodeGraph().NextNodes length = %d; expected 0", len(node.NextNodes))
			}
		})
	}
}

// TestNodeGraph_AddNextNode tests the AddNextNode method
// which adds a child node to the adjacency list
func TestNodeGraph_AddNextNode(t *testing.T) {
	tests := []struct {
		name          string
		initialNodes  int
		nodesToAdd    int
		expectedCount int
	}{
		{
			name:          "add one node to empty graph",
			initialNodes:  0,
			nodesToAdd:    1,
			expectedCount: 1,
		},
		{
			name:          "add multiple nodes",
			initialNodes:  0,
			nodesToAdd:    5,
			expectedCount: 5,
		},
		{
			name:          "add nodes to existing graph",
			initialNodes:  3,
			nodesToAdd:    2,
			expectedCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			root := NewNodeGraph("root", "main")
			for i := 0; i < tt.initialNodes; i++ {
				root.AddNextNode(NewNodeGraph("initial", "init"))
			}

			// Act
			for i := 0; i < tt.nodesToAdd; i++ {
				child := NewNodeGraph("child", "relation")
				root.AddNextNode(child)
			}

			// Assert
			if len(root.NextNodes) != tt.expectedCount {
				t.Errorf("AddNextNode() resulted in %d nodes; expected %d", len(root.NextNodes), tt.expectedCount)
			}
		})
	}
}

// TestNodeGraph_AddNextNode_Nil tests adding a nil node
// which should be safely ignored
func TestNodeGraph_AddNextNode_Nil(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "add nil node should not panic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			root := NewNodeGraph(1, "root")
			initialCount := len(root.NextNodes)

			// Act
			root.AddNextNode(nil)

			// Assert
			if len(root.NextNodes) != initialCount {
				t.Errorf("AddNextNode(nil) changed node count from %d to %d; expected no change", initialCount, len(root.NextNodes))
			}
		})
	}
}

// TestNodeGraph_RemoveNextNode tests the RemoveNextNode method
// which removes a specific child node from the adjacency list
func TestNodeGraph_RemoveNextNode(t *testing.T) {
	tests := []struct {
		name           string
		setupNodes     int
		removeIndex    int
		expectedResult bool
		expectedCount  int
	}{
		{
			name:           "remove existing node",
			setupNodes:     3,
			removeIndex:    1,
			expectedResult: true,
			expectedCount:  2,
		},
		{
			name:           "remove first node",
			setupNodes:     3,
			removeIndex:    0,
			expectedResult: true,
			expectedCount:  2,
		},
		{
			name:           "remove last node",
			setupNodes:     3,
			removeIndex:    2,
			expectedResult: true,
			expectedCount:  2,
		},
		{
			name:           "remove from single node graph",
			setupNodes:     1,
			removeIndex:    0,
			expectedResult: true,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			root := NewNodeGraph("root", "main")
			nodes := make([]*NodeGraph[string, string], tt.setupNodes)
			for i := 0; i < tt.setupNodes; i++ {
				nodes[i] = NewNodeGraph("child", "relation")
				root.AddNextNode(nodes[i])
			}

			// Act
			result := root.RemoveNextNode(nodes[tt.removeIndex])

			// Assert
			if result != tt.expectedResult {
				t.Errorf("RemoveNextNode() = %v; expected %v", result, tt.expectedResult)
			}
			if len(root.NextNodes) != tt.expectedCount {
				t.Errorf("RemoveNextNode() resulted in %d nodes; expected %d", len(root.NextNodes), tt.expectedCount)
			}
		})
	}
}

// TestNodeGraph_RemoveNextNode_NotFound tests removing a node that doesn't exist
// which should return false
func TestNodeGraph_RemoveNextNode_NotFound(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "remove non-existent node returns false",
		},
		{
			name: "remove nil node returns false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			root := NewNodeGraph(1, "root")
			child1 := NewNodeGraph(2, "child1")
			root.AddNextNode(child1)

			otherNode := NewNodeGraph(4, "other")

			// Act
			result1 := root.RemoveNextNode(otherNode)
			result2 := root.RemoveNextNode(nil)

			// Assert
			if result1 {
				t.Error("RemoveNextNode(non-existent) = true; expected false")
			}
			if result2 {
				t.Error("RemoveNextNode(nil) = true; expected false")
			}
			if len(root.NextNodes) != 1 {
				t.Errorf("Node count = %d; expected 1 (unchanged)", len(root.NextNodes))
			}
		})
	}
}

// TestNodeGraph_HasNextNode tests the HasNextNode method
// which checks if a node exists in the adjacency list
func TestNodeGraph_HasNextNode(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() (*NodeGraph[int, string], *NodeGraph[int, string])
		expected bool
	}{
		{
			name: "node exists in adjacency list",
			setup: func() (*NodeGraph[int, string], *NodeGraph[int, string]) {
				root := NewNodeGraph(1, "root")
				child := NewNodeGraph(2, "child")
				root.AddNextNode(child)
				return root, child
			},
			expected: true,
		},
		{
			name: "node does not exist in adjacency list",
			setup: func() (*NodeGraph[int, string], *NodeGraph[int, string]) {
				root := NewNodeGraph(1, "root")
				child := NewNodeGraph(2, "child")
				root.AddNextNode(NewNodeGraph(3, "other"))
				return root, child
			},
			expected: false,
		},
		{
			name: "empty adjacency list",
			setup: func() (*NodeGraph[int, string], *NodeGraph[int, string]) {
				root := NewNodeGraph(1, "root")
				child := NewNodeGraph(2, "child")
				return root, child
			},
			expected: false,
		},
		{
			name: "nil node returns false",
			setup: func() (*NodeGraph[int, string], *NodeGraph[int, string]) {
				root := NewNodeGraph(1, "root")
				root.AddNextNode(NewNodeGraph(2, "child"))
				return root, nil
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			root, child := tt.setup()

			// Act
			result := root.HasNextNode(child)

			// Assert
			if result != tt.expected {
				t.Errorf("HasNextNode() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestNodeGraph_Degree tests the Degree method
// which returns the number of outgoing edges
func TestNodeGraph_Degree(t *testing.T) {
	tests := []struct {
		name          string
		childrenCount int
		expected      int
	}{
		{
			name:          "node with no children has degree 0",
			childrenCount: 0,
			expected:      0,
		},
		{
			name:          "node with one child has degree 1",
			childrenCount: 1,
			expected:      1,
		},
		{
			name:          "node with multiple children",
			childrenCount: 5,
			expected:      5,
		},
		{
			name:          "node with many children",
			childrenCount: 10,
			expected:      10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			root := NewNodeGraph("root", "main")
			for i := 0; i < tt.childrenCount; i++ {
				root.AddNextNode(NewNodeGraph("child", "relation"))
			}

			// Act
			result := root.Degree()

			// Assert
			if result != tt.expected {
				t.Errorf("Degree() = %d; expected %d", result, tt.expected)
			}
		})
	}
}

// TestNodeGraph_IsLeaf tests the IsLeaf method
// which returns true if the node has no children
func TestNodeGraph_IsLeaf(t *testing.T) {
	tests := []struct {
		name          string
		childrenCount int
		expected      bool
	}{
		{
			name:          "node with no children is a leaf",
			childrenCount: 0,
			expected:      true,
		},
		{
			name:          "node with one child is not a leaf",
			childrenCount: 1,
			expected:      false,
		},
		{
			name:          "node with multiple children is not a leaf",
			childrenCount: 5,
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			root := NewNodeGraph("root", "main")
			for i := 0; i < tt.childrenCount; i++ {
				root.AddNextNode(NewNodeGraph("child", "relation"))
			}

			// Act
			result := root.IsLeaf()

			// Assert
			if result != tt.expected {
				t.Errorf("IsLeaf() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestNodeGraph_ClearNextNodes tests the ClearNextNodes method
// which removes all child nodes
func TestNodeGraph_ClearNextNodes(t *testing.T) {
	tests := []struct {
		name              string
		initialChildCount int
	}{
		{
			name:              "clear empty node",
			initialChildCount: 0,
		},
		{
			name:              "clear node with one child",
			initialChildCount: 1,
		},
		{
			name:              "clear node with multiple children",
			initialChildCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			root := NewNodeGraph("root", "main")
			for i := 0; i < tt.initialChildCount; i++ {
				root.AddNextNode(NewNodeGraph("child", "relation"))
			}

			// Act
			root.ClearNextNodes()

			// Assert
			if len(root.NextNodes) != 0 {
				t.Errorf("ClearNextNodes() resulted in %d nodes; expected 0", len(root.NextNodes))
			}
			if !root.IsLeaf() {
				t.Error("ClearNextNodes() node is not a leaf after clearing")
			}
		})
	}
}

// TestNodeGraph_GetNextNodes tests the GetNextNodes method
// which performs DFS traversal with a callback
func TestNodeGraph_GetNextNodes(t *testing.T) {
	tests := []struct {
		name          string
		buildGraph    func() *NodeGraph[int, string]
		expectedOrder []int
	}{
		{
			name: "single node",
			buildGraph: func() *NodeGraph[int, string] {
				return NewNodeGraph(1, "root")
			},
			expectedOrder: []int{1},
		},
		{
			name: "linear chain",
			buildGraph: func() *NodeGraph[int, string] {
				root := NewNodeGraph(1, "root")
				child1 := NewNodeGraph(2, "child1")
				child2 := NewNodeGraph(3, "child2")
				root.AddNextNode(child1)
				child1.AddNextNode(child2)
				return root
			},
			expectedOrder: []int{1, 2, 3},
		},
		{
			name: "tree with multiple children",
			buildGraph: func() *NodeGraph[int, string] {
				root := NewNodeGraph(1, "root")
				child1 := NewNodeGraph(2, "child1")
				child2 := NewNodeGraph(3, "child2")
				root.AddNextNode(child1)
				root.AddNextNode(child2)
				return root
			},
			expectedOrder: []int{1, 2, 3},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			root := tt.buildGraph()
			var visited []int

			// Act
			root.GetNextNodes(func(node *NodeGraph[int, string]) {
				visited = append(visited, node.Data)
			})

			// Assert
			if !reflect.DeepEqual(visited, tt.expectedOrder) {
				t.Errorf("GetNextNodes() visited order = %v; expected %v", visited, tt.expectedOrder)
			}
		})
	}
}

// TestNodeGraph_GetNextNodes_NilHandling tests GetNextNodes with nil inputs
// which should be handled gracefully
func TestNodeGraph_GetNextNodes_NilHandling(t *testing.T) {
	tests := []struct {
		name        string
		node        *NodeGraph[int, string]
		handleNode  func(node *NodeGraph[int, string])
		shouldPanic bool
	}{
		{
			name:        "nil node with valid handler",
			node:        nil,
			handleNode:  func(node *NodeGraph[int, string]) {},
			shouldPanic: false,
		},
		{
			name:        "valid node with nil handler",
			node:        NewNodeGraph(1, "test"),
			handleNode:  nil,
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act & Assert
			defer func() {
				r := recover()
				if (r != nil) != tt.shouldPanic {
					t.Errorf("GetNextNodes() panic = %v; expected panic = %v", r != nil, tt.shouldPanic)
				}
			}()
			tt.node.GetNextNodes(tt.handleNode)
		})
	}
}

// TestNodeGraph_TraverseDFS tests the TraverseDFS method
// which performs depth-first traversal with cycle detection
func TestNodeGraph_TraverseDFS(t *testing.T) {
	tests := []struct {
		name             string
		buildGraph       func() *NodeGraph[int, string]
		expectedCount    int
		expectedMaxLevel int
	}{
		{
			name: "simple tree without cycles",
			buildGraph: func() *NodeGraph[int, string] {
				root := NewNodeGraph(1, "root")
				child1 := NewNodeGraph(2, "child1")
				child2 := NewNodeGraph(3, "child2")
				root.AddNextNode(child1)
				root.AddNextNode(child2)
				return root
			},
			expectedCount:    3,
			expectedMaxLevel: 2,
		},
		{
			name: "graph with cycle",
			buildGraph: func() *NodeGraph[int, string] {
				root := NewNodeGraph(1, "root")
				child := NewNodeGraph(2, "child")
				root.AddNextNode(child)
				child.AddNextNode(root) // Creates cycle
				return root
			},
			expectedCount:    2, // Should visit each node only once
			expectedMaxLevel: 2,
		},
		{
			name: "single node",
			buildGraph: func() *NodeGraph[int, string] {
				return NewNodeGraph(1, "single")
			},
			expectedCount:    1,
			expectedMaxLevel: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			root := tt.buildGraph()
			visited := make(map[int]bool)
			maxLevel := 0
			// Act
			root.TraverseDFS(func(node *NodeGraph[int, string], level int) {
				visited[node.Data] = true
				if level > maxLevel {
					maxLevel = level
				}
			})

			// Assert
			if len(visited) != tt.expectedCount {
				t.Errorf("TraverseDFS() visited %d nodes; expected %d", len(visited), tt.expectedCount)
			}
			if maxLevel != tt.expectedMaxLevel {
				t.Errorf("TraverseDFS() max level = %d; expected %d", maxLevel, tt.expectedMaxLevel)
			}
		})
	}
}

// TestNodeGraph_TraverseBFS tests the TraverseBFS method
// which performs breadth-first traversal with cycle detection
func TestNodeGraph_TraverseBFS(t *testing.T) {
	tests := []struct {
		name          string
		buildGraph    func() *NodeGraph[int, string]
		expectedCount int
		checkOrder    bool
		expectedOrder []int
	}{
		{
			name: "simple tree BFS order",
			buildGraph: func() *NodeGraph[int, string] {
				root := NewNodeGraph(1, "root")
				child1 := NewNodeGraph(2, "child1")
				child2 := NewNodeGraph(3, "child2")
				grandchild := NewNodeGraph(4, "grandchild")
				root.AddNextNode(child1)
				root.AddNextNode(child2)
				child1.AddNextNode(grandchild)
				return root
			},
			expectedCount: 4,
			checkOrder:    true,
			expectedOrder: []int{1, 2, 3, 4},
		},
		{
			name: "graph with cycle",
			buildGraph: func() *NodeGraph[int, string] {
				root := NewNodeGraph(1, "root")
				child := NewNodeGraph(2, "child")
				root.AddNextNode(child)
				child.AddNextNode(root) // Creates cycle
				return root
			},
			expectedCount: 2,
			checkOrder:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			root := tt.buildGraph()
			var visited []int

			// Act
			root.TraverseBFS(func(node *NodeGraph[int, string]) {
				visited = append(visited, node.Data)
			})

			// Assert
			if len(visited) != tt.expectedCount {
				t.Errorf("TraverseBFS() visited %d nodes; expected %d", len(visited), tt.expectedCount)
			}
			if tt.checkOrder && !reflect.DeepEqual(visited, tt.expectedOrder) {
				t.Errorf("TraverseBFS() order = %v; expected %v", visited, tt.expectedOrder)
			}
		})
	}
}

// TestNodeGraph_CountNodes tests the CountNodes method
// which counts all unique reachable nodes
func TestNodeGraph_CountNodes(t *testing.T) {
	tests := []struct {
		name          string
		buildGraph    func() *NodeGraph[int, string]
		expectedCount int
	}{
		{
			name: "single node",
			buildGraph: func() *NodeGraph[int, string] {
				return NewNodeGraph(1, "single")
			},
			expectedCount: 1,
		},
		{
			name: "linear chain",
			buildGraph: func() *NodeGraph[int, string] {
				root := NewNodeGraph(1, "root")
				child1 := NewNodeGraph(2, "child1")
				child2 := NewNodeGraph(3, "child2")
				root.AddNextNode(child1)
				child1.AddNextNode(child2)
				return root
			},
			expectedCount: 3,
		},
		{
			name: "tree structure",
			buildGraph: func() *NodeGraph[int, string] {
				root := NewNodeGraph(1, "root")
				child1 := NewNodeGraph(2, "child1")
				child2 := NewNodeGraph(3, "child2")
				grandchild1 := NewNodeGraph(4, "grandchild1")
				grandchild2 := NewNodeGraph(5, "grandchild2")
				root.AddNextNode(child1)
				root.AddNextNode(child2)
				child1.AddNextNode(grandchild1)
				child2.AddNextNode(grandchild2)
				return root
			},
			expectedCount: 5,
		},
		{
			name: "graph with cycle counts each node once",
			buildGraph: func() *NodeGraph[int, string] {
				root := NewNodeGraph(1, "root")
				child := NewNodeGraph(2, "child")
				root.AddNextNode(child)
				child.AddNextNode(root) // Creates cycle
				return root
			},
			expectedCount: 2,
		},
		{
			name: "diamond pattern (DAG with shared node)",
			buildGraph: func() *NodeGraph[int, string] {
				root := NewNodeGraph(1, "root")
				left := NewNodeGraph(2, "left")
				right := NewNodeGraph(3, "right")
				bottom := NewNodeGraph(4, "bottom")
				root.AddNextNode(left)
				root.AddNextNode(right)
				left.AddNextNode(bottom)
				right.AddNextNode(bottom) // Both paths lead to bottom
				return root
			},
			expectedCount: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			root := tt.buildGraph()

			// Act
			result := root.CountNodes()

			// Assert
			if result != tt.expectedCount {
				t.Errorf("CountNodes() = %d; expected %d", result, tt.expectedCount)
			}
		})
	}
}

// TestNodeGraph_CountNodes_Nil tests CountNodes with nil node
// which should return 0
func TestNodeGraph_CountNodes_Nil(t *testing.T) {
	tests := []struct {
		name     string
		node     *NodeGraph[int, string]
		expected int
	}{
		{
			name:     "nil node returns 0",
			node:     nil,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.node.CountNodes()

			// Assert
			if result != tt.expected {
				t.Errorf("CountNodes() on nil = %d; expected %d", result, tt.expected)
			}
		})
	}
}

// TestNodeGraph_HasCycle tests the HasCycle method
// which detects cycles in the graph
func TestNodeGraph_HasCycle(t *testing.T) {
	tests := []struct {
		name       string
		buildGraph func() *NodeGraph[int, string]
		expected   bool
	}{
		{
			name: "acyclic tree has no cycle",
			buildGraph: func() *NodeGraph[int, string] {
				root := NewNodeGraph(1, "root")
				child1 := NewNodeGraph(2, "child1")
				child2 := NewNodeGraph(3, "child2")
				root.AddNextNode(child1)
				root.AddNextNode(child2)
				return root
			},
			expected: false,
		},
		{
			name: "self-loop creates cycle",
			buildGraph: func() *NodeGraph[int, string] {
				root := NewNodeGraph(1, "root")
				root.AddNextNode(root) // Self-loop
				return root
			},
			expected: true,
		},
		{
			name: "two-node cycle",
			buildGraph: func() *NodeGraph[int, string] {
				root := NewNodeGraph(1, "root")
				child := NewNodeGraph(2, "child")
				root.AddNextNode(child)
				child.AddNextNode(root) // Back edge
				return root
			},
			expected: true,
		},
		{
			name: "longer cycle",
			buildGraph: func() *NodeGraph[int, string] {
				root := NewNodeGraph(1, "root")
				child1 := NewNodeGraph(2, "child1")
				child2 := NewNodeGraph(3, "child2")
				child3 := NewNodeGraph(4, "child3")
				root.AddNextNode(child1)
				child1.AddNextNode(child2)
				child2.AddNextNode(child3)
				child3.AddNextNode(root) // Back to root
				return root
			},
			expected: true,
		},
		{
			name: "DAG with shared descendant has no cycle",
			buildGraph: func() *NodeGraph[int, string] {
				root := NewNodeGraph(1, "root")
				left := NewNodeGraph(2, "left")
				right := NewNodeGraph(3, "right")
				bottom := NewNodeGraph(4, "bottom")
				root.AddNextNode(left)
				root.AddNextNode(right)
				left.AddNextNode(bottom)
				right.AddNextNode(bottom) // Diamond pattern, but no cycle
				return root
			},
			expected: false,
		},
		{
			name: "single node with no edges has no cycle",
			buildGraph: func() *NodeGraph[int, string] {
				return NewNodeGraph(1, "single")
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			root := tt.buildGraph()

			// Act
			result := root.HasCycle()

			// Assert
			if result != tt.expected {
				t.Errorf("HasCycle() = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestNodeGraph_HasCycle_Nil tests HasCycle with nil node
// which should return false
func TestNodeGraph_HasCycle_Nil(t *testing.T) {
	tests := []struct {
		name     string
		node     *NodeGraph[int, string]
		expected bool
	}{
		{
			name:     "nil node has no cycle",
			node:     nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := tt.node.HasCycle()

			// Assert
			if result != tt.expected {
				t.Errorf("HasCycle() on nil = %v; expected %v", result, tt.expected)
			}
		})
	}
}

// TestNodeGraph_WithDifferentTypes tests the generic implementation
// with various type combinations
func TestNodeGraph_WithDifferentTypes(t *testing.T) {
	t.Run("with string data and int relation", func(t *testing.T) {
		// Arrange
		root := NewNodeGraph("root", 1)
		child := NewNodeGraph("child", 2)
		root.AddNextNode(child)

		// Act
		count := root.CountNodes()

		// Assert
		if count != 2 {
			t.Errorf("CountNodes() = %d; expected 2", count)
		}
	})

	t.Run("with struct data type", func(t *testing.T) {
		// Arrange
		type Person struct {
			Name string
			Age  int
		}
		root := NewNodeGraph(Person{"Alice", 30}, "parent")
		child := NewNodeGraph(Person{"Bob", 5}, "child")
		root.AddNextNode(child)

		// Act
		degree := root.Degree()

		// Assert
		if degree != 1 {
			t.Errorf("Degree() = %d; expected 1", degree)
		}
	})

	t.Run("with bool data and empty relation", func(t *testing.T) {
		// Arrange
		root := NewNodeGraph(true, "")
		child := NewNodeGraph(false, "")
		root.AddNextNode(child)

		// Act
		isLeaf := root.IsLeaf()

		// Assert
		if isLeaf {
			t.Error("IsLeaf() = true; expected false")
		}
	})
}

// TestNodeGraph_ComplexScenarios tests complex graph operations
// combining multiple methods
func TestNodeGraph_ComplexScenarios(t *testing.T) {
	t.Run("build and modify graph", func(t *testing.T) {
		// Arrange
		root := NewNodeGraph(1, "root")
		child1 := NewNodeGraph(2, "child1")
		child2 := NewNodeGraph(3, "child2")
		child3 := NewNodeGraph(4, "child3")

		// Act: Build graph
		root.AddNextNode(child1)
		root.AddNextNode(child2)
		root.AddNextNode(child3)

		// Assert: Initial state
		if root.Degree() != 3 {
			t.Errorf("Initial Degree() = %d; expected 3", root.Degree())
		}

		// Act: Remove a node
		removed := root.RemoveNextNode(child2)

		// Assert: After removal
		if !removed {
			t.Error("RemoveNextNode() = false; expected true")
		}
		if root.Degree() != 2 {
			t.Errorf("Degree() after removal = %d; expected 2", root.Degree())
		}
		if root.HasNextNode(child2) {
			t.Error("HasNextNode(child2) = true; expected false after removal")
		}
	})

	t.Run("traversal consistency", func(t *testing.T) {
		// Arrange
		root := NewNodeGraph(1, "root")
		for i := 2; i <= 5; i++ {
			root.AddNextNode(NewNodeGraph(i, "child"))
		}

		// Act: Multiple traversals
		var dfsResult, bfsResult []int
		root.TraverseDFS(func(node *NodeGraph[int, string], level int) {
			dfsResult = append(dfsResult, node.Data)
		})
		root.TraverseBFS(func(node *NodeGraph[int, string]) {
			bfsResult = append(bfsResult, node.Data)
		})

		// Assert: Both should visit all nodes
		if len(dfsResult) != len(bfsResult) {
			t.Errorf("DFS visited %d nodes, BFS visited %d; expected same count", len(dfsResult), len(bfsResult))
		}

		// Verify all unique nodes were visited
		sort.Ints(dfsResult)
		sort.Ints(bfsResult)
		if !reflect.DeepEqual(dfsResult, bfsResult) {
			t.Errorf("DFS and BFS visited different sets of nodes")
		}
	})
}
