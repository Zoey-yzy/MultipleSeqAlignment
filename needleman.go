package main
import "fmt"

type  Move struct {
    direction string
    value float64
}

func GetScoreMatrix2D(seq1 string, seq2 string )[][]float64{

    n1:=len(seq1) + 1 //set number of rows in score matrix
    n2:=len(seq2) + 1 //set number of columns in score matrix
    matrix := make([][]float64,n1)
    for i:=0;i<n1;i++{
        matrix[i] = make([]float64,n2)

    }
    for i:=0;i<n1;i++{ //assigns initial rows
        matrix[i][0] = -1.0 * float64(i) //this should be generalized for whatever the gap penalty is

    }
    for j:=0;j<n2;j++{ //assigns initial cols
        matrix[0][j] = -1.0 * float64(j)

    }
    return matrix
}
func GetTracebackMatrix2D(seq1 string, seq2 string )[][]string{

    n1:=len(seq1) + 1 //this could be common code with score matrix
    n2:=len(seq2) + 1
    matrix := make([][]string,n1) //can we do a matrix of tuples/unit ordered pairs?
    for i:=0;i<n1;i++{
        matrix[i] = make([]string,n2)

    }
    // this is trivial,once we reach the boundary row, we wont need to check direction
    for i:=0;i<n1;i++{
        matrix[i][0] = "LEFT" //assigns the first column to be "LEFT"- don't think this is correct

    }
    // this is trivial,once we reach the boundary column, we wont need to check direction
    for j:=0;j<n2;j++{
        matrix[0][j] = "UP" //assigns the first row to be "UP"- again, don't think this is correct

    }
    PrintTracebackMatrix(matrix)
    return matrix
}

func GetScoreMap(gapPenalty,matchReward,misMatchPenalty float64)map[string](map[string]float64){ //do we need this?
    score := make(map[string](map[string]float64))
    /*
    Assume matchReward = 1, misMatch penalty = -1
     A T G C
    A 1 -1 -1 -1
    T -1 1 -1 -1
    G -1 -1 1 -1
    C -1 -1 -1 1
    
    
    */
    alphabets := [8]string{"A","T","G","C","a","t","g","c"}
    for i:=range(alphabets){
        score[alphabets[i]] = make(map[string]float64)
        for j:=range(alphabets){
            if alphabets[i] == alphabets[j]{
                score[alphabets[i]][alphabets[j]] = matchReward
            }else{
                score[alphabets[i]][alphabets[j]] = misMatchPenalty

            }
        }
    }
    score["gap"] = make(map[string]float64)
    score["gap"]["gap"]  = gapPenalty
    return score
}
func GetMaxValue(list [3]float64) float64{ //gets max value from a set of three

    //help get max value
    max := list[0]
    for i:=0;i<len(list);i++{
        if list[i] > max{
            max = list[i]
        }
    }
    return max
}
func GetMaxMove(matrix [][]float64,i,j int,c1,c2 string,scoreMap map[string](map[string]float64)) Move{
    /*
    for a given cell i,j what is the move that will get me the maximum score given cells to my left (i-1,j)
    up (i,j-1) and diagonal (i-1,j-1)
    */
    diagScore := matrix[i-1][j-1] + scoreMap[c1][c2]
    leftScore := matrix[i][j-1] + scoreMap["gap"]["gap"]
    upScore := matrix[i-1][j] + scoreMap["gap"]["gap"]
    max := GetMaxValue([3]float64{diagScore,leftScore,upScore})
    var move Move
    move.value = max
    // checking which move gives me the max score
    if max == diagScore{
        move.direction = "DIAG"

    }else if max == leftScore{
        move.direction = "LEFT"

    }else{
        move.direction =  "UP"

    }
    return move
}
func ComputeScores(matrix [][]float64,traceBackMatrix [][]string,seq1 string,seq2 string,scoreMap map[string](map[string]float64))([][]float64,[][]string){
    s1 := len(matrix)
    s2 := len(matrix[0])
    for i:=1;i<s1; i++ {
        for j:=1;j<s2; j++{
            move := GetMaxMove(matrix,i,j,string(seq1[i-1]),string(seq2[j-1]),scoreMap)
            matrix[i][j] = move.value
            traceBackMatrix[i][j] = move.direction
        }

    }
    PrintMatrix(matrix)
    PrintTracebackMatrix(traceBackMatrix)
    // fmt.Println("matrix:",matrix)
    // fmt.Println("traceBackMatrix:",traceBackMatrix)
    return matrix,traceBackMatrix

}
func PairwiseAligment(seq1 string,seq2 string) (float64, string, string){
    gapPenalty := -1.0 //not hardcode this?
    matchReward:= 1.0
    misMatchPenalty := -1.0
    s1,s2 := len(seq1)+1,len(seq2)+1
    matrix := GetScoreMatrix2D(seq1, seq2)
    traceBackMatrix := GetTracebackMatrix2D(seq1,seq2)
    scoreMap := GetScoreMap(gapPenalty,matchReward,misMatchPenalty)
    fmt.Println("matrix:",matrix)
    fmt.Println("scoreMap:",scoreMap)
    
    
    matrix,traceBackMatrix = ComputeScores(matrix,traceBackMatrix,seq1,seq2,scoreMap)
    PrintMatrix(matrix)
    PrintTracebackMatrix(traceBackMatrix)
    // fmt.Println("matrix:",matrix)
    // fmt.Println("traceBackMatrix:",traceBackMatrix)
    align1, align2 := TraceSequences(seq1, seq2, traceBackMatrix)
    fmt.Println("align1", align1, "align2", align2)
    alignmentScore := matrix[s1-1][s2-1]
    return alignmentScore, align1, align2

}

