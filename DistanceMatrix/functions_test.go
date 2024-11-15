//Blake Kiefer (bkiefer) Fall 2024 Final Project with Talha Ahmad Khan, William Lu, and Zoey

package main

import (
	"fmt"
	"strconv"
	"os"
	"testing"
	"bufio"
	"io/fs"
	"strings"
	"math"
)

//Test Type Declarations:
type IdentifyShorterSequenceTest struct {
	//input:
	sequence1 string
	sequence2 string

	//output:
	shortSeqLen int
}

type InitializeDistanceMatrixTest struct {
	//input:
	n int

	//output:
	distanceMatrix [][]float64
}

type SequenceSimilarityTest struct {
	//input:
	sequence1 string
	sequence2 string

	//output:
	similarity float64
}

type CreateDistanceMatrixTest struct {
	//input:
	sequences []string

	//output:
	distanceMatrix [][]float64
}

//Tests:

//IdentifyShorterSequence() test:
//Input: TestIdentifyShorterSequence() takes nothing as input besides a command.
//Output: Tells you if your function has passed all of the IdentifyShorterSequence test cases.
func TestIdentifyShorterSequence(t *testing.T) {
	//Reading inputs and resulting outputs of IdentifyShorterSequences test cases.
	tests := ReadIdentifyShorterSequenceTest("Tests/IdentifyShorterSequence/")
	//Entering the inputs from the test cases into our function and making sure the outputs of this function are 
	//equal to the outputs of the test cases.
	for i, test := range tests {
		functionAnswer := IdentifyShorterSequence(test.sequence1, test.sequence2)

		if functionAnswer != test.shortSeqLen {
			t.Errorf("Failed test TestIdentifyShorterSequence")
			fmt.Println("Test Case", strconv.Itoa(i), ": Expected ", strconv.Itoa(test.shortSeqLen), " but got ", strconv.Itoa(functionAnswer), ".")
		}
	}
}

//Input: ReadIdentifyShorterSequenceTest() takes as input a directory of strings ("Test/IdentifyShorterSequence/")
//Output: Returns an array of []IdentifyShorterSequenceTest.
func ReadIdentifyShorterSequenceTest(directory string) []IdentifyShorterSequenceTest {
	//Reading in the input files
	inputFiles := ReadDirectory(directory + "/Inputs")
	numFiles := len(inputFiles)
	tests := make([]IdentifyShorterSequenceTest, numFiles)

	//Assigning each input value in the test case to the correct field of each IdentifyShorterSequenceTest in the array
	for i, inputFile := range inputFiles {
		//Create a list with each indices designated as the particular field of the IdentifyShorterSequenceTest
		inputs := ReadInput(directory + "Inputs/" + inputFile.Name())
		tests[i].sequence1 = inputs[0]
		tests[i].sequence2 = inputs[1]

	}

	//Making sure the number of outputs equals the number of test cases.
	outputFiles := ReadDirectory(directory + "Outputs/")
	if len(outputFiles) != numFiles {
		panic("Error: Number of input and output files do not match.")
	}

	//Assigning each output value in the test case to the correct field of the corresponding IdentifyShorterSequenceTest for
	//each IdentifyShorterSequenceTest object within the array.
	for i, outputFile := range outputFiles {
		//Create a list with each indices designated as the particular field of the IdentifyShorterSequenceTest
		outputs := ReadOutput(directory + "Outputs/" + outputFile.Name())
		output, err := strconv.Atoi(outputs[0])
		Check(err)

		tests[i].shortSeqLen = output
	}

	return tests
}

//InitializeDistanceMatrix() test:
//Input: TestInitializeDistanceMatrix() takes nothing as input besides a command.
//Output: Tells you if your function has passed all of the InitializeDistanceMatrix test cases.
func TestInitializeDistanceMatrix(t *testing.T) {
	//Reading inputs and resulting outputs of IdentifyShorterSequences test cases.
	tests := ReadInitializeDistanceMatrixTest("Tests/InitializeDistanceMatrix/")
	//Entering the inputs from the test cases into our function and making sure the outputs of this function are 
	//equal to the outputs of the test cases.
	for _, test := range tests {
		functionAnswer := InitializeDistanceMatrix(test.n)

		for r := 0; r < test.n; r ++ {
			for c := 0; c < test.n; c ++ {
				if functionAnswer[r][c] != test.distanceMatrix[r][c] {
					t.Errorf("Failed test InitializeDistanceMatrix")
				}
			}
		}
	}
}

