// tree construction using neighbor joining method
type Tree []*Node
type Node struct{
	neighbor1, neighbor2 *Node  
	sequence string
	distance float64 
}
// high level function of neighbor joining method
func NJ(mtx [][]float64, sequences []string)Tree{
	tree := InitializeTree(mtx, sequences)
	num := len(sequences)
	clusters := InitializeClusters(sequences)
	for p:=num;p<2*num-1;p++{
		row, col := FindMinDist(mtx)
		tree[row].distance = CalcDist(mtx, row, col)
		tree[col].distance = CalcDist(mtx, col, row)
		tree[p].neighbor1 = &tree[row]
		tree[p].neighbor2 = &tree[col]
		// tree[p].sequence = TraceBackSeq(tree[p].neighbor1, tree[p].neighbor2) // need to call function in needleman
		
		// first, add a row and column corresponding to new cluster
		
		mtx = AddRowCol(row, col, clusterSize1, clusterSize2, mtx)
		mtx = DeleteRowCol(mtx, row, col)

		// finally, we clean up clusters
		//add current node to end of our clusters
		
		clusters = append(clusters, tree[p])
		clusters = DeleteClusters(clusters, row, col)
	}
	return tree 
}

// Find the smallest distance after adjust
func FindMinDist(mtx [][]float64)(int, int){
	numRows := len(mtx)
	if numRows < 3{
		return 0,1,mtx[0][1]
	}else{
		AdjMatrix := AdjustMatrix(mtx)
		row := 0
		col := 1
		minVal := AdjMatrix[row][col]

		// range over matrix, and see if we can do better than minVal.
		for i := 0; i < len(AdjMatrix)-1; i++ {
		// start column ranging at i + 1
			for j := i + 1; j < len(AdjMatrix[i]); j++ {
			// do we have a winner?
				if AdjMatrix[i][j] < minVal {
					// update all three variables
					minVal = AdjMatrix[i][j]
					row = i
					col = j
					// col will still always be > row.
				}
			}
		}
		return row, col
	}
}
// Calculate the divergence and the adjusted distance
func CalcDivergence(mtx [][]flaot64, numClusters int)[]float64{
	divergence := make([]float64, len(mtx))
	for row := range mtx{
		sum := 0.0
		for _, val := range mtx[row]{
			sum += val
		}
		div := sum / float64((numClusters - 2))
		divergence[row] = div
	}
	return divergence
}

// Adjust the original matrix to gain the adjusted distances
func AdjustMatrix(mtx [][]float64)[][]float64{
	AdjMatrix := make([][]float64, numRows)
	for row := range mtx{
		AdjMatrix[row] = make([]float64, len(mtx[row]))
	}
	Divergence := CalcDivergence(mtx, numRows)
	for i := range mtx{
		for j := range mtx[i]{
			AdjMatrix[i][j] = mtx[i][j] - Divergence[i] - Divergence[j]
		}
	}
	return AdjMatrix 
}

// Calculate the distance between the two children and their parent
func CalcDist(mtx [][]float64, idx1, idx2 int)float64{
	divergence := CalcDivergence(mtx,len(mtx))
	div1 := divergence[idx1]
	div2 := divergence[idx2]
	return (mtx[idx1][idx2] + div1 - div2) / 2.0
}

// Update the distance matrix by adding a parent row and column

func AddRowCol(row, col int, mtx [][]float64) [][]float64 {
	numRows := len(mtx)
	pRow := make([]float64, numRows+1)

	// calculate the distance between the newly added parent node and the other nodes
	
	for r := 0; r < len(pRow)-1; r++ {
	
		// only set a value of the row if it's not at index row or column
	
		if r != row && r != col {
			pRow[r] = (mtx[r][row] + mtx[r][col] - mtx[row][col]) / 2.0
		}
	}

	// append the newly made row for the parent node to the original matrix
	
	mtx = append(mtx, pRow)

	// also add the values of the parent node distance to other nodes
	
	for c := 0; c < numRows; c++ {
		mtx[c] = append(mtx[c], pRow[c])
	}

	return mtx
}

// Delete the original two nodes replaced by a newly parent node

func DeleteRowCol(mtx [][]float64, row, col int) [][]float64 {

	// col > row, so first delete the col one. Delete two rows first
	
	mtx = append(mtx[:col], mtx[col+1:]...)
	mtx = append(mtx[:row], mtx[row+1:]...)

	// then delete the two columns
	
	for r := range mtx {
		mtx[r] = append(mtx[r][:col], mtx[r][col+1:]...)
		mtx[r] = append(mtx[r][:row], mtx[r][row+1:]...)
	}

	return mtx
}

// Update the clusters that track our disconnected nodes of the tree

func DeleteCluster(clusters []*Node, row, col int)[]*Node{
	clusters = append(clusers[:col], clusters[col+1:])
	clusters = apend(clusters[:row], clusters[row+1:])
	return clusters 
}

// Initialize a tree with leaves representing sequences
// we leave space for parent nodes that would be added later

func InitializeTree(sequences []string)Tree{
	num := len(sequences)
	tree := make(Tree, 2*num-1) 
	for i,seq := range sequences{
		var node Node
		node.sequence = seq
		tree[i] = &node
	}
	return tree  
}

// Initialize clusters of the sequences, which are the leaves of the tree

func (tree Tree) InitializeClusters() []*Node {

	// the tree has 2n-1 total nodes, given the number of leaves is n
	// want the first n nodes of the tree which are leaves

	numNodes := len(tree)
	numLeaves := (numNodes + 1) / 2

	clusters := make([]*Node, numLeaves)

	for i := range clusters {
		clusters[i] = tree[i]
	}

	return clusters
}
