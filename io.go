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
	nodeLabel := node.label

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

//COMMENTS PLEASE YOU INCREDIBLY PLEASANT PERSON
//reads in FASTA format sequences from a folder no matter how many text files
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
            // fmt.Println("found a new sequence")
            if currentSequence.info != "" {
                // fmt.Println("let's store prev Sequence")
                fileSequences = append(fileSequences, currentSequence)
                currentSequence = Sequence{}
            }
            // if currentSequence.info != "" && currentSequence.sequence != "" {
            //     // fileSequences = append(fileSequences, currentSequence)
            //     // currentSequence = sequence{}
            // }
            // fmt.Println("saw >")
            currentSequence.info = line[1:]
            // fmt.Println("current sequence info:", currentSequence.info, "current sequence itself:", currentSequence.sequence)
        } else {
            // fmt.Println("amino acid sequence")
            currentSequence.sequence += line
            // fmt.Println("current sequence info:", currentSequence.info, "current sequence itself:", currentSequence.sequence)
            
        }
        // fmt.Println(currentSequence.info, currentSequence.sequence)
    }
    fileSequences = append(fileSequences, currentSequence)
    currentSequence = Sequence{}
    
    fmt.Println("folder has", len(fileSequences))

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return fileSequences, nil
}
// saveAsFasta saves a list of sequences to a FASTA file.
func WriteFASTAOutput(filename string, sequences Sequences) error {
	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write each sequence in FASTA format
	for _, seq := range sequences {
		// Write the identifier line
		_, err := file.WriteString(fmt.Sprintf(">%s\n", seq.info))
		if err != nil {
			return fmt.Errorf("failed to write ID: %w", err)
		}

		// Wrap and write the sequence content (80 characters per line)
		for i := 0; i < len(seq.sequence); i += 80 {
			end := i + 80
			if end > len(seq.sequence) {
				end = len(seq.sequence)
			}
			_, err := file.WriteString(seq.sequence[i:end] + "\n")
			if err != nil {
				return fmt.Errorf("failed to write sequence content: %w", err)
			}
		}
	}
	return nil
}

func PrintSequencesList(sequences Sequences) {
    // fmt.Println("printing root")
    for i := 0; i < len(sequences); i++ {
        fmt.Println(sequences[i].sequence)
    }
}