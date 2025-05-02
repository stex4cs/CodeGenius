package analyzer

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/stex4cs/CodeGenius/internal/llm"
	"github.com/stex4cs/CodeGenius/pkg/config"
	"github.com/stex4cs/CodeGenius/pkg/utils"
)

// AnalyzeCode analyzes the code in the specified file and suggests improvements
func AnalyzeCode(filePath string, cfg *config.Config) (string, error) {
	// Read the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Get the file extension
	ext := filepath.Ext(filePath)
	language := utils.GetLanguageFromExtension(ext)

	// Prepare prompt for the LLM
	prompt := fmt.Sprintf(
		"Analyze the following %s code and suggest improvements for readability, maintainability, and performance. Focus on best practices and potential bugs:\n\n```%s\n%s\n```",
		language,
		language,
		string(content),
	)

	// Get analysis from LLM
	analysis, err := llm.Query(prompt, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to analyze code: %v", err)
	}

	// Format the output
	fileName := filepath.Base(filePath)
	output := fmt.Sprintf("# Code Analysis for %s\n\n%s", fileName, analysis)

	// Save the analysis to a file if enabled in config
	if cfg.SaveResults {
		outputPath := utils.GetOutputPath(filePath, "analysis", ".md")
		err = ioutil.WriteFile(outputPath, []byte(output), 0644)
		if err != nil {
			return output, fmt.Errorf("failed to save analysis to file: %v", err)
		}
		output += fmt.Sprintf("\n\nAnalysis saved to: %s", outputPath)
	}

	return output, nil
}

// AnalyzeCodeQuality performs static code analysis to check for common issues
func AnalyzeCodeQuality(filePath string, cfg *config.Config) ([]string, error) {
	// Read the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	ext := filepath.Ext(filePath)
	language := utils.GetLanguageFromExtension(ext)

	var issues []string

	// Basic static analysis based on language
	switch language {
	case "go":
		issues = analyzeGoCode(string(content))
	case "python":
		issues = analyzePythonCode(string(content))
	case "javascript":
		issues = analyzeJavaScriptCode(string(content))
	default:
		issues = append(issues, fmt.Sprintf("Static analysis not implemented for %s", language))
	}

	return issues, nil
}

// Basic analyzers for different languages
// In a real implementation, these would use proper parsing and static analysis

func analyzeGoCode(content string) []string {
	var issues []string

	// Check for common Go issues
	if strings.Contains(content, "if err != nil {") && !strings.Contains(content, "return") {
		issues = append(issues, "Error handling pattern detected without return statement")
	}

	if strings.Contains(content, "fmt.Println") {
		issues = append(issues, "Consider using structured logging instead of fmt.Println")
	}

	return issues
}

func analyzePythonCode(content string) []string {
	var issues []string

	// Check for common Python issues
	if strings.Contains(content, "except:") {
		issues = append(issues, "Bare except clause detected - specify exceptions to catch")
	}

	if strings.Contains(content, "import *") {
		issues = append(issues, "Wildcard import detected - consider importing specific modules")
	}

	return issues
}

func analyzeJavaScriptCode(content string) []string {
	var issues []string

	// Check for common JavaScript issues
	if strings.Contains(content, "var ") {
		issues = append(issues, "Consider using let or const instead of var")
	}

	if strings.Contains(content, "==") {
		issues = append(issues, "Consider using === for strict equality comparisons")
	}

	return issues
}