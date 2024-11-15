//Blake Kiefer (bkiefer)

package main

//Input: CreateDistanceMatrix takes as input an array of sequences
//Output: It returns a 2D array of integers distanceMatrix representing the number of differences between each sequence to every other sequence.
func CreateDistanceMatrix(sequences []string) [][]float64 {
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
			dissimilarity := 1 - SequenceSimilarity(sequences[r], sequences[c])

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