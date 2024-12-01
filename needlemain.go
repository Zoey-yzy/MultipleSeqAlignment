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

func main() {
    gapScore := -1 //hardcoded/read from commandline/RShiny- add to subMatrix
    //for proteins, try values from -7 to -12; DNA/RNA -2 or -3
    //if you're doing simple scoring- 1 for match, -1 for mismatch- do -1 or -2 
    
    //later change this to be doable through RShiny- folder browsing or dropdown menu
    //reading scoring matrix
    fmt.Println("reading scoring matrix")
    seqType := "Protein"
    fmt.Println("working with", seqType)
	filename := "BLOSUM62.csv"
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
    folderName := "adrenomedullin"
    seqs, err := ReadFASTAInput(seqType, folderName) //seqs is a Sequences object
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("printing final sequences, length", len(seqs))
    for _, seq := range seqs {
        //fmt.Printf("Info: %s, Sequence: %s\n", seq.info, seq.sequence)
        fmt.Println(seq.info, seq.sequence)
    }

    // //to test the alignments- use alignmenttest2 or alignmenttest folders, either in DNA or Protein
    // seqStack, score := NeedlemanWunsch(seqs[0:3], seqs[3:6], subMatrix)
    // PrintSequencesList(seqStack)
    // fmt.Println("score is", score)
    // //has out of bounds error

    //code for multiple NWs on multiple sequences
    // var distanceMatrix Matrix 
    // distanceMatrix := make(Matrix) 
    distanceMatrix := make([][]int, len(seqs)) //hold int scores- results of pairwise alignments
    seqMatrix := make([][]Sequences, len(seqs)) //hold Sequences objects of length 2- no need to redo alignments when filling guide tree nodes, just fill from here
    // var seqMatrix SequenceMatrix 
    for i := 0; i < len(seqs); i++ {
        distanceMatrix[i] = make([]int, len(seqs))
        seqMatrix[i] = make([]Sequences, len(seqs))
        for j := i+1; j < len(seqs); j++ {
            // fmt.Println("printing", seqs[i:i+1], seqs[j:j+1])
            alignedTwo, score := NeedlemanWunsch(seqs[i:i+1], seqs[j:j+1], subMatrix) //returns a Sequences object- the aligned two sequences, int the alignment score
            PrintSequencesList(alignedTwo)
            fmt.Println("alignment score is", score)
            distanceMatrix[i][j] = score*-1
            seqMatrix[i][j] = alignedTwo
            // distanceMatrix.UpdateDistMatrix(seqs[i].info, seqs[j].info, score)
            // innerMap := make(map[string]int)
            // fmt.Println("score passed into update function is", score)
            // innerMap[seqs[j].info] = score
            // distanceMatrix[seqs[i].info] = innerMap //why is the value 0?
            // fmt.Println("distance matrix looks like: ")
            fmt.Println(distanceMatrix)
            // PrintSubMatrix(distanceMatrix)
            // distanceMatrix[seqs[i].info][seqs[j].info] = score
            // seqMatrix[seqs[i].info][seqs[j].info] = alignedTwo
        }
    }

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
    for i := 0; i < len(sequences); i++ {
        fmt.Println(sequences[i].sequence)
    }
}

// func main() {
//     seqs, err := ReadFASTAInput("exampleType", "exampleFolder")
//     if err != nil {
//         fmt.Println("Error:", err)
//         return
//     }

//     for _, seq := range seqs {
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