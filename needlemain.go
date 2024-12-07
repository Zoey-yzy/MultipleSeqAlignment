package main

import (
	"fmt"
	"encoding/csv"
    "os"
    "path/filepath"
    "strconv"
    "sort"
    "strings"
    "bufio"
)

//for getting mismatch penalties for specific combos of nucleotides/amino acids
type Matrix map[string]map[string]int
//write a function in main.go to read in a table- PAM, BLOSUM, or transition/transversion matrix from a .txt file
//and a function to fill in gap penalties into scoring matrix based on a given value in main?
//reusable to represent a distance matrix? will we have floats for that?

type Sequence struct{
    sequence    string
    info        string //we don't know what source yet, so we don't know how the FASTA header is formatted aside from "<"
    //match     string //same length as sequence- "*" if letter at this position matches everybody else, "^" otherwise, ???
}

//input format for MSA- of unaligned sequences, as well as outputs of aligned sequences
//how do we use this for initial sequence to sequence NW?
type Sequences []Sequence

type SequenceMatrix map[string]map[string]Sequences
//to store pairwise alignments themselves after pairwise NW

func main() { //commandline taken from RShiny
    // ./MultipleSeqAlignment Protein alignment BLOSUM62.csv -8 peptidehormones
    
    if len(os.Args) != 6 {
        panic("error: incorrect number of inputs buddy")
    }

    seqType := os.Args[1]

    algType := os.Args[2]
    fmt.Println(algType)

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
        fmt.Println(seq.info, seq.sequence)
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
       
    distanceMatrix := make([][]int, len(inputSeqs)) //hold int scores- results of pairwise alignments
    seqMatrix := make([][]Sequences, len(inputSeqs))
    diffMatrix := make([][]float64, len(inputSeqs))
    if algType == "alignment" {
        distanceMatrix, seqMatrix = NWEverybody(inputSeqs, subMatrix)
        diffMatrix = MatrixToFloats(distanceMatrix)
    } else { //"identity"
        fmt.Println("somebody who understands, if you enter 'identity' as the os.Args[2], then don't do multiple pairwise needleman alignments and build the distance matrix using sequence identity instead")
    }  
    

    fmt.Println("building tree")
    guideTree := NJ(diffMatrix, seqMatrix, inputSeqs, subMatrix)
    // fmt.Println(guideTree[0].sequence)
    root := guideTree[len(guideTree)-1]
    newick := ConvertTreeToNewick(guideTree, root)
	fileName := "tree.newick"
	if err := ExportNewickToFile(newick, fileName); err != nil {
		fmt.Println("Error exporting Newick:", err)
	} else {
		fmt.Println("Newick exported to", fileName)
	}
    PrintSequencesList(guideTree[len(guideTree)-1].sequence)

    // seq1 :=  "MSLTAKDKSVVKAFWGKISGKADVVGAEALGRVLTAYPQTKTYFSHWADLSPGSGPVKKHGGIIMGAIGKAVGLMDDLVGGMSALSDLHAFNLRVDPGNFKILSHNILVTLAIHFPSDFTPEVHIAVDKFLAVVSAALADKYR"[:20] //ykiss_Rainbow_trou
    // seq2 := "MHLTADDKKHIKAIWPSVAAHGDKYGGEALHRMFMCAPKTKTYFPDFDFSEHSKHILAHGKKVSDALNEACNHLDNIAGCLSKLSDLHAYDLRVDPGNFPLLAHQILVVVAIHFPKQFDPATHKALDKFLVSVSNVLTSKYR"[:20] //Xenopus_tropicalis_Western_clawed_frog
    // seq1 := "YRQSMNNFQGLRSFGCRFGTCTVQKLAHQIYQFTDKDKDNVAPRSKISPQGY" //adrenomedullin
    // seq2 := "TQAQLLRVGCVLGTCQVQNLSHRLWQLMGPAGRQDSAPVDPSSPHSY" //adm2
    // alignedSeq1,alignedSeq2,alignmentScore:= Needleman(seq1,seq2)
    // fmt.Println("alignedSeq1:",alignedSeq1)
    // fmt.Println("alignedSeq2:",alignedSeq2)
    // fmt.Println("Score:",alignmentScore)
    // distanceMap := ConstructDistanceMap([]string{seq1,seq2}) //what does this do? we only have two sequences?
    // fmt.Println("distanceMap:",distanceMap)
}

