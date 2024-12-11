// Author: Zhaoyi You
// Andrew ID: zhaoyiy
// Date: 12-01-2024
package main 
import(
	"bufio"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"testing"
)
type FindMinDistTests struct{
	mtx [][]float64
	row int
	col int 
}
func FindMinDistTest(t *testing.T){
	tests := ReadFindMinDistTests("Tests/FindMinDist/")
	for _, test := range tests{ 
// each test is a test struct 
		answer1, answer2 := FindMinDist(test.mtx)
		if answer1 != test.row || answer2 != test.col {
			t.Errorf("FindMinDist(%v) = %v,%v, should be %v, %v", test.mtx, answer1, answer2, test.row, test.col)
		}
	}
}
func ReadFindMinDistTests(directory string)[]FindMinDistTests{
	inputFiles := ReadDirectory(directory + "/input")
	numFiles := len(inputFiles)
	tests := make([]FindMinDistTests, numFiles)
	for i, inputFile := range inputFiles{
		tests[i].mtx = Readmtx(directory + "input/" + inputFile.Name())
	}
	outputFiles := ReadDirectory(directory + "/output")
	if len(outputFiles) != numFiles{
		panic("Error: number of input and output files do not match!") 

	}
	for i, outputFile := range outputFiles{
		tests[i].row, tests[i].col = ReadInt(directory + "output/" + outputFile.Name())
	}
	return tests  
}

// Read a matrix from the file 
// The matrix is written in numbers that are separated by space

func Readmtx(file string)[][]float64{
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	mtx := make([][]float64,0)
	for scanner.Scan(){
		line := scanner.Text()
		newrow := make([]float64,0)
		parts := strings.Split(line," ")
		for _, val := range parts{
			v, err := strconv.ParseFloat(val, 64)
			if err != nil{
				panic(err)
			}
			newrow = append(newrow,v)
		}
		mtx = append(mtx,newrow)
	}
	return mtx 
}

// Read two integers from a file
// the two integers are separated by a space

