package main 
// tree construction using neighbor joining method
import(
	//"strconv"
)

type Tree []*Node

type Node struct{
	neighbors []*Node  
	sequence Sequences //is a slice of Sequence objects
	distance float64 
}

// high level function of neighbor joining method
func NJ(mtx [][]float64, seqMatrix [][]Sequences, sequences Sequences, subMatrix Matrix)Tree{ //should be the original set of input Sequences
	tree := InitializeTree(sequences)
	num := len(sequences)
	clusters := tree.InitializeClusters() //clusters- a slice of pointers to Node objects
	for p := num; p < 2*num - 1; p++ { //so you start with the first internal node in clusters
		// in the last interation when the matrix size is 2, then we just connect the two node
		// serve the second last node as the root 
		if len(mtx) == 2{
			tree[p].distance = mtx[0][1] / 2.0
			tree[p].neighbors = append(tree[p].neighbors, clusters[0])
			tree[p].neighbors = append(tree[p].neighbors, clusters[1])
			//tree[p].sequence = seqMatrix[0][1]
			tree[p].sequence, _ = NeedlemanWunsch(clusters[0].sequence, clusters[1].sequence, subMatrix)
		}else{
			row, col := FindMinDist(mtx) //the indices of this should be the same corresponding to the respective Sequence in Sequences object
			tree[p].neighbors = append(tree[p].neighbors, clusters[row])
			tree[p].neighbors = append(tree[p].neighbors, clusters[col])
			tree[p].neighbors[0].distance = CalcDist(mtx, row, col)
			tree[p].neighbors[1].distance = CalcDist(mtx, col, row)
			//if the new neighbors/children have a sequences of length 1, then we already got the alignment when we built the distance matrix, so look it up
			if len(clusters[row].sequence) == 1 && len(clusters[col].sequence) == 1 { 
				tree[p].sequence = seqMatrix[row][col]
			} else { //you need to get the alignment
				tree[p].sequence, _ = NeedlemanWunsch(clusters[row].sequence, clusters[col].sequence, subMatrix)
			}

			// tree[p].sequence = TraceBackSeq(tree[p].neighbor1, tree[p].neighbor2) // need to call function in needleman
			
			// first, add a row and column corresponding to new cluster
			mtx = AddRowCol(row, col, mtx)
			mtx = DeleteRowCol(mtx, row, col)

			// finally, we clean up clusters
			//add current node to end of our clusters
			clusters = append(clusters, tree[p])
			clusters = DeleteClusters(clusters, row, col)
		}
		
	}
	return tree 
}

// Find the smallest distance after adjust
//returns the row and col indices in the matrix of the minimum value
func FindMinDist(mtx [][]float64)(int, int){
	numRows := len(mtx)
	if numRows < 3{
		return 0,1
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
func CalcDivergence(mtx [][]float64)[]float64{
	divergence := make([]float64, len(mtx))
	numClusters := len(mtx)
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
	numRows := len(mtx)
	AdjMatrix := make([][]float64, numRows)
	for row := range mtx{
		AdjMatrix[row] = make([]float64, len(mtx[row]))
	}
	Divergence := CalcDivergence(mtx)
	for i := range mtx{
		for j := range mtx[i]{
			if i != j{
				AdjMatrix[i][j] = mtx[i][j] - Divergence[i] - Divergence[j]
			}else{
				AdjMatrix[i][j] = mtx[i][j]
			}
			
		}
	}
	return AdjMatrix 
}

// Calculate the distance between the two children and their parent
func CalcDist(mtx [][]float64, idx1, idx2 int)float64{
	divergence := CalcDivergence(mtx)
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

func DeleteClusters(clusters []*Node, row, col int)[]*Node{
	clusters = append(clusters[:col], clusters[col+1:]...)
	clusters = append(clusters[:row], clusters[row+1:]...)
	return clusters 
}

// Initialize a tree with leaves representing sequences
// we leave space for parent nodes that would be added later
func InitializeTree(sequences Sequences)Tree{
	var tree Tree 
	num := len(sequences)
	tree = make([]*Node, 2*num-1) 

	for i := range tree{
		var node Node
		node.neighbors = make([]*Node,0)
		if i < num {
			node.sequence = sequences[i:i+1] //needs to be a slice of Sequences of length 1 in leaves?
		}else{
			// node.sequence = "aligned sequence" + strconv.Itoa(i-num)
			node.sequence = make(Sequences, 0) //this works right?
		}
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

func IstheSame(mtx1,mtx2 [][]float64)bool{
	if len(mtx1) != len(mtx2) || len(mtx1[0]) != len(mtx2[0]){
		panic("unmatched length of the two matrices")
	}
	result := true
	for row := range mtx1{
		for col := range mtx1[0]{
			if mtx1[row][col] != mtx2[row][col]{
				result = false
			}
		}
	}
	return result 
}