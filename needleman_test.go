package main

import (
    "bufio"
    "os"
    "strconv"
    "strings"
    "testing"
	//"fmt"
	"path/filepath"
)

func ReadGetNextCharsInput(filename string) (Sequences, Sequences, string, int, int, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, nil, "", 0, 0, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Scan()
    firstLine := strings.Split(scanner.Text(), " ")
    direction := firstLine[0]
    i, _ := strconv.Atoi(firstLine[1])
    j, _ := strconv.Atoi(firstLine[2])

    scanner.Scan()
    rowSeqsLen, _ := strconv.Atoi(scanner.Text())
    rowSeqs := make(Sequences, rowSeqsLen)
    for k := 0; k < rowSeqsLen; k++ {
        scanner.Scan()
        rowSeqs[k] = Sequence{sequence: scanner.Text()}
    }

    scanner.Scan()
    colSeqsLen, _ := strconv.Atoi(scanner.Text())
    colSeqs := make(Sequences, colSeqsLen)
    for k := 0; k < colSeqsLen; k++ {
        scanner.Scan()
        colSeqs[k] = Sequence{sequence: scanner.Text()}
    }

    return rowSeqs, colSeqs, direction, i, j, nil
}

func ReadGetNextCharsOutput(filename string) ([]string, []string, int, int, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, nil, 0, 0, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Scan()
    newRowChars := strings.Split(scanner.Text(), " ")

    scanner.Scan()
    newColChars := strings.Split(scanner.Text(), " ")

    scanner.Scan()
    lastLine := strings.Split(scanner.Text(), " ")
    i, _ := strconv.Atoi(lastLine[0])
    j, _ := strconv.Atoi(lastLine[1])

    return newRowChars, newColChars, i, j, nil
}

func TestGetNextChars(t *testing.T) {
    inputFiles, err := filepath.Glob("Tests/GetNextChars/input/input*.txt")
    if err != nil {
        t.Fatalf("Error reading input files: %v", err)
    }

    outputFiles, err := filepath.Glob("Tests/GetNextChars/output/output*.txt")
    if err != nil {
        t.Fatalf("Error reading output files: %v", err)
    }

    if len(inputFiles) != len(outputFiles) {
        t.Fatalf("Mismatch between number of input and output files")
    }

    for idx, inputFile := range inputFiles {
        rowSeqs, colSeqs, direction, i, j, err := ReadGetNextCharsInput(inputFile)
        if err != nil {
            t.Fatalf("Error reading input file %s: %v", inputFile, err)
        }

        expectedRowChars, expectedColChars, expectedI, expectedJ, err := ReadGetNextCharsOutput(outputFiles[idx])
        if err != nil {
            t.Fatalf("Error reading output file %s: %v", outputFiles[idx], err)
        }

        newRowChars, newColChars, newI, newJ := GetNextChars(rowSeqs, colSeqs, direction, i, j)

        if !equal(newRowChars, expectedRowChars) || !equal(newColChars, expectedColChars) || newI != expectedI || newJ != expectedJ {
            t.Errorf("Test failed for input file %s\nExpected newRowChars: %v, got: %v\nExpected newColChars: %v, got: %v\nExpected i, j: (%d, %d), got: (%d, %d)", 
                inputFile, expectedRowChars, newRowChars, expectedColChars, newColChars, expectedI, expectedJ, newI, newJ)
        }
    }
}

func equal(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }
    return true
}