//Input: ReadInitializeDistanceMatrixTest() takes as input a directory of strings ("Test/InitializeDistanceMatrix/")
//Output: Returns an array of []InitializeDistanceMatrixTest.
func ReadInitializeDistanceMatrixTest(directory string) []InitializeDistanceMatrixTest {
	//Reading in the input files
	inputFiles := ReadDirectory(directory + "/Inputs")
	numFiles := len(inputFiles)
	tests := make([]InitializeDistanceMatrixTest, numFiles)

	//Assigning each input value in the test case to the correct field of each InitializeDistanceMatrixTest in the array
	for i, inputFile := range inputFiles {
		//Create a list with each indices designated as the particular field of the InitializeDistanceMatrixTest
		inputs := ReadInput(directory + "Inputs/" + inputFile.Name())

		n, err := strconv.Atoi(inputs[0])
		Check(err)

		tests[i].n = n

	}

	//Making sure the number of outputs equals the number of test cases.
	outputFiles := ReadDirectory(directory + "Outputs/")
	if len(outputFiles) != numFiles {
		panic("Error: Number of input and output files do not match.")
	}

	//Assigning each output value in the test case to the correct field of the corresponding InitializeDistanceMatrixTest for
	//each InitializeDistanceMatrixTest object within the array.
	for i, outputFile := range outputFiles {
		//Create a list with each indices designated as the particular field of the InitializeDistanceMatrixTest
		lines := ReadOutput(directory + "Outputs/" + outputFile.Name())

		//Creating an emtpy 2D slice.
		output := make([][]float64, 0)
		//Add the appropriate number of rows to the 2D slice.  The appropriate number is the integer n which is the width and length of the final 2D array.
		for r := 0; r < tests[i].n; r ++ {
			//Get next line from output file
			line := lines[r]
			//Split the line at the spaces to get an array of 10 different string values
			rowVals := strings.Fields(line)

			//Add each value to the row.
			row := make([]float64, 0)
			//The number of values added is also equal to the number of rows because the final 2D array will be a square matrix.
			for c := 0; c < tests[i].n; c ++ {
				//Convert the current string to a float64
				val, err := strconv.ParseFloat(rowVals[c], 64)
				Check(err)

				//Add the value to the row.
				row = append(row, val)
			}

			//Finally add current row to the overall 2D slice
			output = append(output, row)
		}

		//The output is equal to the test distance matrix.  In this case all of the values should be 0.
		tests[i].distanceMatrix = output
	}

	return tests
}

//SequenceSimilarity() test:
//Input: Test=SequenceSimilarity() takes nothing as input besides a command.
//Output: Tells you if your function has passed all of the SequenceSimilarity test cases.
func TestSequenceSimilarity(t *testing.T) {
	//Reading inputs and resulting outputs of SequenceSimilarity test cases.
	tests := ReadSequenceSimilarityTest("Tests/SequenceSimilarity/")
	//Entering the inputs from the test cases into our function and making sure the outputs of this function are 
	//equal to the outputs of the test cases.
	for _, test := range tests {
		functionAnswer := SequenceSimilarity(test.sequence1, test.sequence2)

		if roundFloat(functionAnswer, 2) != test.similarity {
			t.Errorf("Failed test on TestSequenceSimilarity")
			fmt.Println("Expected ", strconv.FormatFloat(test.similarity, 'f', 2, 64), " but got ", strconv.FormatFloat(functionAnswer, 'f', 2, 64), ".")
		}
	}
}

//Input: ReadSequenceSimilarityTest() takes as input a directory of strings ("Test/SequenceSimilarity/")
//Output: Returns an array of []SequenceSimilarityTest.
func ReadSequenceSimilarityTest(directory string) []SequenceSimilarityTest {
	//Reading in the input files
	inputFiles := ReadDirectory(directory + "/Inputs")
	numFiles := len(inputFiles)
	tests := make([]SequenceSimilarityTest, numFiles)

	//Assigning each input value in the test case to the correct field of each SequenceSimilarityTest in the array
	for i, inputFile := range inputFiles {
		//Create a list with each indices designated as the particular field of the SequenceSimilarityTest
		inputs := ReadInput(directory + "Inputs/" + inputFile.Name())
		tests[i].sequence1 = inputs[0]
		tests[i].sequence2 = inputs[1]

	}

	//Making sure the number of outputs equals the number of test cases.
	outputFiles := ReadDirectory(directory + "Outputs/")
	if len(outputFiles) != numFiles {
		panic("Error: Number of input and output files do not match.")
	}

	//Assigning each output value in the test case to the correct field of the corresponding SequenceSimilarityTest for
	//each SequenceSimilarityTest object within the array.
	for i, outputFile := range outputFiles {
		//Create a list with each indices designated as the particular field of the SequenceSimilarityTest
		outputs := ReadOutput(directory + "Outputs/" + outputFile.Name())

		output, err := strconv.ParseFloat(outputs[0], 64)
		Check(err)

		tests[i].similarity = output
	}

	return tests
}

