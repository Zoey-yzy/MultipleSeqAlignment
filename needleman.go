package main
import "fmt"
//import "os"
//import "strings"
//import "bufio"
type  Move struct {
    direction string
    value float64
}

//the function to do Needleman repeatedly- every sequence against every other sequence, to create a distance matrix
//input: inputSeqs Sequences, subMatrix Matrix
//output: a distanceMatrix Matrix and seqMatrix SequenceMatrix
func NWEverybody(inputSeqs Sequences, subMatrix Matrix) ([][]int, [][]Sequences){
    //code for multiple NWs on multiple sequences - this should be modularized
    // var distanceMatrix Matrix 
    // distanceMatrix := make(Matrix) 
    distanceMatrix := make([][]int, len(inputSeqs)) //hold int scores- results of pairwise alignments
    seqMatrix := make([][]Sequences, len(inputSeqs)) //hold Sequences objects of length 2- no need to redo alignments when filling guide tree nodes, just fill from here
    // var seqMatrix SequenceMatrix 
    for i := 0; i < len(inputSeqs); i++ {
        distanceMatrix[i] = make([]int, len(inputSeqs))
        seqMatrix[i] = make([]Sequences, len(inputSeqs))
        for j := i+1; j < len(inputSeqs); j++ {
            // fmt.Println("printing", inputSeqs[i:i+1], inputSeqs[j:j+1])
            alignedTwo, score := NeedlemanWunsch(inputSeqs[i:i+1], inputSeqs[j:j+1], subMatrix) //returns a Sequences object- the aligned two sequences, int the alignment score
            PrintSequencesList(alignedTwo)
            fmt.Println("alignment score is", score)
            distanceMatrix[i][j] = score
            seqMatrix[i][j] = alignedTwo
            // distanceMatrix.UpdateDistMatrix(inputSeqs[i].info, inputSeqs[j].info, score)
            // innerMap := make(map[string]int)
            // fmt.Println("score passed into update function is", score)
            // innerMap[inputSeqs[j].info] = score
            // distanceMatrix[inputSeqs[i].info] = innerMap //why is the value 0?
            // fmt.Println("distance matrix looks like: ")
            fmt.Println(distanceMatrix)
            // PrintSubMatrix(distanceMatrix)
            // distanceMatrix[inputSeqs[i].info][inputSeqs[j].info] = score
            // seqMatrix[inputSeqs[i].info][inputSeqs[j].info] = alignedTwo
        }
    }
    return distanceMatrix, seqMatrix
}

//this is the highest-level needleman algorithm which will be used to align two sequences or two sets of sequences
//Input: the two sequences- rowSeqs and colSeqs, a scoring matrix subMatrix type Matrix
//Output: a combined, aligned Sequences object, as well as the aligment score
func NeedlemanWunsch(rowSeqs, colSeqs Sequences, subMatrix Matrix) (Sequences, int){
    // gapPenalty := -1.0 //how do i deal with you?- it's part of the reading-in function now
    // matchReward:= 1.0
    // misMatchPenalty := -1.0
    // s1,s2 := len(seq1)+1,len(seq2)+1
    
    //rowSeqs is going to be the row labels- the sequences going vertically down the side of matrix
    //colSeqs is going to be the col labels- the sequences going horizontally across the top of matrix

    //add a gap symbol to the start of every sequence, then remove it before you exit this function
    // fmt.Println("length is", len(rowSeqs[0].sequence))
    // fmt.Println("sequence:", rowSeqs[0].sequence)
    for i := 0; i < len(rowSeqs); i++ {
        rowSeqs[i].sequence = "-" + rowSeqs[i].sequence
    }
    for i := 0; i < len(colSeqs); i++ {
        colSeqs[i].sequence = "-" + colSeqs[i].sequence
    }
    // fmt.Println("length is", len(rowSeqs[0].sequence))
    // fmt.Println("sequence:", rowSeqs[0].sequence)
    rowLen, colLen := len(rowSeqs[0].sequence), len(colSeqs[0].sequence)
    //if 1 sequence in row, col, then it's obviously the max, if multiple, they should be aligned and all lengths are equal
    //this is the scoring matrix e.g what is the score for a match,mismatch,gap etc
    // scoreMap := ComputeNeedlemanScoreMap(gapPenalty,matchReward,misMatchPenalty) // this initializes the scoring mechanism
    //should change this part to take in a "scoring lookup table" to get match/mismatch scores from a .csv file in main.go or io.go
    //i have a scorematrix now
    matrix, traceBackMatrix := ComputeNWScores(rowSeqs, colSeqs, subMatrix)
    //fmt.Println("matrix:",matrix)
    //fmt.Println("scoreMap:",scoreMap)
    //PrettyPrintMatrix(traceBackMatrix)
    // for i := 0; i < len(matrix); i++ {
    //     fmt.Println(matrix[i])
    // }
    alignedSeqs := GetAlignments(rowSeqs, colSeqs, traceBackMatrix)
    
    alignmentScore := matrix[rowLen-1][colLen-1] //the score taken from bottom corner is the alignment score- use this to build difference matrix

    //remove the gap symbol from the start of the alignedSeqs
    // seq = "-wefaweioafj"
    // seq = seq[1:]
    for i := 0; i < len(rowSeqs); i++ {
        rowSeqs[i].sequence = rowSeqs[i].sequence[1:]
    }
    for i := 0; i < len(colSeqs); i++ {
        colSeqs[i].sequence = colSeqs[i].sequence[1:]
    }

    return alignedSeqs, alignmentScore
}

