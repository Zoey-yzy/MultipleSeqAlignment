package main

type Tree []*Node

type Node struct{
	neighbors []*Node  
	sequence Sequences //is a slice of Sequence objects
	distance float64 
}


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

//could have implemented code as a 2D matrix of Move structs, containing a value and direction
type  Move struct {
    direction string
    value float64
}