//CreateDistanceMatrix() test:
//Input: TestCreateDistanceMatrix() takes nothing as input besides a command.
//Output: Tells you if your function has passed all of the CreateDistanceMatrix test cases.
func TestCreateDistanceMatrix(t *testing.T) {
	//Reading inputs and resulting outputs of CreateDistanceMatrix test cases.
	tests := ReadCreateDistanceMatrixTest("Tests/CreateDistanceMatrix/")
	//Entering the inputs from the test cases into our function and making sure the outputs of this function are 
	//equal to the outputs of the test cases.
	for _, test := range tests {
		functionAnswer := CreateDistanceMatrix(test.sequences)

		for r := 0; r < len(functionAnswer); r ++ {
			for c := 0; c < len(functionAnswer[r]); c ++ {
				if roundFloat(functionAnswer[r][c], 2) != test.distanceMatrix[r][c] {
					t.Errorf("Failed test CreateDistanceMatrix")
				}
			}
		}
	}
}

//Input: CreateDistanceMatrixTest() takes as input a directory of strings ("Test/CreateDistanceMatrix/")
//Output: Returns an array of []CreateDistanceMatrixTest.
func ReadCreateDistanceMatrixTest(directory string) []CreateDistanceMatrixTest {
	//Reading in the input files
	inputFiles := ReadDirectory(directory + "/Inputs")
	numFiles := len(inputFiles)
	tests := make([]CreateDistanceMatrixTest, numFiles)

	//Assigning each input value in the test case to the correct field of each InitializeDistanceMatrixTest in the array
	for i, inputFile := range inputFiles {
		//Create a list with each indices designated as the particular field of the InitializeDistanceMatrixTest
		inputs := ReadInput(directory + "Inputs/" + inputFile.Name())

		sequences := make([]string, 0)
		for s := 0; s < len(inputs); s ++ {
			sequences = append(sequences, inputs[s])
		}

		tests[i].sequences = sequences

	}

	//Making sure the number of outputs equals the number of test cases.
	outputFiles := ReadDirectory(directory + "Outputs/")
	if len(outputFiles) != numFiles {
		panic("Error: Number of input and output files do not match.")
	}

	//Assigning each output value in the test case to the correct field of the corresponding CreateDistanceMatrixTest for
	//each CreateDistanceMatrixTest object within the array.
	for i, outputFile := range outputFiles {
		//Create a list with each indices designated as the particular field of the CreateDistanceMatrixTest
		lines := ReadOutput(directory + "Outputs/" + outputFile.Name())

		//Creating an emtpy 2D slice.
		output := make([][]float64, 0)
		//Add the appropriate number of rows to the 2D slice.  The appropriate number is the integer n which is the width and length of the final 2D array.
		for r := 0; r < len(lines); r ++ {
			//Get next line from output file
			line := lines[r]
			//Split the line at the spaces to get an array of n different string values
			rowVals := strings.Fields(line)

			//Add each value to the row.
			row := make([]float64, 0)
			//The number of values added is also equal to the number of rows because the final 2D array will be a square matrix.
			for c := 0; c < len(rowVals); c ++ {
				//Convert the current string to a float64
				val, err := strconv.ParseFloat(rowVals[c], 64)
				Check(err)

				//Add the value to the row.
				row = append(row, val)
			}

			//Finally add current row to the overall 2D slice
			output = append(output, row)
		}

		//The output is equal to the test distance matrix.  In this case all of the values should be 0.
		tests[i].distanceMatrix = output
	}

	return tests
}

//Other functions
//Input: ReadInput() takes as input an input file name
//Output: It returns a list where the value within each index corresponds to a particular field of a test case
//type designated within a ReadTestCase function.
func ReadInput(file string) []string {
	f, err := os.Open(file)
	Check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	//Using the scanner to read each individual line (each line has one value) and appending that value into a list
	//The values representing each field of the test are all ordered in the same way in each test case.
	inputs := make([]string, 0)
	for scanner.Scan() {
		input := scanner.Text()
		inputs = append(inputs, input)
	}

	return inputs
}

//Input: ReadOutput() function takes as input a file name (output)
//Output: It returns a list where the value within each index corresponds to a particular field of a test case
//type designated within a ReadTestCase function.
func ReadOutput(file string) []string {
	f, err := os.Open(file)
	Check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	//Using the scanner to read each individual line (each line has one value) and appending that value into a list
	//The values representing each field of the test are all ordered in the same way in each test case.
	outputs := make([]string, 0)
	for scanner.Scan() {
		output := scanner.Text()
		outputs = append(outputs, output)
	}

	return outputs
}

func ReadDirectory(directory string) []fs.DirEntry {
	files, err := os.ReadDir(directory)
	Check(err)
	return files
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}