func ReadInt(file string)(int, int){
	f, err := os.Open(file)
	if err != nil{
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	line := scanner.Text()
	parts := strings.Split(line, " ")
	row, err := strconv.Atoi(parts[0])
	if err != nil{
		panic(err)
	}
	col, err := strconv.Atoi(parts[1])
	if err != nil{
		panic(err)
	}
	return row, col 
}

// test function for calculate average divergence for each cluster
// calculate each row in matrix for their divergence with other clusters 

type CalcDivergenceTests struct{
	mtx [][]float64 
	result []float64 
}
func CalcDivergenceTest(t *testing.T){
	tests := ReadCalcDivergenceTests("Tests/CalcDivergence/")
	for _, test := range tests{
		answer := CalcDivergence(test.mtx)
		for j := range answer{
			if answer[j] != test.result[j]{
				t.Errorf("CalcDivergence(%v) = %v, should be %v", test.mtx, answer, test.result)
			}
		}
	}
}
func ReadCalcDivergenceTests(dir string)[]CalcDivergenceTests{
	inputFiles := ReadDirectory(dir + "input")
	numFiles := len(inputFiles)
	tests := make([]CalcDivergenceTests,numFiles)
	for i,inputFile := range inputFiles{
		tests[i].mtx = Readmtx(dir + "input/" + inputFile.Name())
	}
	outputFiles := ReadDirectory(dir + "output")
	if len(outputFiles) != numFiles{
		panic("Error: number of input and output files do not match!")
	}
	for i, outputFile := range outputFiles{
		tests[i].result = ReadSlice(dir + "output/" + outputFile.Name())
	}
	return tests 
}
func ReadSlice(file string)[]float64{
	f, err := os.Open(file)
	if err != nil{
		panic(err)
	}
	defer f.Close() 
	scanner := bufio.NewScanner(f)
	scanner.Scan() 
	line := scanner.Text() 
	parts := strings.Split(line, " ")
	div := make([]float64,0)
	for _, val := range parts{
		v, err := strconv.ParseFloat(val, 64)
		if err != nil{
			panic(err)
		}
		div = append(div,v)
	}
	return div 
}
type AdjustMatrixTests struct{
	mtx [][]float64 
	result [][]float64 
}
func AdjustMatrixTest(t *testing.T){
	tests := ReadAdjustMatrixTests("Tests/AdjustMatrix/")
	for _, test := range tests{
		answer := AdjustMatrix(test.mtx)
		if !IstheSame(answer, test.result){
			t.Errorf("AdjustMatrix(%v) = %v, should be %v", test.mtx, answer, test.result)
		}
	}
}
func ReadAdjustMatrixTests(dir string)[]AdjustMatrixTests{
	inputFiles := ReadDirectory(dir + "input")
	numFiles := len(inputFiles)
	tests := make([]AdjustMatrixTests, numFiles)
	for i, inputFile := range inputFiles{
		tests[i].mtx = Readmtx(dir + "input/" + inputFile.Name())
	}
	outputFiles := ReadDirectory(dir + "output")
	if len(outputFiles) != numFiles{
		panic("Error: number of input and output files do not match!")
	}
	for i, outputFile := range outputFiles{
		tests[i].result = Readmtx(dir + "output/" + outputFile.Name())
	}
	return tests 
}
type CalcDistTests struct{
	mtx [][]float64 
	idx1 int 
	idx2 int 
	result float64 
}
func CalcDistTest(t *testing.T){
	tests := ReadCalcDistTests("Tests/CalcDist/")
	for _, test := range tests{
		answer := CalcDist(test.mtx, test.idx1, test.idx2)
		if answer != test.result{
			t.Errorf("CalcDist(%v, %v, %v) = %v, should be %v", test.mtx, test.idx1, test.idx2, answer, test.result)
		}
	}
}
func ReadCalcDistTests(dir string)[]CalcDistTests{
	inputFiles := ReadDirectory(dir + "input")
	numFiles := len(inputFiles)
	tests := make([]CalcDistTests, numFiles)
	for i, inputFile := range inputFiles {
		tests[i].idx1, tests[i].idx2, tests[i].mtx = ReadMtxInt(dir + "input/" + inputFile.Name())
	}
	outputFiles := ReadDirectory(dir + "output")
	if len(outputFiles) != numFiles{
		panic("Error: number of input and output files do not match!")
	}
	for i, outputFile := range outputFiles{
		tests[i].result = ReadFloat(dir + "output/" + outputFile.Name())
	}
	return tests 
}

// define the read function for input and output variables, from files to variables

func ReadMtxInt(file string)(int, int, [][]float64){
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close() 
	scanner := bufio.NewScanner(f)
	mtx := make([][]float64,0)
	count := 0
	var idx1, idx2 int 
	for scanner.Scan(){
		if count == 0{
			line := scanner.Text()
			parts := strings.Split(line, " ")
			val1, err1 := strconv.Atoi(parts[0])
			if err1 != nil{
				panic(err1)
			}
			val2, err2 := strconv.Atoi(parts[1])
			if err2 != nil{
				panic(err2)
			}
			idx1, idx2 = val1, val2 
			count += 1
		}else{
			line := scanner.Text()
			parts := strings.Split(line, " ")
			row := make([]float64,len(parts))
			for i, val := range parts{
				row[i], err = strconv.ParseFloat(val, 64)
				if err != nil{
					panic(err)
				}
			}
			mtx = append(mtx,row)
		}
	}
	return idx1, idx2, mtx
}
func ReadFloat(file string)float64{
	f, err := os.Open(file)
	if err != nil{
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	line := scanner.Text()
	result, err := strconv.ParseFloat(line, 64)
	if err != nil{
		panic(err)
	}
	return result 
}
type AddRowColTests struct{
	row int 
	col int 
	mtx [][]float64 
	result [][]float64 
}
func AddRowColTest(t *testing.T){
	tests := ReadAddRowColTest("Tests/AddRowCol/")
	for _, test := range tests{
		answer := AddRowCol(test.row, test.col, test.mtx)
		if !IstheSame(answer, test.result){
			t.Errorf("AddRowCol(%v, %v, %v) = %v, should be %v", test.row, test.col, test.mtx, answer, test.result)
		}
	}
}
func ReadAddRowColTest(dir string)[]AddRowColTests{
	inputFiles := ReadDirectory(dir + "input")
	numFiles := len(inputFiles)
	tests := make([]AddRowColTests, numFiles)
	for i, inputFile := range inputFiles{
		tests[i].row, tests[i].col, tests[i].mtx = ReadMtxInt(dir + "input/" + inputFile.Name())
	}
	outputFiles := ReadDirectory(dir + "output")
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}
	for i, outputFile := range outputFiles{
		tests[i].result = Readmtx(dir + "output/" + outputFile.Name())
	}
	return tests
}
type DeleteRowColTests struct{
	mtx [][]float64 
	row int 
	col int
	result [][]float64  
}
func DeleteRowColTest(t *testing.T){
	tests := ReadDeleteRowColTest("Tests/DeleteRowCol/")
	for _, test := range tests{
		answer := DeleteRowCol(test.mtx, test.row, test.col)
		if !IstheSame(answer, test.result){
			t.Errorf("DeleteRowCol(%v, %v, %v) = %v, should be %v", test.mtx, test.row, test.col, answer, test.result)
		}
	}
}
func ReadDeleteRowColTest(dir string)[]DeleteRowColTests{
	inputFiles := ReadDirectory(dir + "input")
	numFiles := len(inputFiles)
	tests := make([]DeleteRowColTests, numFiles)
	for i, inputFile := range inputFiles{
		tests[i].row, tests[i].col, tests[i].mtx = ReadMtxInt(dir + "input/" + inputFile.Name())
	}
	outputFiles := ReadDirectory(dir + "output")
	if len(outputFiles) != numFiles {
		panic("Error: number of input and output files do not match!")
	}
	for i, outputFile := range outputFiles{
		tests[i].result = Readmtx(dir + "output/" + outputFile.Name())
	}
	return tests
}
func ReadDirectory(dir string) []fs.DirEntry {
	//read in all files in the given directory
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	return files
}