//this function will compute the aligned sequences using the traceback matrix
//Input: Sequences rowSeqs and colSeqs, the traceback matrix which has the the direction taken during the Needleman algorithm
//Output: a Sequences object, consisting of rowSeqs and colSeqs sequences, aligned
func GetAlignments(rowSeqs, colSeqs Sequences, traceBackMatrix [][]string) Sequences{
    n1 := len(rowSeqs[0].sequence) //rowSeq might be diff length from colSeq but all seqs in each are the same length
    n2 := len(colSeqs[0].sequence)//no +1 because gap symbol already added to 1st position- sequence[0] is always "-" which you can use
    fmt.Println("length is", len(rowSeqs[0].sequence))
    //make a slice of sequences (strings) for rowSeq, colSeq, equal to how many there are, set them each to ""
    newRowSeqs := make([]Sequence, len(rowSeqs))
    newColSeqs := make([]Sequence, len(colSeqs)) //make new Sequences objects
    newRowSeq := make([]string, len(rowSeqs))
    newColSeq := make([]string, len(colSeqs)) //will have default value of ""

    //fill each one with the appropriate sequence info
    for i := 0; i < len(newRowSeqs); i++ {
        newRowSeqs[i].info = rowSeqs[i].info
        newRowSeqs[i].sequence = ""
    }
    for i := 0; i < len(newColSeqs); i++ {
        newColSeqs[i].info = colSeqs[i].info
        newColSeqs[i].sequence = ""
    }

    // newSeq1 := ""
    // newSeq2 := ""

    i := n1 - 1
    j := n2 - 1 //i and j are the index-character you're on
    direction := traceBackMatrix[i][j] //start at max row and col of matrix
    //i reflects position of row labels; symbols going down side of matrix
    //j reflects position of column labels; symbols going across top of matrix
    
    for i > 0 || j > 0 {
        //fmt.Println("Direction:",direction,"i:","j", i,j,"seq1:",newSeq1,"seq2:",newSeq2)
        if direction == "DIAG"{ //add the character at that position to everybody, decrement i and j by 1
            if i > 0 {
                for m := 0; m < len(newRowSeq); m++ {
                    newRowSeq[m] = string(rowSeqs[m].sequence[i]) + newRowSeq[m]
                    // fmt.Println(string(rowSeqs[m].sequence[i]))
                }
                // newSeq1 =  string(seq1[i-1]) + newSeq1
            }
            if j > 0 {
                for n := 0; n < len(newColSeq); n++ {
                    newColSeq[n] = string(colSeqs[n].sequence[j]) + newColSeq[n]
                }
                // newSeq2 =  string(seq2[j-1]) + newSeq2
            }         
            i--
            j--
        } else if direction == "LEFT"{ //add character at that position to column sequence, "-" to row; decrement the column count by 1
            for m := 0; m < len(newRowSeq); m++ {
                newRowSeq[m] = "-" + newRowSeq[m]
            }
            // newSeq1 =  "-" + newSeq1
            if j > 0 {
                for n := 0; n < len(newColSeq); n++ {
                    newColSeq[n] = string(colSeqs[n].sequence[j]) + newColSeq[n]
                }
                // newSeq2 =  string(seq2[j-1]) + newSeq2
            }
            j--
        } else if direction == "UP"{ //add character at that position to row sequence, "-" to column; decrement the row count by 1
            if i > 0 {
                for m := 0; m < len(newRowSeq); m++ {
                    newRowSeq[m] = string(rowSeqs[m].sequence[i]) + newRowSeq[m]
                }
                // newSeq1 =  string(seq1[i-1]) + newSeq1
            }
            for n := 0; n < len(newColSeq); n++ {
                newColSeq[n] = "-" + newColSeq[n]
            }
            // newSeq2 =  "-"+ newSeq2
            i-- 
        }
        // fmt.Println(newRowSeqs[0], newColSeqs[0])
        direction = traceBackMatrix[i][j]    
    }

    //now, replace the old sequences in rowSeqs and colSeqs with the sequences in new ones
    for i := 0; i < len(rowSeqs); i++ {
        newRowSeqs[i].sequence = newRowSeq[i]
    }
    for j := 0; j < len(colSeqs); j++ {
        newColSeqs[j].sequence = newColSeq[j]
    }
    //combine rowSeqs and colSeqs into one sequences object
    mergedSeqs := append(newRowSeqs, newColSeqs...)

    return mergedSeqs
}

