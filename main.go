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
    // EXAMPLE COMMAND: ./MultipleSeqAlignment Protein alignment BLOSUM62.csv -8 peptidehormones
    // Main steps of our main function are:
    // 1: Read command line args, we expect 5 (other than the os.Args[0]) as shown in example command above
    // 2: Read scoring matrix input file
    // 3: Read sequences file supplied as .fasta or .txt in FASTA format
    // 4: Create Distance Matrix (DM) using distMethod
    // 5: Create guide tree off DM and build MSA using Neighbor join (NJ) tree method. Export guide tree in newick format
    // 6: Build a phylogentic tree off aligned sequences from NJ and export tree in newick format
    // 7: You can expect output to be stored in output/
    if len(os.Args) != 6 {
        panic("ERROR:error: incorrect number of inputs buddy")
    }
    //1: Read command line args, we expect 5 (other than the os.Args[0]) as shown in example command above
    seqType := os.Args[1]
    distMethod := os.Args[2]
    filename := os.Args[3]
    //Gap penalty: for proteins, try values from -7 to -12; DNA/RNA -2 or -3
    //if you're doing simple scoring- 1 for match, -1 for mismatch- do -1 or -2 
    gapScore, err := strconv.Atoi(os.Args[4])
    //folderName := os.Args[5]
    filePath := os.Args[5]
    fmt.Println("INFO:Read arguments:",os.Args[1:])
    
    
    //2: Read scoring matrix input file
    fmt.Println("INFO:Reading scoring matrix")
    subMatrix, err := ReadScoringMatrix(seqType, filename, gapScore) //subMatrix contains the reference scoring matrix
    if err != nil {
        fmt.Println("ERROR:Error:", err)
        return
    }
    //3: Read sequences file supplied as .fasta or .txt in FASTA format
    fmt.Println("INFO:Reading sequences")
    inputSeqs, err := readFASTAFile(filePath)
    if err != nil {
        fmt.Println("ERROR:Error:", err)
        return
    }

    fmt.Println("INFO:Read final sequences, number of sequences is", len(inputSeqs))
    

    //4: Create Distance Matrix (DM) using distMethod
    fmt.Println("INFO:Creating distance matrix")
    distMatrix := CreateDistanceMatrix(distMethod, inputSeqs, subMatrix)   
    fmt.Println("INFO:Building guide tree")
    guideTree := NJ(distMatrix, inputSeqs, subMatrix)
    
    //5: Create guide tree off DM and build MSA using Neighbor join (NJ) tree method. Export guide tree in newick format
    fmt.Println("INFO:Converting Guide Tree to newick format")
    root := guideTree[len(guideTree)-1]
	newick := ConvertTreeToNewick(guideTree, root)
	fileName := "output/guidetree.newick"
	err = ExportNewickToFile(newick, fileName)
	if err != nil{
		fmt.Println("ERROR:Error exporting Newick:", err)
	} else {
		fmt.Println("INFO:Newick exported to", fileName)
	}

    fmt.Println("INFO:Writing our MSA to output/msa.fasta")
	WriteFASTAOutput("output/msa.fasta", root.sequence) //output aligned sequences


    //6: Build a phylogentic tree off aligned sequences from NJ and export tree in newick format
    fmt.Println("INFO:Building a phylogenetic tree based on aligned Sequences using distMethod=naivehomology")
    fileName = "output/phylotree.newick"
    alignedDistMatrix := CreateDistanceMatrix("naivehomology", root.sequence, subMatrix)
    finalTree := NJ(alignedDistMatrix, root.sequence, subMatrix)
    finalroot := finalTree[len(finalTree)-1]
    fmt.Println("INFO:Converting Phylo Tree to newick format")
	newick = ConvertTreeToNewick(finalTree, finalroot)
	err = ExportNewickToFile(newick, fileName)
	if err != nil{
		fmt.Println("ERROR:Error exporting Newick:", err)
	} else {
		fmt.Println("INFO:Newick exported to", fileName)
	}
    
}

