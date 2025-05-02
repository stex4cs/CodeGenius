package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/stex4cs/CodeGenius/internal/analyzer"
	"github.com/stex4cs/CodeGenius/internal/documentation"
	"github.com/stex4cs/CodeGenius/internal/explanation"
	"github.com/stex4cs/CodeGenius/internal/optimization"
	"github.com/stex4cs/CodeGenius/pkg/config"
)

func main() {
	// Define command-line flags
	analyzeCmd := flag.Bool("analyze", false, "Analyze code and suggest improvements")
	docCmd := flag.Bool("doc", false, "Generate documentation for the code")
	explainCmd := flag.Bool("explain", false, "Explain complex code segments")
	optimizeCmd := flag.Bool("optimize", false, "Suggest code optimizations")
	configPath := flag.String("config", "", "Path to configuration file")
	filePath := flag.String("file", "", "Path to the file to process")
	dirPath := flag.String("dir", "", "Path to the directory to process")
	verbose := flag.Bool("verbose", false, "Enable verbose output")

	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Enable verbose mode if flag is set
	if *verbose {
		cfg.Verbose = true
	}

	// Check if file or directory is provided
	if *filePath == "" && *dirPath == "" {
		fmt.Println("Error: You must provide either a file (-file) or directory (-dir) to process")
		flag.Usage()
		os.Exit(1)
	}

	// Determine files to process
	var filesToProcess []string
	if *filePath != "" {
		filesToProcess = append(filesToProcess, *filePath)
	} else {
		files, err := getFilesInDirectory(*dirPath, cfg.FileExtensions)
		if err != nil {
			log.Fatalf("Error reading directory: %v", err)
		}
		filesToProcess = files
	}

	// Execute the requested command
	if *analyzeCmd {
		for _, file := range filesToProcess {
			if cfg.Verbose {
				fmt.Printf("Analyzing file: %s\n", file)
			}
			result, err := analyzer.AnalyzeCode(file, cfg)
			if err != nil {
				log.Printf("Error analyzing file %s: %v", file, err)
				continue
			}
			fmt.Println(result)
		}
	} else if *docCmd {
		for _, file := range filesToProcess {
			if cfg.Verbose {
				fmt.Printf("Generating documentation for file: %s\n", file)
			}
			result, err := documentation.GenerateDocumentation(file, cfg)
			if err != nil {
				log.Printf("Error generating documentation for file %s: %v", file, err)
				continue
			}
			fmt.Println(result)
		}
	} else if *explainCmd {
		for _, file := range filesToProcess {
			if cfg.Verbose {
				fmt.Printf("Explaining file: %s\n", file)
			}
			result, err := explanation.ExplainCode(file, cfg)
			if err != nil {
				log.Printf("Error explaining file %s: %v", file, err)
				continue
			}
			fmt.Println(result)
		}
	} else if *optimizeCmd {
		for _, file := range filesToProcess {
			if cfg.Verbose {
				fmt.Printf("Optimizing file: %s\n", file)
			}
			result, err := optimization.SuggestOptimizations(file, cfg)
			if err != nil {
				log.Printf("Error suggesting optimizations for file %s: %v", file, err)
				continue
			}
			fmt.Println(result)
		}
	} else {
		fmt.Println("No command specified. Please use one of the available commands.")
		flag.Usage()
		os.Exit(1)
	}
}

// getFilesInDirectory returns all files with the specified extensions in a directory and its subdirectories
func getFilesInDirectory(dir string, extensions []string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			for _, allowedExt := range extensions {
				if ext == allowedExt {
					files = append(files, path)
					break
				}
			}
		}
		return nil
	})

	return files, err
}