//this function is like a forward pass of the needleman algorthm whereby we compute the score for a given subsequence ending at i,j
//Input: the input Sequences (single or multiple aligned) rowSeqs and colSeqs, the scoring dictionary subMatrix
//Output: two matrices, the first matrix represents a matrix of alignment scores whereas the second matrix consists of the direction of travel in the forward pass
// of needleman algorithm e.g which direcion we took from a given cell i,j
func ComputeNWScores(rowSeqs, colSeqs Sequences, subMatrix Matrix) ([][]int, [][]string){
    matrix := InitializeNWMatrix2D(rowSeqs, colSeqs, subMatrix)
    traceBackMatrix := InitializeNWTracebackMatrix2D(rowSeqs, colSeqs)
    s1 := len(matrix)
    s2 := len(matrix[0])
    for i := 1; i < s1; i++ {
        for j := 1; j < s2; j++{
            value, direction := ComputeNWMaxMove(matrix, rowSeqs, colSeqs, i, j, subMatrix)
            matrix[i][j] = value
            traceBackMatrix[i][j] = direction
            
            // move := ComputeNeedlemanMaxMove(matrix, rowSeqs, colSeqs, rowInd, colInd, subMatrix)
            // matrix[i][j] = move.value
            // traceBackMatrix[i][j] = move.direction
        }
    }
    //fmt.Println("matrix:",matrix)
    //fmt.Println("traceBackMatrix:",traceBackMatrix)
    return matrix, traceBackMatrix
}

