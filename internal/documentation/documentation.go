package documentation

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/stex4cs/CodeGenius/internal/llm"
	"github.com/stex4cs/CodeGenius/pkg/config"
	"github.com/stex4cs/CodeGenius/pkg/utils"
)

// GenerateDocumentation generates documentation for the code in the specified file
func GenerateDocumentation(filePath string, cfg *config.Config) (string, error) {
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
		"Generate comprehensive documentation for the following %s code. Include function descriptions, parameters, return values, and usage examples:\n\n```%s\n%s\n```",
		language,
		language,
		string(content),
	)

	// Get documentation from LLM
	docs, err := llm.Query(prompt, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to generate documentation: %v", err)
	}

	// Format the output
	fileName := filepath.Base(filePath)
	output := fmt.Sprintf("# Documentation for %s\n\n%s", fileName, docs)

	// Save the documentation to a file if enabled in config
	if cfg.SaveResults {
		outputPath := utils.GetOutputPath(filePath, "docs", ".md")
		err = ioutil.WriteFile(outputPath, []byte(output), 0644)
		if err != nil {
			return output, fmt.Errorf("failed to save documentation to file: %v", err)
		}
		output += fmt.Sprintf("\n\nDocumentation saved to: %s", outputPath)
	}

	return output, nil
}

// ExtractFunctions extracts functions/methods from the code
func ExtractFunctions(filePath string) ([]FunctionInfo, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	ext := filepath.Ext(filePath)
	language := utils.GetLanguageFromExtension(ext)

	var functions []FunctionInfo

	switch language {
	case "go":
		functions = extractGoFunctions(string(content))
	case "python":
		functions = extractPythonFunctions(string(content))
	case "javascript":
		functions = extractJavaScriptFunctions(string(content))
	}

	return functions, nil
}

// FunctionInfo represents information about a function or method
type FunctionInfo struct {
	Name       string
	Signature  string
	DocString  string
	LineNumber int
}

// Extract functions for different languages
// In a real implementation, these would use proper parsing

func extractGoFunctions(content string) []FunctionInfo {
	var functions []FunctionInfo

	// Simple regex for Go functions - in a real implementation, use a proper parser
	funcRegex := regexp.MustCompile(`(?m)^func\s+(\w+)([^{]+){`)
	matches := funcRegex.FindAllStringSubmatchIndex(content, -1)

	lines := strings.Split(content, "\n")
	lineCount := 0

	for i, match := range matches {
		if len(match) >= 4 {
			start := match[0]
			end := match[1]

			// Find the line number
			lineNumber := 0
			currentPos := 0
			for _, line := range lines {
				currentPos += len(line) + 1 // +1 for newline
				lineNumber++
				if currentPos > start {
					break
				}
			}

			// Extract the function name and signature
			fullMatch := content[match[0]:match[1]]
			nameMatch := content[match[2]:match[3]]
			sigMatch := content[match[4]:match[5]]

			// Extract doc string (comments above the function)
			docString := ""
			if lineNumber > 1 {
				for j := lineNumber - 2; j >= 0 && j >= lineNumber-5; j-- {
					line := strings.TrimSpace(lines[j])
					if strings.HasPrefix(line, "//") {
						docString = line[2:] + "\n" + docString
					} else if line == "" {
						continue
					} else {
						break
					}
				}
			}

			functions = append(functions, FunctionInfo{
				Name:       nameMatch,
				Signature:  fullMatch,
				DocString:  strings.TrimSpace(docString),
				LineNumber: lineNumber,
			})
		}
	}

	return functions
}

