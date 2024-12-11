//Blake Kiefer (bkiefer), William Lu, Talha Khan

package main

import (
	"fmt"
	// "math"
)

//"naivehomology"- naive
//"alignedhomology" - aligned homology
//"alignmentscore" - needleman

func CreateDistanceMatrix(distMethod string, inputSeqs Sequences, subMatrix Matrix) [][]float64 {
	distanceMatrix := make([][]int, len(inputSeqs)) //hold int scores- results of pairwise alignments
    // seqMatrix := make([][]Sequences, len(inputSeqs)) //get rid of this or comment out
    diffMatrix := make([][]float64, len(inputSeqs))
    if distMethod == "alignmentscore" {
        distanceMatrix = NWDistMatrix(inputSeqs, subMatrix) //maybe output a [][]float64 from here and not cast
        diffMatrix = MatrixToFloats(distanceMatrix)
    } else if distMethod == "alignedhomology" { //"identity"
		diffMatrix = CreateAlignedHomologyMatrix(inputSeqs, subMatrix)
    } else if distMethod == "naivehomology" {
		diffMatrix = CreateNaiveHomologyMatrix(inputSeqs) //fix data type
	} else { //if none of the above
		panic("invalid distance matrix creation type")
	}
	fmt.Println("diffMatrix is")
	for n := 0; n < len(diffMatrix); n++ {
		fmt.Println(diffMatrix[n])
	}
	normDistMatrix := MinMaxNormalize(diffMatrix)
	fmt.Println("normeddistMatrix is")
	for n := 0; n < len(normDistMatrix); n++ {
		fmt.Println(normDistMatrix[n])
	}
	return normDistMatrix
}

//the function to do Needleman repeatedly- every sequence against every other sequence, to create a distance matrix
//input: inputSeqs Sequences, subMatrix Matrix
//output: a distanceMatrix Matrix and seqMatrix SequenceMatrix
func NWDistMatrix(inputSeqs Sequences, subMatrix Matrix) [][]int{
    //code for multiple NWs on multiple sequences - this should be modularized
    // var distanceMatrix Matrix 
    // distanceMatrix := make(Matrix) 
    fmt.Println("starting to needleman everybody", inputSeqs)
    distanceMatrix := make([][]int, len(inputSeqs)) //hold int scores- results of pairwise alignments
    for i := 0; i < len(distanceMatrix); i++ {
        distanceMatrix[i] = make([]int, len(inputSeqs))
    }
    // seqMatrix := make([][]Sequences, len(inputSeqs)) //hold Sequences objects of length 2- no need to redo alignments when filling guide tree nodes, just fill from here
    // for i := 0; i < len(seqMatrix); i++ {
    //     seqMatrix[i] = make([]Sequences, len(inputSeqs))
    // }
    // var seqMatrix SequenceMatrix 
    for i := 0; i < len(inputSeqs); i++ {
        // distanceMatrix[i] = make([]int, len(inputSeqs))
        // seqMatrix[i] = make([]Sequences, len(inputSeqs))
        for j := i; j < len(inputSeqs); j++ {
            if j == i { //this is janky
                distanceMatrix[j][i] = 0
            } else {
                // fmt.Println("printing", inputSeqs[i:i+1], inputSeqs[j:j+1])
                // if i >= j {
                _, score := NeedlemanWunsch(inputSeqs[i:i+1], inputSeqs[j:j+1], subMatrix) //returns a Sequences object- the aligned two sequences, int the alignment score
                // PrintSequencesList(alignedTwo)
                // fmt.Println("alignment score is", score)
                distanceMatrix[i][j] = score
                distanceMatrix[j][i] = score
                // seqMatrix[i][j] = alignedTwo
            // distanceMatrix.UpdateDistMatrix(inputSeqs[i].info, inputSeqs[j].info, score)
            // innerMap := make(map[string]int)
            // fmt.Println("score passed into update function is", score)
            // innerMap[inputSeqs[j].info] = score
            // distanceMatrix[inputSeqs[i].info] = innerMap //why is the value 0?
            // fmt.Println("distance matrix looks like: ")
                // fmt.Println("distance matrix is:")
                // for n := 0; n < len(distanceMatrix); n++ {
                //     fmt.Println(distanceMatrix[n])
                // }
            // PrintSubMatrix(distanceMatrix)
            // distanceMatrix[inputSeqs[i].info][inputSeqs[j].info] = score
            // seqMatrix[inputSeqs[i].info][inputSeqs[j].info] = alignedTwo
            // }
            }  
        }
    }
    return distanceMatrix
}

