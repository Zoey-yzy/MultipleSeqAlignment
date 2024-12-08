if (!require("tidyverse")) {
  install.packages("tidyverse")
}
if (!require("ape")) {
  install.packages("ape")
}
tree <- ape::read.tree("output/tree.newick") 
plot(tree)
edgelabels(
  round(tree$edge.length, 2), # Round edge lengths for cleaner display
  col = "blue",               # Edge label color
  frame = "none",             # Remove frames around labels
  adj = c(0,1.2), 
  cex = 0.8                   # Font size of edge labels
)