//This function is used to initialize the first row and first column
//with gap penalties
//Input: two seq1 and seq2 string
//Output: a 2D matrix
func InitializeNWMatrix2D(rowSeqs, colSeqs Sequences, subMatrix Matrix) [][]int {
    n1 := len(rowSeqs[0].sequence) 
    n2 := len(colSeqs[0].sequence)//no +1 because gap symbol already added to 1st position- sequence[0] is always "-" which you can use
    matrix := make([][]int, n1)
    for i := 0; i < n1; i++ {
        matrix[i] = make([]int, n2)
    } //made a matrix to hold scores- if all match/mismatch/gap are ints, couldn't this be int?

    //if you're SumOfPairs of gap for "UP" or "LEFT" use index 0 for the respective index
    matrix[0][0] = SumOfPairs(rowSeqs, colSeqs, 0, 0, subMatrix)
    for i := 1; i < n1; i++ {
        matrix[i][0] = SumOfPairs(rowSeqs, colSeqs, i, 0, subMatrix) + matrix[i-1][0] //the only neighbor is to your left
        // matrix[i][0] = gapPenalty * float64(i)
    }
    
    for j := 1; j < n2; j++ {
        matrix[0][j] = SumOfPairs(rowSeqs, colSeqs, 0, j, subMatrix) + matrix[0][j-1] //the only neighbor is above you
    }
    return matrix
}

//SumOfPairs
//input: rowSeqs, colSeqs Sequences, rowInd and colInd the index of the letter we're looking at, subMatrix Matrix for looking up scores
//return: the integer score of resulting from performing sum of pairs operation at this index
func SumOfPairs(rowSeqs, colSeqs Sequences, rowInd, colInd int, subMatrix Matrix) int{
    sopScore := 0
    //for every seq in rowSeqs
    for i:= 0; i < len(rowSeqs); i++ {
        for j := 0; j < len(colSeqs); j++ { //i and j are indexes of the lists of sequences on the row and on the column
            // fmt.Println("i, j", i, j)
            // fmt.Println("rowInd, colInd, rowSeqs", rowInd, colInd, len(rowSeqs))
            symbol1 := string(rowSeqs[i].sequence[rowInd]) //rowInd indicates which letter in the sequences to use
            symbol2 := string(colSeqs[j].sequence[colInd])
            score := subMatrix[symbol1][symbol2]
            sopScore += score
        }
    }//if you only have 1 sequence in rowSeqs, colSeqs, it just returns the score for that single matchup- it works
    return sopScore
}

//this function will intialize the traceback matrix
//Input: rowSeqs, colSeqs Sequences- just lengths
//Output: a traceback matrix- 2D slice of strings
func InitializeNWTracebackMatrix2D(rowSeqs, colSeqs Sequences) [][]string {
    n1 := len(rowSeqs[0].sequence) 
    n2 := len(colSeqs[0].sequence)//no +1 because gap symbol already added to 1st position- sequence[0] is always "-" which you can use
    matrix := make([][]string, n1)
    for i := 0; i < n1; i++ {
        matrix[i] = make([]string, n2)
    } //made a matrix of scores

    //does 0 0 need a direction associated with it?

    // this is trivial,once we reach the boundary row, we wont need to check direction
    //fill first col
    for i := 1; i < n1; i++ {
        matrix[i][0] = "UP"
    }
    // this is trivial,once we reach the boundary column, we wont need to check direction
    //fill first row
    for j := 0; j < n2; j++ {
        matrix[0][j] = "LEFT"
    }
    return matrix
}

//this is a helper function that essentially helps us find the maximum value in ta 3 element list
//Input: a list of three ints
//Output: an int which is the maximum in the input list
func ComputeNWMax(list [3]int) int{
    //help get max value
    max := list[0]
    for i := 0; i < len(list); i++ {
        if list[i] > max{
            max = list[i]
        }
    }
    return max
}
//this is a helper function to compute the maximum value for three possible cells: UP,DIAGONAL,LEFT
//Input: a 2D matrix ints containing scores, rowSeqs and colSeqs the Sequences being aligned, the current indices rowInd, colInd to look at and a Matrix submatrix specifying the scoring mechanism
//Output: int value for the maximum score from the three specified directions and the string value for the direction
func ComputeNWMaxMove(matrix [][]int, rowSeqs, colSeqs Sequences, rowInd, colInd int, subMatrix Matrix) (int, string){
    /*
    for a given cell i,j what is the move that will get me the maximum score given cells to my left (i-1,j)
    up (i,j-1) and diagonal (i-1,j-1)
    */
    var direction string
    upScore := matrix[rowInd-1][colInd] + SumOfPairs(rowSeqs, colSeqs, 0, colInd, subMatrix) //0 is always gap
    diagScore := matrix[rowInd-1][colInd-1] + SumOfPairs(rowSeqs, colSeqs, rowInd, colInd, subMatrix)
    leftScore := matrix[rowInd][colInd-1] + SumOfPairs(rowSeqs, colSeqs, rowInd, 0, subMatrix) //0 is always gap
    // diagScore := matrix[i-1][j-1] + scoreMap[c1][c2]
    // leftScore := matrix[i][j-1] + scoreMap["gap"]["gap"]
    // upScore := matrix[i-1][j] + scoreMap["gap"]["gap"]
    //in our algorithm, penalties are negative so get direction in max
    // max := ComputeNWMax([3]int{diagScore, upScore, leftScore,})
    max := ComputeNWMax([3]int{upScore, diagScore, leftScore,})
    // var move Move
    // move.value = max
    // checking which move gives me the max score
    if max == diagScore{
        direction = "DIAG"
    } 
    if max == leftScore{
        direction = "LEFT"
    }
    if max == upScore{
        direction =  "UP"
    }
    return max, direction
}