//Input: CreateDistanceMatrix takes as input an array of sequences
//Output: It returns a 2D array of integers distanceMatrix representing the number of differences between each sequence to every other sequence.
func CreateAlignedHomologyMatrix(inputSeqs Sequences, subMatrix Matrix) [][]float64 { //func name was changed from CreateDistanceMatrix
	//Create an distance matrix with a number of rows equal to the number of sequences and a number of columns equal to the number of sequences.
	n := len(inputSeqs)
	distanceMatrix := InitializeDistanceMatrix(n)

	//Each row represents a sequence.  Each column represents a sequence.  Range through each row and each column within that row to calculate the distance metric between the sequence of that row with ever other sequence.
	//The diagonal values remain 0 because the distance between a sequence and itself is 0.0.
	//The matrix is symmetrical because when comparing two different sequences to each other, the difference between sequence1 and sequence2 is the same as the difference between sequence2 and sequence1.
	//Start column index at one more index above the index of the row because all the diagonals are going to be 0.0 anyway.
	//Only have to range rows to the second to last row since after filling that row the bottom row will already be complete due to the matrix symmetry.
	for r := 0; r < n - 1; r ++ {
		for c := r + 1; c < n; c ++ {
			alignedTwo, _ := NeedlemanWunsch(inputSeqs[r:r+1], inputSeqs[c:c+1], subMatrix)
			dissimilarity := 1 - SequenceSimilarity(alignedTwo[0].sequence, alignedTwo[1].sequence)

			//This value should be entered into the two corresponding indices on either side of the diagonal since the matrix is symmetrical.
			distanceMatrix[r][c] = dissimilarity
			distanceMatrix[c][r] = dissimilarity
		}
	}
	return distanceMatrix
}

//Input: CreateDistanceMatrix takes as input an array of sequences
//Output: It returns a 2D array of integers distanceMatrix representing the number of differences between each sequence to every other sequence.
func CreateNaiveHomologyMatrix(sequences Sequences) [][]float64 { //func name was changed from CreateDistanceMatrix
	//Create an distance matrix with a number of rows equal to the number of sequences and a number of columns equal to the number of sequences.
	n := len(sequences)
	distanceMatrix := InitializeDistanceMatrix(n)

	//Each row represents a sequence.  Each column represents a sequence.  Range through each row and each column within that row to calculate the distance metric between the sequence of that row with ever other sequence.
	//The diagonal values remain 0 because the distance between a sequence and itself is 0.0.
	//The matrix is symmetrical because when comparing two different sequences to each other, the difference between sequence1 and sequence2 is the same as the difference between sequence2 and sequence1.
	//Start column index at one more index above the index of the row because all the diagonals are going to be 0.0 anyway.
	//Only have to range rows to the second to last row since after filling that row the bottom row will already be complete due to the matrix symmetry.
	for r := 0; r < n - 1; r ++ {
		for c := r + 1; c < n; c ++ {
			dissimilarity := 1 - SequenceSimilarity(sequences[r].sequence, sequences[c].sequence)

			//This value should be entered into the two corresponding indices on either side of the diagonal since the matrix is symmetrical.
			distanceMatrix[r][c] = dissimilarity
			distanceMatrix[c][r] = dissimilarity
		}
	}

	return distanceMatrix
}

//Input: InititalizeDistanceMatrix takes in an integer n representing the number of sequences.
//Ouput: It returns a 2D array with n rows and n columns.  Each float64 within the array is by default 0.
func InitializeDistanceMatrix(n int) [][]float64 {
	distanceMatrix := make([][]float64, n)

	for r := 0; r < n; r ++ {
		distanceMatrix[r] = make([]float64, n)
	}

	return distanceMatrix
}

