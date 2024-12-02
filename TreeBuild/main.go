package main
import(
	"fmt"
	/*"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/encoding/dot"
	"gonum.org/v1/gonum/simple"
	"gonum.org/v1/gonum/traverse"
	"log"
	"os"*/
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)
func main() {
	mtx := [][]float64{
		[]float64{0, 4, 5, 10},
		[]float64{4, 0, 7, 12},
		[]float64{5, 7, 0, 9},
		[]float64{10, 12, 9, 0},
	}
	seq := []string{"ATGC", "TGAC", "CGTA", "GGCA"}
	tree := NJ(mtx, seq)
	root := tree[len(tree)-1]
	PrintTree(root, "", false)

	// Plot the tree
	PlotTree(root)
}
	/*
	 g := simple.NewUndirectedGraph()
	 nodes := map[int64]graph.Node{}
	 for id, node := range tree {
		 nodes[id] = g.NewNode()
		 g.AddNode(nodes[id])
	 }
 
	 for id, node := range tree {
		 for _, neighbor := range node.Neighbors {
			 g.SetEdge(g.NewEdge(nodes[id], nodes[neighbor]))
		 }
	 }
 
	 // Output to DOT format for rendering
	 dotData, err := dot.Marshal(g, "Tree", "", "  ")
	 if err != nil {
		 log.Fatalf("failed to create DOT data: %v", err)
	 }
 
	 // Save to file for visualization
	 outputFile := "tree.dot"
	 err = os.WriteFile(outputFile, []byte(dotData), 0644)
	 if err != nil {
		 log.Fatalf("failed to save DOT data: %v", err)
	 }
 
	 log.Printf("Tree visualization saved to %s", outputFile)
}
var PositionMap = make(map[*Node]plotter.XY)

// PrintTree recursively prints the tree structure
func PrintTree(node *Node, prefix string, isLeft bool) {
	if node == nil {
		return
	}

	// Print the current node
	if isLeft {
		fmt.Printf("%s├── ", prefix)
	} else {
		fmt.Printf("%s└── ", prefix)
	}

	if node.sequence != "" {
		fmt.Printf("Leaf: %s\n", node.sequence)
	} else {
		fmt.Printf("Node: %.2f\n", node.distance)
	}

	// Update prefix and print child nodes
	childPrefix := prefix
	if isLeft {
		childPrefix += "│   "
	} else {
		childPrefix += "    "
	}

	for i, neighbor := range node.neighbors {
		PrintTree(neighbor, childPrefix, i == 0)
	}
}

// AssignPositions assigns positions to each node based on index and depth
func AssignPositions(node *Node, index int, depth int) {
	if node == nil {
		return
	}

	// Assign position based on the index and depth
	PositionMap[node] = plotter.XY{X: float64(index), Y: float64(depth)}

	// Recursively assign positions to neighbors
	for _, neighbor := range node.neighbors {
		AssignPositions(neighbor, index+1, depth+1)
		index++
	}
}

// AddNodes adds nodes to the plot based on the external position map
func addNodes(node *Node, nodes *plotter.Scatter, lines *plotter.Line) {
	if node == nil {
		return
	}

	// Retrieve the pre-calculated position for this node
	pos, exists := PositionMap[node]
	if exists {
		// Add the current node's position to the scatter plot
		nodes.XYs = append(nodes.XYs, pos)

		// Add lines for connections to the neighbors
		for _, neighbor := range node.neighbors {
			posNeighbor := PositionMap[neighbor]
			lines.XYs = append(lines.XYs, plotter.XY{X: pos.X, Y: pos.Y})
			lines.XYs = append(lines.XYs, posNeighbor)
			addNodes(neighbor, nodes, lines)
		}
	}
}

// PlotTree visualizes the tree
func PlotTree(root *Node) {
	// Initialize a new plot
	p := plot.New()
	p.Title.Text = "Phylogenetic Tree"

	// Create scatter plot for nodes and line plot for edges
	nodes, err := plotter.NewScatter(plotter.XYs{})
	if err != nil {
		panic(err)
	}

	lines, err := plotter.NewLine(plotter.XYs{})
	if err != nil {
		panic(err)
	}

	// Assign positions based on the tree structure
	AssignPositions(root, 0, 0)

	// Add nodes and edges based on positions
	addNodes(root, nodes, lines)

	// Add plots to the chart
	p.Add(nodes, lines)

	// Save the plot to a file
	if err := p.Save(8*vg.Inch, 6*vg.Inch, "tree_with_dot_line.png"); err != nil {
		panic(err)
	}
}
*/
var PositionMap = make(map[*Node]plotter.XY)

// PrintTree recursively prints the tree structure
func PrintTree(node *Node, prefix string, isLeft bool) {
	if node == nil {
		return
	}

	// Print the current node
	if isLeft {
		fmt.Printf("%s├── ", prefix)
	} else {
		fmt.Printf("%s└── ", prefix)
	}

	if node.sequence != "" {
		fmt.Printf("%s, %.2f\n", node.sequence, node.distance)
	} // else {
		//fmt.Printf("Node: %.2f\n", node.distance)
	

	// Update prefix and print child nodes
	childPrefix := prefix
	if isLeft {
		childPrefix += "│   "
	} else {
		childPrefix += "    "
	}

	for i, neighbor := range node.neighbors {
		PrintTree(neighbor, childPrefix, i == 0)
	}
}

// AssignPositions assigns positions to each node based on index and depth
func AssignPositions(node *Node, index int, depth int) {
	if node == nil {
		return
	}

	// Assign position based on the index and depth
	PositionMap[node] = plotter.XY{X: float64(index), Y: float64(depth)}

	// Recursively assign positions to neighbors
	for _, neighbor := range node.neighbors {
		AssignPositions(neighbor, index+1, depth+1)
		index++
	}
}

// AddNodes adds nodes to the plot based on the external position map
func addNodes(node *Node, nodes *plotter.Scatter, lines *plotter.Line) {
	if node == nil {
		return
	}

	// Retrieve the pre-calculated position for this node
	pos, exists := PositionMap[node]
	if exists {
		// Add the current node's position to the scatter plot
		nodes.XYs = append(nodes.XYs, pos)

		// Add lines for connections to the neighbors
		for _, neighbor := range node.neighbors {
			posNeighbor := PositionMap[neighbor]
			lines.XYs = append(lines.XYs, plotter.XY{X: pos.X, Y: pos.Y})
			lines.XYs = append(lines.XYs, posNeighbor)
			addNodes(neighbor, nodes, lines)
		}
	}
}

// PlotTree visualizes the tree
func PlotTree(root *Node) {
	// Initialize a new plot
	p := plot.New()
	p.Title.Text = "Phylogenetic Tree"

	// Create scatter plot for nodes and line plot for edges
	nodes, err := plotter.NewScatter(plotter.XYs{})
	if err != nil {
		panic(err)
	}

	lines, err := plotter.NewLine(plotter.XYs{})
	if err != nil {
		panic(err)
	}

	// Assign positions based on the tree structure
	AssignPositions(root, 0, 0)

	// Add nodes and edges based on positions
	addNodes(root, nodes, lines)

	// Add plots to the chart
	p.Add(nodes, lines)

	// Save the plot to a file
	if err := p.Save(8*vg.Inch, 6*vg.Inch, "tree_with_dot_line.png"); err != nil {
		panic(err)
	}
}
