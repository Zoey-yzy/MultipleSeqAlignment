// tree construction using neighbor joining method
type Tree []*Node
type Node struct{
	neighbor1, neighbor2 *Node  
	sequence string 
}
// high level function of neighbor joining method
func NJ(mtx [][]float64, sequences []string)Tree{
	tree := InitializeTree(mtx, sequences)
	num := len(sequences)
	row, col, sum := FindMinSum(mtx)
	for p:=num;p<2*num-1;p++{
		tree[p].neighbor1 = &tree[row]
		tree[p].neighbor2 = &tree[col]
		tree[p].sequence = TraceBackSeq(tree[p].neighbor1, tree[p].neighbor2) // need to call function in needleman
		mtx = UpdateDistance(mtx,row,col)
	}
	return tree 
}
// initialize tree of all sequences
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
// each time choose a pair of neighbors and update the distance matrix
func UpdateDistance(mtx [][]float64, row, col int)[][]float64{
	n := len(mtx)
	NewRow := make([]float64, n)
	for i:=0;i<n;i++{
		if i == row || i == col{
			NewRow[i] = 0
		}else{
			NewRow[i] = (mtx[i][row] + mtx[i][col])/2.0 
		}
	}
	mtx = append(mtx,NewRow)
	for j:=0;j<n+1;j++{
		mtx[j] = append(mtx,NewRow[j])
	}
	mtx = append(mtx[:col], mtx[col+1:]...)
	mtx = append(mtx[:row], mtx[row+1:]...)
	for r := range mtx{
		mtx[r] = append(mtx[r][:col], mtx[r][col+1:]...)
		mtx[r] = append(mtx[r][:row], mtx[r][row+1:]...)
	}
	return mtx 
}
// to get the pair of nodes (row and col) that we need to add a branch
func FindMinSum(mtx [][]float64)(int,int,float64){
	n := len(mtx)
	row := 0
	col := 1
	sum := CalculateSum(mtx,row,col)
	for i:=0;i<n-1;i++{
		for j:=i+1;j<n;j++{
			if sum > CalculateSum(mtx,i,j){
				sum = CalculateSum
				row = i
				col = j
			}
		}
	}
	return row, col, sum 
}
// to calculate the sum of distance for each chosen pair of nodes
func CalculateSum(mtx [][]float64, row, col int)float64{
	n := len(mtx)
	sum1, sum3 := 0.0, 0.0
	for i:=0;i<n-1;i++{
		for j:=i+1;j<n;j++{
			if i!=row && i!=col && j!=row && j!=col{
				sum += mtx[i][j]
			}
		}
	}
	// the sum of other n-2 nodes inside the larger subtree
	sum3 /= (n-2) 
	// the distance of the chosen two nodes
	sum2 := 0.5 * mtx[row][col] 
	// the distance of the two nodes from the smaller subtree to the other nodes from the other subtree
	for k:=0;k<n;k++{
		sum1 += mtx[row][k] + mtx[col][k]
	}
	sum1 /= (2*(n-2))
	return sum1 + sum2 + sum3 
}