func extractPythonFunctions(content string) []FunctionInfo {
	var functions []FunctionInfo

	// Simple regex for Python functions - in a real implementation, use a proper parser
	funcRegex := regexp.MustCompile(`(?m)^def\s+(\w+)([^:]+):`)
	matches := funcRegex.FindAllStringSubmatchIndex(content, -1)

	lines := strings.Split(content, "\n")

	for _, match := range matches {
		if len(match) >= 4 {
			start := match[0]
			end := match[1]

			// Find the line number
			lineNumber := 0
			currentPos := 0
			for _, line := range lines {
				currentPos += len(line) + 1 // +1 for newline
				lineNumber++
				if currentPos > start {
					break
				}
			}

			// Extract the function name and signature
			fullMatch := content[match[0]:match[1]]
			nameMatch := content[match[2]:match[3]]
			sigMatch := content[match[4]:match[5]]

			// Look for docstring (triple-quoted string after function definition)
			docString := ""
			if lineNumber < len(lines) {
				for i := lineNumber; i < min(lineNumber+10, len(lines)); i++ {
					line := strings.TrimSpace(lines[i])
					if strings.HasPrefix(line, "\"\"\"") || strings.HasPrefix(line, "'''") {
						// Found the start of a docstring
						startQuote := line[:3]
						startIndex := i
						endIndex := -1

						// Find the end of the docstring
						for j := startIndex + 1; j < min(startIndex+20, len(lines)); j++ {
							if strings.Contains(lines[j], startQuote) {
								endIndex = j
								break
							}
						}

						if endIndex != -1 {
							// Extract the docstring
							for k := startIndex; k <= endIndex; k++ {
								docString += lines[k] + "\n"
							}
							docString = strings.Trim(docString, startQuote)
							docString = strings.TrimSpace(docString)
							break
						}
					}
				}
			}

			functions = append(functions, FunctionInfo{
				Name:       nameMatch,
				Signature:  fullMatch,
				DocString:  docString,
				LineNumber: lineNumber,
			})
		}
	}

	return functions
}

func extractJavaScriptFunctions(content string) []FunctionInfo {
	var functions []FunctionInfo

	// Simple regex for JavaScript functions - in a real implementation, use a proper parser
	funcRegex := regexp.MustCompile(`(?m)(function\s+(\w+)|\bconst\s+(\w+)\s*=\s*function|\blet\s+(\w+)\s*=\s*function|\bvar\s+(\w+)\s*=\s*function|\b(\w+)\s*:\s*function)`)
	matches := funcRegex.FindAllStringSubmatchIndex(content, -1)

	lines := strings.Split(content, "\n")

	for _, match := range matches {
		if len(match) >= 2 {
			start := match[0]
			end := match[1]

			// Find the line number
			lineNumber := 0
			currentPos := 0
			for _, line := range lines {
				currentPos += len(line) + 1 // +1 for newline
				lineNumber++
				if currentPos > start {
					break
				}
			}

			// Extract the function name and signature
			fullMatch := content[match[0]:match[1]]
			nameMatch := ""
			for i := 2; i < len(match); i += 2 {
				if match[i] != -1 && match[i+1] != -1 {
					nameMatch = content[match[i]:match[i+1]]
					break
				}
			}

			// Extract JSDoc comments (if any)
			docString := ""
			if lineNumber > 1 {
				// Look for JSDoc style comments
				inComment := false
				for j := lineNumber - 2; j >= 0 && j >= lineNumber-10; j-- {
					line := strings.TrimSpace(lines[j])
					if strings.HasPrefix(line, "*/") {
						inComment = true
					} else if strings.HasPrefix(line, "/**") {
						inComment = false
						break
					} else if inComment {
						commentText := strings.TrimSpace(line)
						if strings.HasPrefix(commentText, "*") {
							commentText = strings.TrimSpace(commentText[1:])
						}
						docString = commentText + "\n" + docString
					} else if line == "" {
						continue
					} else {
						break
					}
				}
			}

			functions = append(functions, FunctionInfo{
				Name:       nameMatch,
				Signature:  fullMatch,
				DocString:  strings.TrimSpace(docString),
				LineNumber: lineNumber,
			})
		}
	}

	return functions
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}