//Input: SequenceSimilarity() takes as input two strings that are two sequences.
//Output: It returns a sequence similarity score representing the similarity between the two sequences.
func SequenceSimilarity(sequence1, sequence2 string) float64 {
	//A rough comparision is sufficient to create the guide tree, and the sequences are not aligned yet.  It is sufficient to range only throught the length of the shorter sequence to determine the similarity of the sequences.
	shortSeqLen := IdentifyShorterSequence(sequence1, sequence2)

	commonNucleotides := 0
	//Range through the length of the shorter sequence.  Start at index 0 of both sequences since strings can be indexed. If the nucleotide in both sequences is the same, add one to the number of commonNucleotides encountered thus far.
	for i := 0; i < shortSeqLen; i ++ {
		if sequence1[i] == sequence2[i] {
			commonNucleotides ++
		}
	}

	//The final similarity score between two sequences is the number of commonNucleotides counted between them divided by the number of nucleotides that have been ranged over (the length of the shorter sequence)
	return float64(commonNucleotides)/float64(shortSeqLen)
}

//Input: IdentifyShorterSequence() takes as input two strings representing two sequences.
//Output: It returns the length of the shorter sequence between the two.
func IdentifyShorterSequence(sequence1, sequence2 string) int {
	//If sequence one is shorter than sequence two, return the length of sequence one.  This is what we will range through when counting a rough similarity score between the two sequences.
	if len(sequence1) < len(sequence2) {
		return len(sequence1)
	
	//If sequence two is shorter than sequence one, return the length of sequence 2.  We will also return the length of sequence two if the lengths of the two sequences are identical.  At that point it doesn't matter which length of which sequence we return because they are identical.  Here we have arbitrarily picked to return the length of sequence two in such a case.
	} else {
		return len(sequence2)
	}
}

//MinMaxNormalize normalizes values of a distance matrix such that the min val is 0 and max val is 1, and all values distributed between
//input: a distanceMatrix 2D slice of floats
//returns: another 2D slice of floats, but min-max normalized
func MinMaxNormalize(distanceMatrix [][]float64) [][]float64 {
	normMatrix := make([][]float64, len(distanceMatrix))
	min := FindMinVal(distanceMatrix)
	max := FindMaxVal(distanceMatrix)
	for r := range distanceMatrix {
		normMatrix[r] = make([]float64, len(distanceMatrix[r]))
		for c := range distanceMatrix[r] {
			if r == c {
				normMatrix[r][c] = 0 //the distance between two identical sequences should be 0
			} else {
				normMatrix[r][c] = (distanceMatrix[r][c] - min)/(max-min)
			}
		}
	}
	return normMatrix
}

//FindMinVal finds the minimum val in a 2D slice of floats
//input: a distanceMatrix 2D slice of floats
//returns the minimum value in that 2D slice
func FindMinVal(distanceMatrix [][]float64) float64 {
	min := distanceMatrix[0][0]
	for r := range distanceMatrix {
		for c := range distanceMatrix[r] {
			if distanceMatrix[r][c] < min {
				min = distanceMatrix[r][c]//the distance between two identical sequences should be 0
			}
		}
	}
	return min
} 

//FindMaxVal finds the maximum val in a 2D slice of floats
//input: a distanceMatrix 2D slice of floats
//returns the maximum value in that 2D slice
func FindMaxVal(distanceMatrix [][]float64) float64 {
	max := distanceMatrix[0][0]
	for r := range distanceMatrix {
		for c := range distanceMatrix[r] {
			if distanceMatrix[r][c] > max {
				max = distanceMatrix[r][c]//the distance between two identical sequences should be 0
			}
		}
	}
	return max
} 

//MatrixToFloats converts all values in a 2D slice of ints into float64
//input: a distanceMatrix 2D slice of ints
//returns that distanceMatrix, all values cast to floats
func MatrixToFloats(distanceMatrix [][]int) [][]float64{
    //we can clean this later- but take distanceMatrix and convert to [][]float64
    diffMatrix := make([][]float64, len(distanceMatrix))
    for i := 0; i < len(distanceMatrix); i++ {
        diffMatrix[i] = make([]float64, len(distanceMatrix))
        for j := 0; j < len(distanceMatrix); j++ {
            diffMatrix[i][j] = float64(distanceMatrix[i][j])
        }
    }
    return diffMatrix
}