// ReadComputeNWMaxInput reads the input for ComputeNWMax from a file.
func ReadComputeNWMaxInput(filename string) ([3]int, error) {
    file, err := os.Open(filename)
    if err != nil {
        return [3]int{}, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Scan()
    parts := strings.Split(scanner.Text(), " ")
    var list [3]int
    for i, part := range parts {
        list[i], err = strconv.Atoi(part)
        if err != nil {
            return [3]int{}, err
        }
    }

    return list, nil
}

// ReadComputeNWMaxOutput reads the expected output for ComputeNWMax from a file.
func ReadComputeNWMaxOutput(filename string) (int, error) {
    file, err := os.Open(filename)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Scan()
    return strconv.Atoi(scanner.Text())
}

func TestComputeNWMax(t *testing.T) {
    inputFiles, err := filepath.Glob("Tests/ComputeNWMax/input/input*.txt")
    if err != nil {
        t.Fatalf("Error reading input files: %v", err)
    }

    outputFiles, err := filepath.Glob("Tests/ComputeNWMax/output/output*.txt")
    if err != nil {
        t.Fatalf("Error reading output files: %v", err)
    }

    if len(inputFiles) != len(outputFiles) {
        t.Fatalf("Mismatch between number of input and output files")
    }

    for idx, inputFile := range inputFiles {
        list, err := ReadComputeNWMaxInput(inputFile)
        if err != nil {
            t.Fatalf("Error reading input file %s: %v", inputFile, err)
        }

        expectedMax, err := ReadComputeNWMaxOutput(outputFiles[idx])
        if err != nil {
            t.Fatalf("Error reading output file %s: %v", outputFiles[idx], err)
        }

        actualMax := ComputeNWMax(list)

        if actualMax != expectedMax {
            t.Errorf("Test failed for input file %s\nExpected max: %d, got: %d", inputFile, expectedMax, actualMax)
        }
    }
}

func ReadSumOfPairsInput(filename string) (Sequences, Sequences, int, int, Matrix, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, nil, 0, 0, nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Scan()
    indices := strings.Split(scanner.Text(), " ")
    rowInd, _ := strconv.Atoi(indices[0])
    colInd, _ := strconv.Atoi(indices[1])

    scanner.Scan()
    matrixInfo := strings.Split(scanner.Text(), " ")
    sequenceType := matrixInfo[0]
    matrixFilename := matrixInfo[1]
    gapScore, _ := strconv.Atoi(matrixInfo[2])

    subMatrix, err := ReadScoringMatrix(sequenceType, matrixFilename, gapScore)
    if err != nil {
        return nil, nil, 0, 0, nil, err
    }

    scanner.Scan()
    rowSeqsLen, _ := strconv.Atoi(scanner.Text())
    rowSeqs := make(Sequences, rowSeqsLen)
    for i := 0; i < rowSeqsLen; i++ {
        scanner.Scan()
        rowSeqs[i] = Sequence{sequence: scanner.Text()}
    }

    scanner.Scan()
    colSeqsLen, _ := strconv.Atoi(scanner.Text())
    colSeqs := make(Sequences, colSeqsLen)
    for i := 0; i < colSeqsLen; i++ {
        scanner.Scan()
        colSeqs[i] = Sequence{sequence: scanner.Text()}
    }

    return rowSeqs, colSeqs, rowInd, colInd, subMatrix, nil
}

func ReadSumOfPairsOutput(filename string) (int, error) {
    file, err := os.Open(filename)
    if err != nil {
        return 0, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Scan()
    return strconv.Atoi(scanner.Text())
}

func TestSumOfPairs(t *testing.T) {
    inputFiles, err := filepath.Glob("Tests/SumOfPairs/input/input*.txt")
    if err != nil {
        t.Fatalf("Error reading input files: %v", err)
    }

    outputFiles, err := filepath.Glob("Tests/SumOfPairs/output/output*.txt")
    if err != nil {
        t.Fatalf("Error reading output files: %v", err)
    }

    if len(inputFiles) != len(outputFiles) {
        t.Fatalf("Mismatch between number of input and output files")
    }

    for idx, inputFile := range inputFiles {
        rowSeqs, colSeqs, rowInd, colInd, subMatrix, err := ReadSumOfPairsInput(inputFile)
        if err != nil {
            t.Fatalf("Error reading input file %s: %v", inputFile, err)
        }

        expectedSopScore, err := ReadSumOfPairsOutput(outputFiles[idx])
        if err != nil {
            t.Fatalf("Error reading output file %s: %v", outputFiles[idx], err)
        }

        actualSopScore := SumOfPairs(rowSeqs, colSeqs, rowInd, colInd, subMatrix)

        if actualSopScore != expectedSopScore {
            t.Errorf("Test failed for input file %s\nExpected sopScore: %d, got: %d", inputFile, expectedSopScore, actualSopScore)
        }
    }
}

func ReadComputeNWMaxMoveInput(filename string) (Sequences, Sequences, int, int, Matrix, [][]int, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, nil, 0, 0, nil, nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Scan()
    indices := strings.Split(scanner.Text(), " ")
    rowInd, _ := strconv.Atoi(indices[0])
    colInd, _ := strconv.Atoi(indices[1])

    scanner.Scan()
    matrixInfo := strings.Split(scanner.Text(), " ")
    sequenceType := matrixInfo[0]
    matrixFilename := matrixInfo[1]
    gapScore, _ := strconv.Atoi(matrixInfo[2])

    subMatrix, err := ReadScoringMatrix(sequenceType, matrixFilename, gapScore)
    if err != nil {
        return nil, nil, 0, 0, nil, nil, err
    }

    scanner.Scan()
    rowSeqsLen, _ := strconv.Atoi(scanner.Text())
    rowSeqs := make(Sequences, rowSeqsLen)
    for i := 0; i < rowSeqsLen; i++ {
        scanner.Scan()
        rowSeqs[i] = Sequence{sequence: scanner.Text()}
    }

    scanner.Scan()
    colSeqsLen, _ := strconv.Atoi(scanner.Text())
    colSeqs := make(Sequences, colSeqsLen)
    for i := 0; i < colSeqsLen; i++ {
        scanner.Scan()
        colSeqs[i] = Sequence{sequence: scanner.Text()}
    }

    var matrix [][]int
    for scanner.Scan() {
        row := strings.Split(scanner.Text(), " ")
        intRow := make([]int, len(row))
        for i, val := range row {
            intRow[i], err = strconv.Atoi(val)
            if err != nil {
                return nil, nil, 0, 0, nil, nil, err
            }
        }
		//fmt.Println(intRow)
        matrix = append(matrix, intRow)
    }

    return rowSeqs, colSeqs, rowInd, colInd, subMatrix, matrix, nil
}

func ReadComputeNWMaxMoveOutput(filename string) (int, string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return 0, "", err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Scan()
    parts := strings.Split(scanner.Text(), " ")
    max, err := strconv.Atoi(parts[0])
    if err != nil {
        return 0, "", err
    }
    direction := parts[1]

    return max, direction, nil
}

func TestComputeNWMaxMove(t *testing.T) {
    inputFiles, err := filepath.Glob("Tests/ComputeNWMaxMove/input/input*.txt")
    if err != nil {
        t.Fatalf("Error reading input files: %v", err)
    }

    outputFiles, err := filepath.Glob("Tests/ComputeNWMaxMove/output/output*.txt")
    if err != nil {
        t.Fatalf("Error reading output files: %v", err)
    }

    if len(inputFiles) != len(outputFiles) {
        t.Fatalf("Mismatch between number of input and output files")
    }

    for idx, inputFile := range inputFiles {
        rowSeqs, colSeqs, rowInd, colInd, subMatrix, matrix, err := ReadComputeNWMaxMoveInput(inputFile)
        if err != nil {
            t.Fatalf("Error reading input file %s: %v", inputFile, err)
        }

        expectedMax, expectedDirection, err := ReadComputeNWMaxMoveOutput(outputFiles[idx])
        if err != nil {
            t.Fatalf("Error reading output file %s: %v", outputFiles[idx], err)
        }

        actualMax, actualDirection := ComputeNWMaxMove(matrix, rowSeqs, colSeqs, rowInd, colInd, subMatrix)

        if actualMax != expectedMax || actualDirection != expectedDirection {
            t.Errorf("Test failed for input file %s\nExpected max: %d, direction: %s\nGot max: %d, direction: %s", 
                inputFile, expectedMax, expectedDirection, actualMax, actualDirection)
        }
    }
}

func ReadInitializeNWMatrix2DInput(filename string) (Sequences, Sequences, Matrix, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, nil, nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Scan()
    matrixInfo := strings.Split(scanner.Text(), " ")
    sequenceType := matrixInfo[0]
    matrixFilename := matrixInfo[1]
    gapScore, _ := strconv.Atoi(matrixInfo[2])

    subMatrix, err := ReadScoringMatrix(sequenceType, matrixFilename, gapScore)
    if err != nil {
        return nil, nil, nil, err
    }

    scanner.Scan()
    rowSeqsLen, _ := strconv.Atoi(scanner.Text())
    rowSeqs := make(Sequences, rowSeqsLen)
    for i := 0; i < rowSeqsLen; i++ {
        scanner.Scan()
        rowSeqs[i] = Sequence{sequence: scanner.Text()}
    }

    scanner.Scan()
    colSeqsLen, _ := strconv.Atoi(scanner.Text())
    colSeqs := make(Sequences, colSeqsLen)
    for i := 0; i < colSeqsLen; i++ {
        scanner.Scan()
        colSeqs[i] = Sequence{sequence: scanner.Text()}
    }

    return rowSeqs, colSeqs, subMatrix, nil
}

func ReadInitializeNWMatrix2DOutput(filename string) ([][]int, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var matrix [][]int
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        row := strings.Split(scanner.Text(), " ")
        intRow := make([]int, len(row))
        for i, val := range row {
            intRow[i], err = strconv.Atoi(val)
            if err != nil {
                return nil, err
            }
        }
        matrix = append(matrix, intRow)
    }

    return matrix, nil
}

func TestInitializeNWMatrix2D(t *testing.T) {
    inputFiles, err := filepath.Glob("Tests/InitializeNWMatrix2D/input/input*.txt")
    if err != nil {
        t.Fatalf("Error reading input files: %v", err)
    }

    outputFiles, err := filepath.Glob("Tests/InitializeNWMatrix2D/output/output*.txt")
    if err != nil {
        t.Fatalf("Error reading output files: %v", err)
    }

    if len(inputFiles) != len(outputFiles) {
        t.Fatalf("Mismatch between number of input and output files")
    }

    for idx, inputFile := range inputFiles {
        rowSeqs, colSeqs, subMatrix, err := ReadInitializeNWMatrix2DInput(inputFile)
        if err != nil {
            t.Fatalf("Error reading input file %s: %v", inputFile, err)
        }

        expectedMatrix, err := ReadInitializeNWMatrix2DOutput(outputFiles[idx])
        if err != nil {
            t.Fatalf("Error reading output file %s: %v", outputFiles[idx], err)
        }

        actualMatrix := InitializeNWMatrix2D(rowSeqs, colSeqs, subMatrix)

        if !equalMatrices(actualMatrix, expectedMatrix) {
            t.Errorf("Test failed for input file %s\nExpected matrix: %v\nGot matrix: %v", 
                inputFile, expectedMatrix, actualMatrix)
        }
    }
}

func equalMatrices(a, b [][]int) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if len(a[i]) != len(b[i]) {
            return false
        }
        for j := range a[i] {
            if a[i][j] != b[i][j] {
                return false
            }
        }
    }
    return true
}

