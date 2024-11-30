package main

func reallyMain() int{
	//command line- ./needleman DNA BLASTDNA
	//or proteins- ./needleman Protein BLOSUM62 
	//read in directory containing FASTAs of MSA sequences into sequence, then sequences type
	//specify if working with DNA or protein? to find directory of scoring matrices to use
	//read in a scoring matrix to use
	//MSA(sequences, scoringmatrix)
	//create a new distMatrix Matrix object for distance matrix- need type floats?
	//for every i sequence in sequences
	//	for j is i to len(sequences)	
	//		sequencePair, score := Needleman(i, j, matrix) 
	//		distMatrix[i][j] = score
	// can sequence pair be saved somehow for guide tree construction?

	return 3
}