# MultipleSeqAlignment
project for Programming for Scientist 2024 fall @ cmu
References:
Start with 1987 CLUSTAL paper - https://link.springer.com/article/10.1007/BF02603120 
Then probably 1994 CLUSTALW paper - https://www.ncbi.nlm.nih.gov/pmc/articles/PMC308517/ 
(1970 Needleman-Wunsch- https://www.sciencedirect.com/science/article/pii/0022283670900574?via%3Dihub) 
Use Multiple sequence alignment modeling: methods and applications | Briefings in Bioinformatics | Oxford Academic to find additional algorithms

<u>Input arguments are formatted as:</u>  
./MultipleSeqAlignment seqType algType filename gapScore folderName

<u>seqType</u>: the type of sequence being aligned  
either "DNA" or "Protein"- case-sensitive

<u>algType</u>: the method used to construct the guide tree  
either "alignment"- pairwise Needleman-Wunsch of each sequence with every other to obtain an alignment score  
or "identity"- calculated using sequence identity

<u>filename</u>: string indicating the name of the .csv file to read a scoring lookup matrix from  
if seqType = "DNA", list of files is under "ScoringMatrices/DNA"  
if seqType = "Protein", list of files is under "ScoringMatrices/Protein"- e.g. "BLOSUM62.csv"

<u>gapScore</u>: the linear gap penalty to use  
for proteins, try gapScore values from -7 to -12; DNA/RNA -2 or -3
if you're doing simple scoring- (1 for match, -1 for mismatch)- do gapScore =  -1 or -2  

<u>folderName</u>: string indicating the name of the folder to access pre-loaded FASTA sequence files  
if seqType = "DNA", list of files is under "Sequences/DNA"  
    mostly tests, we don't have anything here really

if seqType = "Protein", list of files is under "Sequences/Protein"  
    contains tests and some real sequences:  
        "adrenomedullin"- amino acid sequences of human adrenomedullin and adrenomedullin-2 - length ~30 aa  
        "peptidehormones"- amino acid sequences of 9 different glucagon/secretin- related peptide hormones - length ~30-40 aa  
        "enkephalins"- amino acid sequences of human Met- and Leu-enkephalin - length 5 aa  
        "opioidpeptides"- amino acid sequences of 16 different opioid receptor peptide ligands - length 3-7 aa


Scoring/substitution matrix sources:
DNA: https://www.stat.berkeley.edu/~hhuang/STAT141/STATC141-lecture7.pdf
i remember another set of slides but can't find

Proteins:
PAM250: original paper- Dayhoff 1978 https://profs.scienze.univr.it/liptak/ALBioinfo/files/pam1.pdf
https://rosalind.info/glossary/pam250/

BLOSUM62: original paper- Henikoff and Henikoff 1992 https://www.pnas.org/doi/pdf/10.1073/pnas.89.22.10915
https://rosalind.info/glossary/blosum62/

Mutation Distance: original paper- Fitch and Margoliash 1968 https://www.jstor.org/stable/1720651?seq=2
note- they apparently did this before standard single-letter amino acid abbreviations

additional explainers and info:
http://compbio.pbworks.com/w/page/16252911/Pairwise%20Sequence%20Alignment%20and%20Substitution%20Matrices
https://www.ibi.vu.nl/teaching/mnw_2year/2007/mnw2yr_lec8_2007.pdf
https://rob-p.github.io/CSE549F17/lectures/Lec10.pdf
http://www.binf.gmu.edu/vaisman/binf630/lec06s16.pdf 

also: we're using a linear gap penalty, but affine gap penalty is more widely used: https://rosalind.info/problems/gaff/
https://dept.math.lsa.umich.edu/~dburns/seqalmit2.pdf
typical linear gap penalty values: -2 or -3 for DNA, -7 to -12 for proteins
in affine, gap extension penalty can be -1 to -2
https://www.sciencedirect.com/topics/engineering/gap-penalty#:~:text=In%20this%20regard%2C%20the%20simplest,3%20in%20DNA%20sequence%20alignment. Pairwise Sequence Alignment
Miguel Rocha, Pedro G. Ferreira, in 
Bioinformatics Algorithms
, 2018

different versions of PAM/BLOSUM for different similarities of aligned sequences- e.g. PAM120 https://www.calstatela.edu/sites/default/files/pam120.txt


