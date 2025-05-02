package explanation

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/stex4cs/CodeGenius/internal/llm"
	"github.com/stex4cs/CodeGenius/pkg/config"
	"github.com/stex4cs/CodeGenius/pkg/utils"
)

// ExplainCode explains the complex parts of the code in the specified file
func ExplainCode(filePath string, cfg *config.Config) (string, error) {
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
		"Explain the following %s code in detail, focusing on how it works and any complex or non-obvious parts. Use simple language that a junior developer would understand:\n\n```%s\n%s\n```",
		language,
		language,
		string(content),
	)

	// Get explanation from LLM
	explanation, err := llm.Query(prompt, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to explain code: %v", err)
	}

	// Format the output
	fileName := filepath.Base(filePath)
	output := fmt.Sprintf("# Code Explanation for %s\n\n%s", fileName, explanation)

	// Save the explanation to a file if enabled in config
	if cfg.SaveResults {
		outputPath := utils.GetOutputPath(filePath, "explanation", ".md")
		err = ioutil.WriteFile(outputPath, []byte(output), 0644)
		if err != nil {
			return output, fmt.Errorf("failed to save explanation to file: %v", err)
		}
		output += fmt.Sprintf("\n\nExplanation saved to: %s", outputPath)
	}

	return output, nil
}

// IdentifyComplexParts identifies particularly complex parts of the code that might need explanation
func IdentifyComplexParts(filePath string, cfg *config.Config) ([]ComplexCodeSection, error) {
	// Read the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Get the file extension
	ext := filepath.Ext(filePath)
	language := utils.GetLanguageFromExtension(ext)

	// Prepare prompt for the LLM
	prompt := fmt.Sprintf(
		"Analyze the following %s code and identify the 3-5 most complex parts that would benefit from explanation. For each part, provide the line numbers or function names and a brief description of why it's complex:\n\n```%s\n%s\n```",
		language,
		language,
		string(content),
	)

	// Get analysis from LLM
	analysis, err := llm.Query(prompt, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to identify complex parts: %v", err)
	}

	// Parse the response to extract the complex parts
	// This is a simplified approach - in a real implementation, you would parse the response more carefully
	sections := parseComplexSections(analysis, language)

	return sections, nil
}

// ComplexCodeSection represents a complex section of code
type ComplexCodeSection struct {
	Description   string
	LineStart     int
	LineEnd       int
	FunctionName  string
	Complexity    string  // e.g., "High", "Medium", "Low"
	ComplexityScore float64 // 1-10 scale
}

// parseComplexSections parses the LLM response to extract complex code sections
// This is a simplified implementation - in a real implementation, you would parse the response more carefully
func parseComplexSections(analysis string, language string) []ComplexCodeSection {
	var sections []ComplexCodeSection

	// Split the analysis into lines
	lines := strings.Split(analysis, "\n")

	currentSection := ComplexCodeSection{}
	inSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)// Look for patterns that might indicate a section
		if strings.Contains(line, "Line") || strings.Contains(line, "Function") || strings.Contains(line, "Method") || 
		   (strings.Contains(line, ":") && (strings.Contains(line, "complex") || strings.Contains(line, "difficult"))) {
			
			// If we were already processing a section, add it to the list
			if inSection && currentSection.Description != "" {
				sections = append(sections, currentSection)
				currentSection = ComplexCodeSection{}
			}
			
			inSection = true
			
			// Try to extract line numbers
			lineStartIdx := strings.Index(line, "Line")
			if lineStartIdx != -1 {
				// Try to parse line range like "Line 10-15" or "Lines 10-15"
				lineRange := line[lineStartIdx:]
				var start, end int
				_, err := fmt.Sscanf(lineRange, "Line %d-%d", &start, &end)
				if err != nil {
					_, err = fmt.Sscanf(lineRange, "Lines %d-%d", &start, &end)
				}
				
				if err == nil {
					currentSection.LineStart = start
					currentSection.LineEnd = end
				} else {
					// Try to parse a single line like "Line 10"
					var lineNum int
					_, err := fmt.Sscanf(lineRange, "Line %d", &lineNum)
					if err == nil {
						currentSection.LineStart = lineNum
						currentSection.LineEnd = lineNum
					}
				}
			}
			
			// Try to extract function name
			funcIdx := strings.Index(line, "Function")
			if funcIdx != -1 {
				parts := strings.Split(line[funcIdx:], " ")
				if len(parts) >= 2 {
					currentSection.FunctionName = strings.Trim(parts[1], ":`(),.;")
				}
			}
			
			// Try to determine complexity level
			if strings.Contains(strings.ToLower(line), "high complexity") || 
			   strings.Contains(strings.ToLower(line), "very complex") {
				currentSection.Complexity = "High"
				currentSection.ComplexityScore = 8.0
			} else if strings.Contains(strings.ToLower(line), "medium complexity") || 
			         strings.Contains(strings.ToLower(line), "moderately complex") {
				currentSection.Complexity = "Medium"
				currentSection.ComplexityScore = 5.0
			} else if strings.Contains(strings.ToLower(line), "low complexity") || 
			         strings.Contains(strings.ToLower(line), "slightly complex") {
				currentSection.Complexity = "Low"
				currentSection.ComplexityScore = 3.0
			} else {
				// Default to medium if not specified
				currentSection.Complexity = "Medium"
				currentSection.ComplexityScore = 5.0
			}
			
			// Start accumulating the description
			currentSection.Description = line
		} else if inSection {
			// Continue accumulating the description
			currentSection.Description += "\n" + line
		}
	}
	
	// Don't forget to add the last section if we were processing one
	if inSection && currentSection.Description != "" {
		sections = append(sections, currentSection)
	}
	
	return sections
}

