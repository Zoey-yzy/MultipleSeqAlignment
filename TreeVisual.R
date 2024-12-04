library(tidyverse)
library(ape)

tree <- ape::read.tree("c:\Users\William\Documents\GitHub\MultipleSeqAlignment\tree.newick")
plot(tree)
edgelabels(
  round(tree$edge.length, 2), # Round edge lengths for cleaner display
  col = "blue",               # Edge label color
  frame = "none",             # Remove frames around labels
  adj = c(0,1.2), 
  cex = 0.8                   # Font size of edge labels
)

# 6:12 --------------------------------------------------------------------