func ConvertTreeToNewick(tree Tree, root *Node) string {
	visited := make(map[*Node]bool)
	return buildNewickString(root, visited) + ";"
}

// Recursive helper to build the Newick string.
func buildNewickString(node *Node, visited map[*Node]bool) string {
	visited[node] = true
	var parts []string

	// Iterate over neighbors to recursively build Newick string.
	for _, neighbor := range node.neighbors {
		if !visited[neighbor] {
			subtree := buildNewickString(neighbor, visited)
			branch := subtree + ":" + strconv.FormatFloat(neighbor.distance, 'f', 6, 64)
			parts = append(parts, branch)
		}
	}

	// Extract sequence string from Sequences type as a slice representation.
	nodeLabel := extractSequenceSliceLabel(node.sequence)

	// Build Newick string for this node.
	if len(parts) > 0 {
		return "(" + strings.Join(parts, ",") + ")" + nodeLabel
	}
	return nodeLabel
}

// Extracts a slice representation of a node's sequences.
func extractSequenceSliceLabel(sequences Sequences) string {
	if len(sequences) > 0 {
		var sequenceStrings []string
		for _, seq := range sequences {
			sequenceStrings = append(sequenceStrings, seq.sequence)
		}
		return "[" + strings.Join(sequenceStrings, ",") + "]" // Join sequences into a slice representation.
	}
	return "[]" // Fallback if no sequence is available.
}

// ExportNewickToFile writes the Newick string to a specified file.
func ExportNewickToFile(newick string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(newick)
	return err
}

func MatrixToFloats(distanceMatrix [][]int) [][]float64{
    //we can clean this later- but take distanceMatrix and convert to [][]float64
    diffMatrix := make([][]float64, len(distanceMatrix))
    for i := 0; i < len(distanceMatrix); i++ {
        diffMatrix[i] = make([]float64, len(distanceMatrix))
        for j := 0; j < len(distanceMatrix); j++ {
            diffMatrix[i][j] = float64(-1*distanceMatrix[i][j])
        }
    }
    return diffMatrix
}