// ExplainAlgorithm generates a detailed explanation of a specific algorithm or complex function
func ExplainAlgorithm(filePath string, functionName string, cfg *config.Config) (string, error) {
	// Read the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}
	
	// Extract the specific function if functionName is provided
	codeToExplain := string(content)
	if functionName != "" {
		// This is a simplified approach - in a real implementation, you would use proper parsing
		lines := strings.Split(string(content), "\n")
		var functionLines []string
		inFunction := false
		bracketCount := 0
		
		for _, line := range lines {
			if !inFunction {
				// Look for function definition
				if strings.Contains(line, "func "+functionName) || 
				   strings.Contains(line, "function "+functionName) || 
				   strings.Contains(line, "def "+functionName) {
					inFunction = true
					functionLines = append(functionLines, line)
					bracketCount += countChars(line, '{')
					bracketCount -= countChars(line, '}')
					continue
				}
			} else {
				functionLines = append(functionLines, line)
				bracketCount += countChars(line, '{')
				bracketCount -= countChars(line, '}')
				
				// Check if we've reached the end of the function
				if bracketCount <= 0 && (strings.TrimSpace(line) == "}" || strings.TrimSpace(line) == "end") {
					break
				}
			}
		}
		
		if len(functionLines) > 0 {
			codeToExplain = strings.Join(functionLines, "\n")
		}
	}
	
	// Get the file extension
	ext := filepath.Ext(filePath)
	language := utils.GetLanguageFromExtension(ext)
	
	// Prepare prompt for the LLM
	prompt := fmt.Sprintf(
		"Explain the following %s code algorithm in detail. Break down how it works step by step, explain any complex logic, and include relevant computer science concepts:\n\n```%s\n%s\n```",
		language,
		language,
		codeToExplain,
	)
	
	// Get explanation from LLM
	explanation, err := llm.Query(prompt, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to explain algorithm: %v", err)
	}
	
	// Format the output
	fileName := filepath.Base(filePath)
	var output string
	if functionName != "" {
		output = fmt.Sprintf("# Algorithm Explanation: %s in %s\n\n%s", functionName, fileName, explanation)
	} else {
		output = fmt.Sprintf("# Algorithm Explanation for %s\n\n%s", fileName, explanation)
	}
	
	// Save the explanation to a file if enabled in config
	if cfg.SaveResults {
		var outputName string
		if functionName != "" {
			outputName = fmt.Sprintf("algorithm_%s", functionName)
		} else {
			outputName = "algorithm"
		}
		outputPath := utils.GetOutputPath(filePath, outputName, ".md")
		err = ioutil.WriteFile(outputPath, []byte(output), 0644)
		if err != nil {
			return output, fmt.Errorf("failed to save algorithm explanation to file: %v", err)
		}
		output += fmt.Sprintf("\n\nExplanation saved to: %s", outputPath)
	}
	
	return output, nil
}

// countChars counts occurrences of a character in a string
func countChars(s string, c byte) int {
	count := 0
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			count++
		}
	}
	return count
}