package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "path/filepath"
    "strconv"
)

// type ScoringMatrix map[string]map[string]int

func ReadScoringMatrix(filename string) (ScoringMatrix, error) {
    filepath := filepath.Join("ScoringMatrices", filename)
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

    matrix := make(ScoringMatrix)
    headers := records[0][1:] // Column labels

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
    }

    return matrix, nil
}

func main() {
    filename := "standardDNA.csv"
    matrix, err := ReadScoringMatrix(filename)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(matrix)
	fmt.Println("A, A", matrix["A"]["A"], "A, T", matrix["A"]["T"])
}
