package main

// import "fmt"
//import "os"
//import "strings"
//import "bufio"

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
    // fmt.Println("matrix:",matrix)
    // fmt.Println("scoreMap:",scoreMap)
    // PrettyPrintMatrix(traceBackMatrix)
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
    // fmt.Println("length of rowSeq is", len(rowSeqs[0].sequence))
    // fmt.Println("length of colSeq is", len(colSeqs[0].sequence))
    //make a slice of sequences (strings) for rowSeq, colSeq, equal to how many there are, set them each to ""
    newRowSeqs := make([]Sequence, len(rowSeqs))
    newColSeqs := make([]Sequence, len(colSeqs)) //make new Sequences objects
    // newRowSeq := make([]string, len(rowSeqs))
    // newColSeq := make([]string, len(colSeqs)) //will have default value of ""

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
    
    for i > 0 || j > 0 { //i refers to the ith character in rowSeqs; j to the jth character in rowSeqs
        newRowChars := make([]string, len(rowSeqs)) //new slice of characters, one for every sequence in rowSeqs
        newColChars := make([]string, len(colSeqs)) //new slice of characters, one for every sequence in rowSeqs
        newRowChars, newColChars, i, j = GetNextChars(rowSeqs, colSeqs, direction, i, j)
        if i >= 0 {
            for m := 0; m < len(newRowSeqs); m++ { //m counts the number of Sequence objects in Sequences
                newRowSeqs[m].sequence = newRowChars[m] + newRowSeqs[m].sequence //add this character
                // fmt.Println(string(rowSeqs[m].sequence[i]))
            }
            // newSeq1 =  string(seq1[i-1]) + newSeq1
        }
        if j >= 0 {
            for n := 0; n < len(newColSeqs); n++ { //n counts the number of Sequence objects in Sequences
                newColSeqs[n].sequence = newColChars[n] + newColSeqs[n].sequence //add this character to every Sequence in rowSeqs
                // fmt.Println(string(colSeqs[n].sequence[j]))
            }
            // newSeq2 =  string(seq2[j-1]) + newSeq2
        }
        // i, j = 
        direction = traceBackMatrix[i][j]
        direction = traceBackMatrix[i][j] //extra line?   
    }
    //combine rowSeqs and colSeqs into one sequences object
    mergedSeqs := append(newRowSeqs, newColSeqs...)
    return mergedSeqs
}

//GetNextChars determines the next character/nucleotide/amino acid that 
func GetNextChars(rowSeqs, colSeqs Sequences, direction string, i, j int) ([]string, []string, int, int) {
    newRowChars := make([]string, len(rowSeqs)) //new slice of characters, one for every sequence in rowSeqs
    newColChars := make([]string, len(colSeqs)) //new slice of characters, one for every sequence in rowSeqs
    if direction == "DIAG"{ //add the character at that position to everybody, decrement i and j by 1
        for m := 0; m < len(rowSeqs); m++ { //m counts the number of Sequence objects in Sequences
            newRowChars[m] = string(rowSeqs[m].sequence[i])
        }
        for n := 0; n < len(colSeqs); n++ { //n counts the number of Sequence objects in Sequences
            newColChars[n] = string(colSeqs[n].sequence[j])
        }
        i--
        j-- //decrement both i and j
    } else if direction == "LEFT"{ //add character at that position to column sequence, "-" to row; decrement the column count by 1
        for m := 0; m < len(rowSeqs); m++ { //m counts the number of Sequence objects in Sequences
            newRowChars[m] = "-"
        }
        for n := 0; n < len(colSeqs); n++ { //n counts the number of Sequence objects in Sequences
            newColChars[n] = string(colSeqs[n].sequence[j])
        }
        j--
    } else if direction == "UP"{ //add character at that position to row sequence, "-" to column; decrement the row count by 1
        for m := 0; m < len(rowSeqs); m++ { //m counts the number of Sequence objects in Sequences
            newRowChars[m] = string(rowSeqs[m].sequence[i])
        }
        for n := 0; n < len(colSeqs); n++ { //n counts the number of Sequence objects in Sequences
            newColChars[n] = "-"
        }
        i-- 
    }
    // fmt.Println("newColChars in function", newColChars)
    // fmt.Println("newColChars in function", len(newColChars))
    return newRowChars, newColChars, i, j
}


//this function is like a forward pass of the needleman algorithm whereby we compute the score for a given subsequence ending at i,j
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

//this is a helper function that essentially helps us find the maximum value in a 3 element list
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
    // fmt.Println("current rowSeqs letter is", string(rowSeqs[0].sequence[rowInd]), rowInd, "current colSeqs letter is", string(colSeqs[0].sequence[colInd]), colInd)
    upScore := matrix[rowInd-1][colInd] + SumOfPairs(rowSeqs, colSeqs, rowInd, 0, subMatrix) //0 is always gap
    diagScore := matrix[rowInd-1][colInd-1] + SumOfPairs(rowSeqs, colSeqs, rowInd, colInd, subMatrix)
    leftScore := matrix[rowInd][colInd-1] + SumOfPairs(rowSeqs, colSeqs, 0, colInd, subMatrix) //0 is always gap
    // diagScore := matrix[i-1][j-1] + scoreMap[c1][c2]
    // leftScore := matrix[i][j-1] + scoreMap["gap"]["gap"]
    // upScore := matrix[i-1][j] + scoreMap["gap"]["gap"]
    //in our algorithm, penalties are negative so get direction in max
    // max := ComputeNWMax([3]int{diagScore, upScore, leftScore,})  
    neighborVals := [3]int{diagScore, leftScore, upScore,}
    max := ComputeNWMax(neighborVals)
    // checking which move gives me the max score

    if max == diagScore{
        direction = "DIAG"
    } else if max == leftScore{
        direction = "LEFT"
    } else if max == upScore{
        direction =  "UP"
    }
   
    return max, direction
}



//everything after this is the old set of functions

