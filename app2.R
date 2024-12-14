# Install necessary packages if not already installed
if (!require("shiny")) install.packages("shiny")
if (!require("ape")) install.packages("ape")
if (!requireNamespace("BiocManager", quietly=TRUE)){
  install.packages("BiocManager")
  
}
if (!require("DECIPHER")) {
  BiocManager::install("DECIPHER")
  
}
if (!require("msa")){
  BiocManager::install("msa")
}
if (!require("here")){
  install.packages("here")
  
}

library(shiny)
library(ape)
library(DECIPHER)
library(here)
library(msa)

# Define the Shiny UI
ui <- fluidPage(
  titlePanel("MSA Visualization"),
  
  sidebarLayout(
    sidebarPanel(
      fileInput("sequenceFile", "Upload Sequence File",
                accept = c(".fasta",".txt")),
      selectInput("seqType", "Sequence Type",
                  choices = c("DNA" = "DNA",
                              "Protein" = "Protein")),
      selectInput("distanceMetric", "Distance Metric",
                  choices = c("Needleman" = "alignmentscore",
                              "SequenceHomology" = "naivehomology",
                              "AlignedSequenceHomology"="alignedhomology")),
      sliderInput("gapPenalty", "gapPenalty", 
                  min = -10, max = -1, value = -1),
      selectInput("scoring", "ScoringMatrix",
                  choices = c(
                    "BLAST" = "DNABLAST.csv",
                    "StandardDNA" = "standardDNA.csv",
                    "BLOSUM62" = "BLOSUM62.csv",
                    "PAM250" = "PAM250.csv")),
      
      selectInput("treeType", "Tree Layout",
                  choices = c("Phylogram" = "phylogram",
                              "Cladogram" = "cladogram",
                              "Unrooted" = "unrooted",
                              "Radial" = "radial")),
      selectInput("algorithm1", "Choose Alignment Algorithm 1",
                  choices = c("Progressive (ours)")),
      selectInput("algorithm2", "Choose Alignment Algorithm 2",
                  choices = c("ClustalW (msa)")),
      actionButton("runGoMSA", "Run MSA")
    ),
    
    mainPanel(
      p("Guide Tree"),
      plotOutput("guideTreePlot", height = "600px"),
      p("Phylogenetic Tree"),
      plotOutput("phyloTreePlot", height = "600px"),
      p("Alignment 1 Visualization"),
      uiOutput("alignedSeqPlot", height = "600px"),
      p("Alignment 2 Visualization"),
      uiOutput("alignedSeqPlot2", height = "600px"),
      
      
    )
  )
)

