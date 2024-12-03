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


# Check if shiny is installed; if not, install it
if (!require("DT")) {
  install.packages("DT")
}
# Check if ggplot2 is installed; if not, install it
if (!require("ggplot2")) {
  install.packages("ggplot2")
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
      textInput("fileURL", "Or enter the URL of the sequencesa file:"),
      sliderInput("gapPenalty", "gapPenalty", 
                  min = -5, max = -1, value = -1),
      sliderInput("match", "match", 
                  min = 1, max = 5, value = 1),
      sliderInput("mismatch", "mismatch", 
                  min = -5, max = -1, value = -1),
      actionButton("runGoCode", "Run Sequence Alignment")
    ),
    mainPanel(imageOutput("image"))
    
  )
)



server <- function(input, output) {
  
  observeEvent(input$runGoCode, {
    # Determine the file path or URL to use
    filePath <- NULL
    
    if (!is.null(input$genomeFile)) {
      filePath <- input$genomeFile$datapath
    } else if (input$fileURL != "") {
      tempFile <- tempfile(fileext = ".fasta")
      download.file(input$fileURL, tempFile, mode = "wb")
      filePath <- tempFile
    }
    
    req(filePath)  # Ensure a file path or URL is provided
    
    # Compile the Go program
    compile_result <- system("go build", intern = TRUE)
    print(compile_result)  # For debugging, to see compile output
    
    # Run the compiled Go program
    print("Running go engine....")
    run_result <- system(paste("./MultipleSeqAlignment", filePath), intern = TRUE)
    print(run_result)  # For debugging, to see runtime output
    
    # Read the skew array from the CSV file
    print("Running msa.csv file")
    sequences_data <- read.csv("output/msa.csv", stringsAsFactors = FALSE)
    
    # Extract the sequences as a vector
    sequences <- sequences_data$Sequence
    print("Visualizing MSA")
    output$image <- renderImage( 
      { 
        list(src = "output/msa.png", height = "100%") 
      }, 
      deleteFile = FALSE 
    ) 
  
  })
}

# Run the application 
shinyApp(ui = ui, server = server)
