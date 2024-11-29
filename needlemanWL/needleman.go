package main
import "fmt"
//import "os"
//import "strings"
//import "bufio"
type  Move struct {
    direction string
    value float64
}
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
func InitializeNeedlemanTracebackMatrix2D(seq1 string, seq2 string )[][]string{

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
    reverseSeq1 := ""
    reverseSeq2 := ""
    i := n1 - 1
    j := n2 - 1
    direction := traceBackMatrix[i][j] //start at max row and col of matrix
    
    for i > 0 ||  j > 0 {
        //fmt.Println("Direction:",direction,"i:","j", i,j,"seq1:",reverseSeq1,"seq2:",reverseSeq2)
        if direction == "DIAG"{
            if i > 0 {
                reverseSeq1 =  string(seq1[i-1]) + reverseSeq1
            }
            if j > 0 {
                reverseSeq2 =  string(seq2[j-1]) + reverseSeq2
            }         
            i -=1
            j-=1
        } else if direction == "LEFT"{
            reverseSeq1 =  "-" + reverseSeq1
            if j > 0 {
                reverseSeq2 =  string(seq2[j-1]) + reverseSeq2
            }
            j-=1
        } else if direction == "UP"{
            if i > 0 {
                reverseSeq1 =  string(seq1[i-1]) + reverseSeq1
            }
            reverseSeq2 =  "-"+ reverseSeq2
            i -=1 
        }
        direction = traceBackMatrix[i][j]    
    }
    return reverseSeq1,reverseSeq2
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
