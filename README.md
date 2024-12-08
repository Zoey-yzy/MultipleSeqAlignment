# MultipleSeqAlignment
project for Programming for Scientist 2024 fall @ cmu


<ins>Input arguments are formatted as:</ins>  
```
./MultipleSeqAlignment seqType distMethod filename gapScore folderName
```
An example command would be:
```
```
./MultipleSeqAlignment DNA alignedhomology standardDNA.csv -1 slidesexample
```
``` 
<ins>seqType</ins>: the type of sequence being aligned  
either "DNA" or "Protein"- case-sensitive

<ins>distMethod</ins>: the method used to construct the guide tree  
- "naivehomology"- calculated using sequence identity of unaligned pairs of sequences
- "alignedhomology" - calculated using sequence identity of Needleman-Wunsch-aligned pairs of sequences
- "alignmentscore" - pairwise Needleman-Wunsch of each sequence with every other to obtain an alignment score  

<ins>filename</ins>: string indicating the name of the .csv file to read a scoring lookup matrix from  
if seqType = "DNA", list of files is under "ScoringMatrices/DNA"  
if seqType = "Protein", list of files is under "ScoringMatrices/Protein"- e.g. "BLOSUM62.csv"

<ins>gapScore</ins>: the linear gap penalty to use  
for proteins, try gapScore values from -7 to -12; DNA/RNA -2 or -3
if you're doing simple scoring- (1 for match, -1 for mismatch)- do gapScore =  -1 or -2  

<ins>folderName</ins>: string indicating the name of the folder to access pre-loaded FASTA sequence files  
- if seqType = "DNA", list of files is under "Sequences/DNA"  
    mostly tests, we don't have anything here really

- if seqType = "Protein", list of files is under "Sequences/Protein"  
    contains tests and some real sequences:  
        "adrenomedullin"- amino acid sequences of human adrenomedullin and adrenomedullin-2 - length ~30 aa  
        "peptidehormones"- amino acid sequences of 9 different glucagon/secretin- related peptide hormones - length ~30-40 aa  
        "enkephalins"- amino acid sequences of human Met- and Leu-enkephalin - length 5 aa  
        "opioidpeptides"- amino acid sequences of 16 different opioid receptor peptide ligands - length 3-7 aa