// ReadScoringMatrix reads in rows, columns, and values from a .csv file and stores them in a Matrix object
// input: sequenceType string, filename string, gapScore int
// output: Matrix object containing the values from the .csv file
func ReadScoringMatrix(sequenceType, filename string, gapScore int) (Matrix, error) {
    filepath := filepath.Join("ScoringMatrices", sequenceType, filename)
    file, err := os.Open(filepath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    matrix := make(Matrix)
    headers := records[0][1:] // Column labels

    // Add the gap column to headers
    headers = append(headers, "-")

    for _, row := range records[1:] {
        rowLabel := row[0]
        matrix[rowLabel] = make(map[string]int)
        for i, value := range row[1:] {
            intValue, err := strconv.Atoi(value)
            if err != nil {
                return nil, err
            }
            matrix[rowLabel][headers[i]] = intValue
        }
        // Add the gap score to each row
        matrix[rowLabel]["-"] = gapScore
    }

    // Add the gap row
    matrix["-"] = make(map[string]int)
    for _, header := range headers {
        if header == "-" {
            matrix["-"][header] = 0
        } else {
            matrix["-"][header] = gapScore
        }
    }

    return matrix, nil
}

//PrintSubMatrix uses fmt.Print and fmt.Println to print a Matrix object as a table with each value separated by a space 
//input: a Matrix object, representing a table with row labels represented by outer key and column labels represented by inner key
func PrintSubMatrix(matrix Matrix) {
    // Get the column labels and sort them
    var columns []string
    for col := range matrix {
        // fmt.Println(col)
        columns = append(columns, col)
    }
    sort.Strings(columns)

    // Print the column labels
    fmt.Print(" ")
    for _, col := range columns {
        fmt.Printf(" %s", col)
    }
    fmt.Println()

    // Print each row in alphabetical order
    for _, rowLabel := range columns {
        fmt.Printf("%s", rowLabel)
        for _, col := range columns {
            fmt.Printf(" %d", matrix[rowLabel][col])
        }
        fmt.Println()
    }
}

func ReadFASTAInput(sequenceType, folderName string) (Sequences, error) {
    dirPath := filepath.Join("Sequences", sequenceType, folderName)
    fmt.Println("directory path is", dirPath)
    files, err := os.ReadDir(dirPath)
    if err != nil {
        return nil, err
    }

    var allSequences Sequences

    for _, file := range files {
        if !file.IsDir() && strings.HasSuffix(file.Name(), ".txt") {
            filePath := filepath.Join(dirPath, file.Name())
            fileSequences, err := readFASTAFile(filePath)
            if err != nil {
                return nil, err
            }
            allSequences = append(allSequences, fileSequences...)
        }
        // fmt.Println(len(allSequences))
    }

    return allSequences, nil
}

func readFASTAFile(filePath string) (Sequences, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var fileSequences Sequences
    scanner := bufio.NewScanner(file)
    var currentSequence Sequence
    
    for scanner.Scan() {
        line := scanner.Text()
        // fmt.Println("new line is", line)
        if strings.HasPrefix(line, ">") {
            // if currentSequence.info != "" && currentSequence.sequence != "" {
            //     // fileSequences = append(fileSequences, currentSequence)
            //     // currentSequence = sequence{}
            // }
            // fmt.Println("saw >")
            currentSequence.info = line[1:]
            // fmt.Println("current sequence info:", currentSequence.info, "current sequence itself:", currentSequence.sequence)
        } else {
            currentSequence.sequence = line
            // fmt.Println("current sequence info:", currentSequence.info, "current sequence itself:", currentSequence.sequence)
            fileSequences = append(fileSequences, currentSequence)
            currentSequence = Sequence{}
        }
        // fmt.Println(currentSequence.info, currentSequence.sequence)
    }
    
    fmt.Println("folder has", len(fileSequences))

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return fileSequences, nil
}

func PrintSequencesList(sequences Sequences) {
    // fmt.Println("printing root")
    for i := 0; i < len(sequences); i++ {
        fmt.Println(sequences[i].sequence)
    }
}

// func main() {
//     inputSeqs, err := ReadFASTAInput("exampleType", "exampleFolder")
//     if err != nil {
//         fmt.Println("Error:", err)
//         return
//     }

//     for _, seq := range inputSeqs {
//         fmt.Printf("Info: %s, Sequence: %s\n", seq.info, seq.sequence)
//     }
// }

// //this is a helper function to read FASTA files
// //Input: the filepath 
// //Output: a map containing the name of the sequence and value cotaining the sequence and an error if any
// func readFASTAFile(filename string) (map[string]string, error) {
// 	sequences := make(map[string]string)
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()

// 	var id string
// 	var seq strings.Builder
// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		line := scanner.Text()
// 		if strings.HasPrefix(line, ">") {
// 			if id != "" {
// 				sequences[id] = seq.String()
// 				seq.Reset()
// 			}
// 			id = line[1:] // Remove ">" from the sequence ID
// 		} else {
// 			seq.WriteString(line)
// 		}
// 	}
// 	if id != "" {
// 		sequences[id] = seq.String()
// 	}
// 	return sequences, scanner.Err()
// }