//TraceSequences uses the traceback matrix to assemble a pairwise-aligned seq1 and seq2
//input: original sequences strings seq1, seq2; traceBackMatrix
//returns: align1, align2 strings incorporating gaps as '-' where appropriate
func TraceSequences(seq1, seq2 string, traceBackMatrix [][]string) (string, string) {
    fmt.Println("aligning sequences:\n", seq1, "\n", seq2)
    // length := MaxLength(seq1, seq2)
    align1, align2 := make([]byte, 0), make([]byte, 0) //make empty strings (slices of bytes) to store aligned sequences
    ind1, ind2 := len(seq1), len(seq2) //start tracing back at the end of each sequence
    //elements of seq1 are the row labels, elements of seq2 are the column labels
    for ind1 > 0 && ind2 > 0 {
        if traceBackMatrix[ind1][ind2] == "DIAG" { //getting the index of previous step
            fmt.Println("ind1", ind1, "ind2", ind2)
            align1 = append([]byte{byte(seq1[ind1-1])}, align1...) //pre-appends a slice of bytes of length 1 to the existing string (slice of bytes)
            align2 = append([]byte{byte(seq2[ind2-1])}, align2...)
            ind1--
            ind2-- //index is one lower next time around
            fmt.Println("align1", string(align1))
            fmt.Println("align2", string(align2))
        } else if traceBackMatrix[ind1][ind2] == "UP" {
            fmt.Println("ind1", ind1, "ind2", ind2)
            align1 = append([]byte{byte(seq1[ind1-1])}, align1...)
            align2 = append([]byte{'-'}, align2...) //'-' is gap symbol
            ind1-- //don't change the index indicating the horizontal thing
            fmt.Println("align1", string(align1))
            fmt.Println("align2", string(align2))
        } else { //if traceBackMatrix[ind1][ind2] == "LEFT" {
            fmt.Println("ind1", ind1, "ind2", ind2)
            align1 = append([]byte{'-'}, align1...)
            align2 = append([]byte{byte(seq2[ind2-1])}, align2...)
            ind2--
            fmt.Println("align1", string(align1))
            fmt.Println("align2", string(align2))
        }
    }
    return string(align1), string(align2)
}

// func MaxLength(seq1, seq2 string) int {
//     if len(seq1) > len(seq2) {
//         return len(seq1)
//     } else {
//         return len(seq2)
//     }
// }

func PrintMatrix(matrix [][]float64) {
    for i := range matrix {
        fmt.Println(matrix[i])
    }
}

func PrintTracebackMatrix(matrix [][]string) {
    for i := range matrix {
        fmt.Println(matrix[i])
    }
}

func main(){
    seq1 := "atatat"
    seq2 := "att"
    scores, align1, align2 := PairwiseAligment(seq1,seq2)
    fmt.Println("Score:",scores)
    fmt.Println(align1)
    fmt.Println(align2)

}