// func (distMatrix *Matrix)UpdateDistMatrix(rowLabel, colLabel string, score int) {
//     innerMap := make(map[string]int)
//     fmt.Println("score passed into update function is", score)
//     innerMap[colLabel] = score
//     distMatrix[rowLabel] = innerMap
// }

//everything after this is the old set of functions

//This function is used to initialize the first row and first column
//with gap penalties
//Input: two seq1 and seq2 string
//Output: a 2D matrix
func InitializeNeedlemanMatrix2D(seq1 string, seq2 string,gapPenalty float64 )[][]float64{

    n1:=len(seq1) + 1
    n2:=len(seq2) + 1
    matrix := make([][]float64,n1)
    for i:=0;i<n1;i++{
        matrix[i] = make([]float64,n2)

    }
    for i:=0;i<n1;i++{
        matrix[i][0] = gapPenalty * float64(i)

    }
    for j:=0;j<n2;j++{
        matrix[0][j] = gapPenalty * float64(j)

    }
    return matrix
}

//this function will intialize the traceback matrix
//Input: seq1 and seq2
//Output: a list of traceback matrix 
func InitializeNeedlemanTracebackMatrix2D(seq1 string, seq2 string) [][]string {

    n1:=len(seq1) + 1
    n2:=len(seq2) + 1
    matrix := make([][]string,n1)
    for i:=0;i<n1;i++{
        matrix[i] = make([]string,n2)

    }
    // this is trivial,once we reach the boundary row, we wont need to check direction
    for i:=0;i<n1;i++{
        matrix[i][0] = "UP"

    }
    // this is trivial,once we reach the boundary column, we wont need to check direction
    for j:=0;j<n2;j++{
        matrix[0][j] = "LEFT"

    }
    return matrix
}
//this is a helper function which essentially assigns in a dictionary the penalties for match,mismatch and a gap
//Input: gap,mismatch penalities and match reward
//Output: a dictionary or a map that specifies what the reward/penalty is for a given token
//e.g map['A']['A'] = matchReward,map['A']['B'] = mismatchPenalty, map['gap']['gap'] = gapPenalty,
func ComputeNeedlemanScoreMap(gapPenalty,matchReward,misMatchPenalty float64)map[string](map[string]float64){
    score := make(map[string](map[string]float64))
    /*
    Assume matchReward = 1, misMatch penalty = -1
    A T G C
    A 1 -1 -1 -1
    T -1 1 -1 -1
    G -1 -1 1 -1
    C -1 -1 -1 1
    
    
    */
    // alphabets := [8]string{"A","T","G","C","a","t","g","c"}
    alphabets := []string{"A","C","D","E","F","G","H","I","K","L","M","N","P","Q","R","S","T","V","W","Y"}
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
//this is a helper function that essentially helps us find the maximum value in ta 3 element list
//Input: a list of three floats
//Output: a float which is the maximum in the input list
func ComputeNeedlemanMaxValue(list [3]float64) float64{

    //help get max value
    max := list[0]
    for i:=0;i<len(list);i++{
        if list[i] > max{
            max = list[i]
        }
    }
    return max
}
//this is a helper function to compute the maximum value for three possible cells: UP,DIAGONAL,LEFT
//Input: a 2D matrix of scores, the current indices i,j to look at and a dictionary specifying the scoring mechanism
//Output: Outputs a Move struct which basically holds the value for the maximum score from the three specified directions as
//well as the string value for the direction
func ComputeNeedlemanMaxMove(matrix [][]float64,i,j int,c1,c2 string,scoreMap map[string](map[string]float64)) Move{
    /*
    for a given cell i,j what is the move that will get me the maximum score given cells to my left (i-1,j)
    up (i,j-1) and diagonal (i-1,j-1)
    */
    diagScore := matrix[i-1][j-1] + scoreMap[c1][c2]
    leftScore := matrix[i][j-1] + scoreMap["gap"]["gap"]
    upScore := matrix[i-1][j] + scoreMap["gap"]["gap"]
    max := ComputeNeedlemanMaxValue([3]float64{diagScore,leftScore,upScore})
    var move Move
    move.value = max
    // checking which move gives me the max score
    if max == diagScore{
        move.direction = "DIAG"

    } 
    if max == leftScore{
        move.direction = "LEFT"

    }
    if max == upScore{
        move.direction =  "UP"

    }
    return move
}
//this function is like a forward pass of the needleman algorthm whereby we compute the score for a given subsequence ending at i,j
//Input: the input sequences seq1 and seq2 and the scoring dictionary
//Output: two matrices, the first matrix represents a matrix of alignment scores whereas the second matrix consists of the direction of travel in the forward pass
// of needleman algorithm e.g which direcion we took from a given cell i,j
func ComputeNeedlemanScores(seq1 string,seq2 string,scoreMap map[string](map[string]float64))([][]float64,[][]string){
    matrix := InitializeNeedlemanMatrix2D(seq1, seq2,scoreMap["gap"]["gap"])
    traceBackMatrix := InitializeNeedlemanTracebackMatrix2D(seq1,seq2)
    s1 := len(matrix)
    s2 := len(matrix[0])
    for i:=1;i<s1; i++ {
        for j:=1;j<s2; j++{
            move := ComputeNeedlemanMaxMove(matrix,i,j,string(seq1[i-1]),string(seq2[j-1]),scoreMap)
            matrix[i][j] = move.value
            traceBackMatrix[i][j] = move.direction
        }

    }
    //fmt.Println("matrix:",matrix)
    //fmt.Println("traceBackMatrix:",traceBackMatrix)
    return matrix,traceBackMatrix

}

//this function will compute the aligned sequences using the traceback matrix
//Input: sequences seq1 and seq1 and the traceback matrix which has the the direction taken during the Needleman algorithm
//Output: a pair of aligned sequences
func ComputeNeedlemanAlignments(seq1 string,seq2 string,traceBackMatrix [][]string)(string,string){
    n1:=len(seq1) + 1
    n2:=len(seq2) + 1
    newSeq1 := ""
    newSeq2 := ""
    i := n1 - 1
    j := n2 - 1
    direction := traceBackMatrix[i][j] //start at max row and col of matrix
    
    for i > 0 ||  j > 0 {
        //fmt.Println("Direction:",direction,"i:","j", i,j,"seq1:",newSeq1,"seq2:",newSeq2)
        if direction == "DIAG"{
            if i > 0 {
                newSeq1 =  string(seq1[i-1]) + newSeq1
                // fmt.Println(string(seq1[i-1]))
            }
            if j > 0 {
                newSeq2 =  string(seq2[j-1]) + newSeq2
            }         
            i -=1
            j-=1
        } else if direction == "LEFT"{
            newSeq1 =  "-" + newSeq1
            if j > 0 {
                newSeq2 =  string(seq2[j-1]) + newSeq2
            }
            j-=1
        } else if direction == "UP"{
            if i > 0 {
                newSeq1 =  string(seq1[i-1]) + newSeq1
            }
            newSeq2 =  "-"+ newSeq2
            i -=1 
        }
        direction = traceBackMatrix[i][j]    
    }
    return newSeq1,newSeq2
}
//just a helper function to visualize matrix nicely
func PrettyPrintMatrix(m [][]string){
    a := len(m)
    for i:=0;i<a;i++{
        fmt.Println(m[i])

    }
}

//this is the highest-level needleman algorithm which will be used to align two sequences
//Input: the two sequences to align seq1 and seq2
//Output: the aligned sequences for seq1,seq2 as well as the aligment score
func Needleman(seq1 string,seq2 string) (string,string,float64){
    gapPenalty := -1.0
    matchReward:= 1.0
    misMatchPenalty := -1.0
    s1,s2 := len(seq1)+1,len(seq2)+1
    //this is the scoring matrix e.g what is the score for a match,mismatch,gap etc
    scoreMap := ComputeNeedlemanScoreMap(gapPenalty,matchReward,misMatchPenalty) // this initializes the scoring mechanism
    //should change this part to take in a "scoring lookup table" to get match/mismatch scores from a .csv file in main.go or io.go
    matrix,traceBackMatrix := ComputeNeedlemanScores(seq1,seq2,scoreMap)
    //fmt.Println("matrix:",matrix)
    //fmt.Println("scoreMap:",scoreMap)
    //PrettyPrintMatrix(traceBackMatrix)
    alignedSeq1,alignedSeq2 := ComputeNeedlemanAlignments(seq1,seq2,traceBackMatrix)
    
    alignmentScore := matrix[s1-1][s2-1] //the score taken from bottom corner is the alignment score- use this to build difference matrix

    return alignedSeq1,alignedSeq2,alignmentScore

}

//this function will be useful when building the guide tree needed for alignments
//Input: list of sequences to align
//Output: a map of maps 
//The map will represent the dissimilarity between every pair of sequences
//e.g map[seq1][seq2] = 1 - needleman(seq1,seq2)
func ConstructDistanceMap(seqs []string)map[string](map[string]float64){
    score := make(map[string](map[string]float64))
    /*
    Assume matchReward = 1, misMatch penalty = -1
      A T G C
    A 1 -1 -1 -1
    T -1 1 -1 -1
    G -1 -1 1 -1
    C -1 -1 -1 1
    
    
    */
    for i:=range(seqs){
        score[seqs[i]] = make(map[string]float64)
        for j:=range(seqs){
            if seqs[i] == seqs[j]{
                score[seqs[i]][seqs[j]] = 0
            }else{
                _,_,alignmentScore := Needleman(seqs[i],seqs[j])
                score[seqs[i]][seqs[j]] = -alignmentScore

            }
        }
    }
    return score

}


// func main(){
//     seq1 :=  "MSLTAKDKSVVKAFWGKISGKADVVGAEALGRVLTAYPQTKTYFSHWADLSPGSGPVKKHGGIIMGAIGKAVGLMDDLVGGMSALSDLHAFNLRVDPGNFKILSHNILVTLAIHFPSDFTPEVHIAVDKFLAVVSAALADKYR"[:20] //ykiss_Rainbow_trou
//     seq2 := "MHLTADDKKHIKAIWPSVAAHGDKYGGEALHRMFMCAPKTKTYFPDFDFSEHSKHILAHGKKVSDALNEACNHLDNIAGCLSKLSDLHAYDLRVDPGNFPLLAHQILVVVAIHFPKQFDPATHKALDKFLVSVSNVLTSKYR"[:20] //Xenopus_tropicalis_Western_clawed_frog
//     alignedSeq1,alignedSeq2,alignmentScore:= Needleman(seq1,seq2)
//     fmt.Println("alignedSeq1:",alignedSeq1)
//     fmt.Println("alignedSeq2:",alignedSeq2)
//     fmt.Println("Score:",alignmentScore)
//     distanceMap := ConstructDistanceMap([]string{seq1,seq2})
//     fmt.Println("distanceMap:",distanceMap)

   

// }
