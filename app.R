#
# This is a Shiny web application. You can run the application by clicking
# the 'Run App' button above.
#
# Find out more about building applications with Shiny here:
#
#    https://shiny.posit.co/
#

# Check if shiny is installed; if not, install it
if (!require("shiny")) {
  install.packages("shiny")
}



# Check if ggplot2 is installed; if not, install it
if (!require("ggplot2")) {
  install.packages("ggplot2")
}

if (!require("tidyverse")) {
  install.packages("tidyverse")
}
if (!require("ape")) {
  install.packages("ape")
}




library(shiny)
library(ggplot2)


# First, set your working directory to source file location.

# Define UI
ui <- fluidPage(
  titlePanel("Multiple Sequence Alignment Go"),
  
  
  
  sidebarLayout(
    sidebarPanel(
      fileInput("Sequences", "Upload Sequences", accept = c(".txt", ".fasta", ".fa", "csv")),
      sliderInput("gapPenalty", "gapPenalty", 
                  min = -5, max = -1, value = -1),
      sliderInput("match", "match", 
                  min = 1, max = 5, value = 1),
      sliderInput("mismatch", "mismatch", 
                  min = -5, max = -1, value = -1),
      actionButton("runGoCode", "Run Sequence Alignment")
    ),
    mainPanel(plotOutput("outputPlot"))  # Display the plot in the app
    
    
  )
)



server <- function(input, output) {
  
  observeEvent(input$runGoCode, {
    # Determine the file path or URL to use
    filePath <- NULL
    
    if (!is.null(input$genomeFile)) {
      filePath <- input$genomeFile$datapath
    } 
    
    req(filePath)  # Ensure a file path or URL is provided
    
    # Compile the Go program
    compile_result <- system("go build", intern = TRUE)
    print(compile_result)  # For debugging, to see compile output
    
    # Run the compiled Go program
    print("Running go engine....")
    seqType <- "Protein"
    distanceMetric <- "alignedhomology"
    scoring<- "BLOSUM62.csv"
    gapPenalty <- input$gapPenalty
    inputFile < -"covidspikeprotein"
    run_result <- system(paste("./MultipleSeqAlignment", seqType, distanceMetric,scoring,gapPenalty,inputFile), intern = TRUE)
    print(run_result)  # For debugging, to see runtime output
    
    # Read the skew array from the CSV file
    print("Running msa.csv file")
    tree <- ape::read.tree("output/tree.newick") 
    treePlot <- plotTree(tree)
    output$outputPlot <- renderPlot({
    tree})
    
    # Extract the sequences as a vector
    print("Visualizing Tree")
    
   
  })
}

# Run the application 
shinyApp(ui = ui, server = server)
