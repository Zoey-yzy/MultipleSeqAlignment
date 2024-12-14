# MultipleSeqAlignment
project for Programming for Scientist 2024 fall @ cmu
Demo link: https://drive.google.com/drive/folders/1xEtGnSn3T2a_1ChPEXnvrB3qYVItbGt2?usp=sharing

<ins>Input arguments are formatted as:</ins>  
```
./MultipleSeqAlignment seqType distMethod scoringMatrixfilename gapScore inputSequencesFilepath
```
An example command would be:
```
./MultipleSeqAlignment Protein alignmentscore BLOSUM62.csv -8 Sequences/Protein/covidspikeprotein/spikeproteins2.txt
``` 
<ins>seqType</ins>: the type of sequence being aligned  
either "DNA" or "Protein"- case-sensitive
- Example 1: DNA
- Example 2: Protein

<ins>distMethod</ins>: the method used to construct the guide tree  
- "naivehomology"- calculated using sequence identity of unaligned pairs of sequences
- "alignedhomology" - calculated using sequence identity of Needleman-Wunsch-aligned pairs of sequences
- "alignmentscore" - pairwise Needleman-Wunsch of each sequence with every other to obtain an alignment score  
- Example 1:naivehomology
- Example 2:alignedhomology
- Example 3:alignmentscore

<ins>scoringMatrixfilename</ins>: string indicating the name of the .csv file to read a scoring lookup matrix from  
if seqType = "DNA", list of files is under "ScoringMatrices/DNA"  
if seqType = "Protein", list of files is under "ScoringMatrices/Protein"- e.g. "BLOSUM62.csv"
- Example 1: BLOSUM62.csv (if seqType = Protein)
- Example 2: BLASTDNA.csv  (if seqType = Protein)

<ins>gapScore</ins>: the linear gap penalty to use  
for proteins, try gapScore values from -7 to -12; DNA/RNA -2 or -3
if you're doing simple scoring- (1 for match, -1 for mismatch)- do gapScore =  -1 or -2 
- Example 1: -1
- Example 2 : -2 

<ins>inputSequencesFilepath</ins>: string indicating the path of the file to FASTA sequence files  
- Example1 : Sequences/DNA/slidesexample/example.txt
- Example2 : Sequences/Protein/covidspikeprotein/spikeproteins.txt
- if seqType = "DNA", list of files is under "Sequences/DNA"  
    mostly tests, we don't have anything here really

- if seqType = "Protein", list of files is under "Sequences/Protein"  
    contains tests and some real sequences:  
        "adrenomedullin"- amino acid sequences of human adrenomedullin and adrenomedullin-2 - length ~30 aa  
        "peptidehormones"- amino acid sequences of 9 different glucagon/secretin- related peptide hormones - length ~30-40 aa  
        "enkephalins"- amino acid sequences of human Met- and Leu-enkephalin - length 5 aa  
        "opioidpeptides"- amino acid sequences of 16 different opioid receptor peptide ligands - length 3-7 aa