func ReadInitializeNWTracebackMatrix2DInput(filename string) (Sequences, Sequences, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Scan()
    rowSeqsLen, _ := strconv.Atoi(scanner.Text())
    rowSeqs := make(Sequences, rowSeqsLen)
    for i := 0; i < rowSeqsLen; i++ {
        scanner.Scan()
        rowSeqs[i] = Sequence{sequence: scanner.Text()}
    }

    scanner.Scan()
    colSeqsLen, _ := strconv.Atoi(scanner.Text())
    colSeqs := make(Sequences, colSeqsLen)
    for i := 0; i < colSeqsLen; i++ {
        scanner.Scan()
        colSeqs[i] = Sequence{sequence: scanner.Text()}
    }

    return rowSeqs, colSeqs, nil
}

func ReadInitializeNWTracebackMatrix2DOutput(filename string) ([][]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var matrix [][]string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        row := strings.Split(scanner.Text(), " ")
        for i, val := range row {
            if val == "-" {
                row[i] = ""
            }
        }
        matrix = append(matrix, row)
    }

    return matrix, nil
}

func TestInitializeNWTracebackMatrix2D(t *testing.T) {
    inputFiles, err := filepath.Glob("Tests/InitializeNWTracebackMatrix2D/input/input*.txt")
    if err != nil {
        t.Fatalf("Error reading input files: %v", err)
    }

    outputFiles, err := filepath.Glob("Tests/InitializeNWTracebackMatrix2D/output/output*.txt")
    if err != nil {
        t.Fatalf("Error reading output files: %v", err)
    }

    if len(inputFiles) != len(outputFiles) {
        t.Fatalf("Mismatch between number of input and output files")
    }

    for idx, inputFile := range inputFiles {
        rowSeqs, colSeqs, err := ReadInitializeNWTracebackMatrix2DInput(inputFile)
        if err != nil {
            t.Fatalf("Error reading input file %s: %v", inputFile, err)
        }

        expectedMatrix, err := ReadInitializeNWTracebackMatrix2DOutput(outputFiles[idx])
        if err != nil {
            t.Fatalf("Error reading output file %s: %v", outputFiles[idx], err)
        }

        actualMatrix := InitializeNWTracebackMatrix2D(rowSeqs, colSeqs)

        if !equalStringMatrices(actualMatrix, expectedMatrix) {
            t.Errorf("Test failed for input file %s\nExpected matrix: %v\nGot matrix: %v", 
                inputFile, expectedMatrix, actualMatrix)
        }
    }
}

func equalStringMatrices(a, b [][]string) bool {
    if len(a) != len(b) {
        return false
    }
    for i := range a {
        if len(a[i]) != len(b[i]) {
            return false
        }
        for j := range a[i] {
            if a[i][j] != b[i][j] {
                return false
            }
        }
    }
    return true
}