# Define the Shiny server
server <- function(input, output, session) {
  alignedSeqs2 <- reactiveVal(NULL)  # Store the alignment
  cwd <- here()
  print(paste("cwd",cwd))
  observe({
    if (input$seqType == "DNA") {
      updateSelectInput(session, "scoring", 
                        choices = c("BLAST" = "DNABLAST.csv","StandardDNA" = "standardDNA.csv"))
    } 
    else if (input$seqType == "Protein") {
      updateSelectInput(session, "scoring", 
                        choices = c("BLOSUM62" = "BLOSUM62.csv","PAM250" = "PAM250.csv"))
    }
  })
  
  
  observeEvent(input$runGoMSA, {
    
    # Determine the file path or URL to use
    filePath <- input$sequenceFile$datapath
    if (is.null(input$sequenceFile)){
      showModal(modalDialog(
        title = "Error",
        "Please upload a file before submitting.",
        easyClose = TRUE,
        footer = NULL
      ))
      print(paste("No input file provided...exiting event:",cwd))
      return("")
      
    }
    print(paste("Running go MSA:",cwd))
    
    alignedSeqs <- NULL
    subMatrixPath <- paste(cwd,"/ScoringMatrices/",input$seqType,"/","BLOSUM63.csv",sep="")
    print(subMatrixPath)
    msaOutputPath <- paste( cwd, "/output/msa.fasta",sep="")
    msaPlotOutputPath <- paste(cwd , "/www/msa.html",sep="")
    msaPlotOutputPath2 <- paste(cwd , "/www/msa2.html",sep="")
    guideTreeOutputPath <- paste(cwd, "/output/guidetree.newick",sep="")
    phyloTreeOutputPath <- paste(cwd,"/output/phylotree.newick",sep="")
    print(paste("phyloTreeOutputPath:",phyloTreeOutputPath))
    withProgress(message = "Running MSA in go....", value = 0, {
      
      goCommand1 <- "go build"
      
      goCommand2 <- "./MultipleSeqAlignment Protein alignmentscore BLOSUM62.csv -8 covidspikeprotein"
      goCommand3 = paste("./MultipleSeqAlignment",input$seqType,input$distanceMetric,input$scoring,input$gapPenalty,filePath)
      print(goCommand3)
      incProgress(0.25, detail = paste("Preparing go commands"))
      print(goCommand3)
      incProgress(0.25, detail = paste("Preparing to run go build"))
      run_result1 <- system(goCommand1, intern = TRUE)
      incProgress(0.25, detail = paste("Preparing to run go msa"))
      run_result2 <- system(goCommand3, intern = TRUE)
      incProgress(0.25, detail = paste("Ran go msa successfully"))
      print(run_result1)  # For debugging, to see runtime output
      print(run_result2)  # For debugging, to see runtime output
      # Reactive expression to load the tree from file
      # Perform alignment based on selected algorithm
      if (input$seqType == "Protein"){
        sequences <- readAAStringSet(filePath,format="fasta")
      }else{
        sequences <- readDNAStringSet(filePath,format="fasta")
        
      }
      
      substitution_matrixW <- ""
      if (input$scoring == "BLOSUM62.csv"){
        substitution_matrixW <- "blosum"

      }
      else if (input$scoring == "PAM250.csv"){
        substitution_matrixW <- "pam"

      }
      else{
        substitution_matrixW <- "default"
        
      }
      
      alignment <- switch(input$algorithm2,
                          #"ClustalOmega (msa)" = msa(sequences, "ClustalOmega",gapOpening = gapExtension=input$gapPenalty,substitutionMatrix=substitution_matrixOmega ),
                          "ClustalW (msa)" = msa(sequences, "ClustalW",gapOpening=input$gapPenalty,gapExtension=input$gapPenalty,cluster="nj",substitutionMatrix=substitution_matrixW ),
                          #"Muscle (msa)" = msa(sequences, "Muscle",cluster="nj",gapExtension=input$gapPenalty,)
                          )
      aligned_matrix <- as.character(alignment)  # Alignment matrix
      
      if (input$seqType == "Protein"){
        aligned <- AAStringSet(aligned_matrix)
      }else{
        aligned <- DNAStringSet(aligned_matrix)
        
      }
      alignedSeqs2(aligned)  # Store aligned sequences
      #alignment_result(alignment)  # Save alignment result
      
    })
    alignedSeqs <- reactive({
      req(msaOutputPath)  # Ensure file is uploaded

      if (input$seqType == "Protein"){
        sequences <- readAAStringSet(msaOutputPath,format="fasta")
      }else{
        sequences <- readDNAStringSet(msaOutputPath,format="fasta")
        
      }
      sequences
    })
    treeData <- reactive({
      req(guideTreeOutputPath)  # Ensure file is uploaded
      read.tree(guideTreeOutputPath)
    })
    # Render the tree plot
    output$guideTreePlot <- renderPlot({
      req(treeData())  # Ensure tree data is available
      
      tree <- treeData()
      
      # Assign node labels if missing
      if (is.null(tree$node.label)) {
        tree$node.label <- paste("Node", 1:tree$Nnode)
      }
      
      # Plot the tree with the selected layout and options
      plot.phylo(tree, 
                 type = input$treeType, 
                 show.tip.label = TRUE, 
                 show.node.label = FALSE, 
                 cex = 1, 
                 edge.width = 2)
    })
    
    output$alignedSeqPlot <- renderUI({
      req(alignedSeqs())
      
      # Save alignment plot as an HTML file
      BrowseSeqs(alignedSeqs(), htmlFile = msaPlotOutputPath, openURL = FALSE)
      
      # Display the HTML file in the Shiny app
      tags$iframe(src = "msa.html", width = "100%", height = "600px", frameborder = 0)
    })
    
    output$alignedSeqPlot2 <- renderUI({
      req(alignedSeqs2())
      
      # Save alignment plot as an HTML file
      BrowseSeqs(alignedSeqs2(), htmlFile = msaPlotOutputPath2, openURL = FALSE)
      
      # Display the HTML file in the Shiny app
      tags$iframe(src = "msa2.html", width = "100%", height = "600px", frameborder = 0)
    })
    
    
    # Reactive expression to load the tree from file
    phyloTreeData <- reactive({
      req(phyloTreeOutputPath)  # Ensure file is uploaded
      read.tree(phyloTreeOutputPath)
    })
    
    # Render the tree plot
    output$phyloTreePlot <- renderPlot({
      req(phyloTreeData())  # Ensure tree data is available
      
      tree <- phyloTreeData()
      
      # Assign node labels if missing
      if (is.null(tree$node.label)) {
        tree$node.label <- paste("Node", 1:tree$Nnode)
      }
      
      # Plot the tree with the selected layout and options
      plot.phylo(tree, 
                 type = input$treeType, 
                 show.tip.label = TRUE, 
                 show.node.label = FALSE, 
                 cex = 1, 
                 edge.width = 2)
    })
  })
}

# Run the Shiny app
shinyApp(ui = ui, server = server)