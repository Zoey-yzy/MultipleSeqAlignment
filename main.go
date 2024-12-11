package main

import (
	"fmt"
	// "encoding/csv"
    "os"
    // "path/filepath"
    "strconv"
    // "sort"
    // "strings"
    // "bufio"
)

func main() { //commandline taken from RShiny
    // ./MultipleSeqAlignment Protein alignment BLOSUM62.csv -8 peptidehormones
    
    if len(os.Args) != 6 {
        panic("error: incorrect number of inputs buddy")
    }

    seqType := os.Args[1]

    distMethod := os.Args[2]
    fmt.Println(distMethod)

    filename := os.Args[3]

    gapScore, err := strconv.Atoi(os.Args[4])

    folderName := os.Args[5]
    
    // gapScore := -2 //hardcoded/read from commandline/RShiny- add to subMatrix
    //for proteins, try values from -7 to -12; DNA/RNA -2 or -3
    //if you're doing simple scoring- 1 for match, -1 for mismatch- do -1 or -2 
    
    //later change this to be doable through RShiny- folder browsing or dropdown menu
    //reading scoring matrix
    fmt.Println("reading scoring matrix")
    // seqType := "Protein"
    fmt.Println("working with", seqType)
	// filename := "Simple.csv"
    subMatrix, err := ReadScoringMatrix(seqType, filename, gapScore) //subMatrix contains the reference scoring matrix
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(subMatrix)
    PrintSubMatrix(subMatrix)
	// fmt.Println("A, A", matrix["A"]["A"], "A, T", matrix["A"]["T"])

    //reading sequences
    fmt.Println("reading sequences")
    // folderName := "peptidehormones"
    inputSeqs, err := ReadFASTAInput(seqType, folderName) //inputSeqs is a Sequences object
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("printing final sequences, number of sequences is", len(inputSeqs))
    for _, seq := range inputSeqs {
        //fmt.Printf("Info: %s, Sequence: %s\n", seq.info, seq.sequence)
        fmt.Println("the sequence is:", seq.info, seq.sequence)
    }

    //START WITH THESE TWO GUYS TO DO TEST CASES

    // //to test the alignments- use foldername alignmenttest2 or alignmenttest, either in DNA or Protein
    // seqStack, score := NeedlemanWunsch(inputSeqs[0:3], inputSeqs[3:6], subMatrix)
    // PrintSequencesList(seqStack)
    // fmt.Println("score is", score)

    // //to test the alignments of alignment to sequence- use foldername alignmenttestDurand, either in DNA or Protein
    // //use a scoring matrix of mismatch = -3, gap = -2, match = 0- filename Durand
    // // compare against the durand sumofpairs slides
    // seqStack, score := NeedlemanWunsch(inputSeqs[0:2], inputSeqs[2:3], subMatrix)
    // PrintSequencesList(seqStack)
    // fmt.Println("score is", score)
  
    // sekwences, alignmentScorr := NeedlemanWunsch(inputSeqs[0:1], inputSeqs[1:2], subMatrix)
    // PrintSequencesList(sekwences)
    // fmt.Println("score was", alignmentScorr)

    distMatrix := CreateDistanceMatrix(distMethod, inputSeqs, subMatrix)
    // fmt.Println("distance matrix looks like: ")
    fmt.Println("distance matrix is:")
    for n := 0; n < len(distMatrix); n++ {
        fmt.Println(distMatrix[n])
    }
    fmt.Println("building guide tree")
    guideTree := NJ(distMatrix, inputSeqs, subMatrix)
    // fmt.Println(guideTree[0].sequence)
    root := guideTree[len(guideTree)-1]
	
	// WriteNewickToFile(guideTree, "output", "tree.newick")
	newick := ConvertTreeToNewick(guideTree, root)
	fmt.Println("newick:",newick)
	fileName := "output/covidguidetree.newick"
	err = ExportNewickToFile(newick, fileName)
	if err != nil{
		fmt.Println("Error exporting Newick:", err)
	} else {
		fmt.Println("Newick exported to", fileName)
	}
    PrintSequencesList(guideTree[len(guideTree)-1].sequence)
	WriteFASTAOutput("output/msa.fasta", root.sequence) //output aligned sequences

    fileName = "output/covidtree.newick"
    alignedDistMatrix := CreateNaiveHomologyMatrix(root.sequence)
    fmt.Println("building phylogenetic tree")
    finalTree := NJ(alignedDistMatrix, root.sequence, subMatrix)
    // fmt.Println(guideTree[0].sequence)
    finalroot := finalTree[len(finalTree)-1]
	
	// WriteNewickToFile(guideTree, "output", "tree.newick")
	newick = ConvertTreeToNewick(finalTree, finalroot)
	fmt.Println("newick:",newick)
	// fileName := "output/covidfinaltree.newick"
	err = ExportNewickToFile(newick, fileName)
	if err != nil{
		fmt.Println("Error exporting Newick:", err)
	} else {
		fmt.Println("Newick exported to", fileName)
	}
    
}

