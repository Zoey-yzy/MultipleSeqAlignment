References:
Start with 1987 CLUSTAL paper - https://link.springer.com/article/10.1007/BF02603120 
Then probably 1994 CLUSTALW paper - https://www.ncbi.nlm.nih.gov/pmc/articles/PMC308517/ 
(1970 Needleman-Wunsch- https://www.sciencedirect.com/science/article/pii/0022283670900574?via%3Dihub) 
Use Multiple sequence alignment modeling: methods and applications | Briefings in Bioinformatics | Oxford Academic to find additional algorithms

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
