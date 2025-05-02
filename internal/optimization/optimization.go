package optimization

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/stex4cs/CodeGenius/internal/llm"
	"github.com/stex4cs/CodeGenius/pkg/config"
	"github.com/stex4cs/CodeGenius/pkg/utils"
)

// SuggestOptimizations suggests optimizations for the code in the specified file
func SuggestOptimizations(filePath string, cfg *config.Config) (string, error) {
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
		"Suggest optimizations for the following %s code. Focus on performance improvements, reducing complexity, and making the code more efficient. Include code examples where appropriate:\n\n```%s\n%s\n```",
		language,
		language,
		string(content),
	)

	// Get optimizations from LLM
	optimizations, err := llm.Query(prompt, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to suggest optimizations: %v", err)
	}

	// Format the output
	fileName := filepath.Base(filePath)
	output := fmt.Sprintf("# Optimization Suggestions for %s\n\n%s", fileName, optimizations)

	// Save the optimizations to a file if enabled in config
	if cfg.SaveResults {
		outputPath := utils.GetOutputPath(filePath, "optimizations", ".md")
		err = ioutil.WriteFile(outputPath, []byte(output), 0644)
		if err != nil {
			return output, fmt.Errorf("failed to save optimizations to file: %v", err)
		}
		output += fmt.Sprintf("\n\nOptimizations saved to: %s", outputPath)
	}

	return output, nil
}

// AnalyzePerformance analyzes the performance of the code and identifies bottlenecks
func AnalyzePerformance(filePath string, cfg *config.Config) ([]PerformanceIssue, error) {
	// Read the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Get the file extension
	ext := filepath.Ext(filePath)
	language := utils.GetLanguageFromExtension(ext)

	var issues []PerformanceIssue

	// Basic performance analysis based on language
	switch language {
	case "go":
		issues = analyzeGoPerformance(string(content))
	case "python":
		issues = analyzePythonPerformance(string(content))
	case "javascript":
		issues = analyzeJavaScriptPerformance(string(content))
	default:
		issues = append(issues, PerformanceIssue{
			Description:  fmt.Sprintf("Performance analysis not implemented for %s", language),
			LineStart:    0,
			LineEnd:      0,
			Severity:     "Info",
			Optimization: fmt.Sprintf("Consider using a profiler specific to %s to analyze performance", language),
		})
	}

	// If we have an API key, use the LLM for more advanced analysis
	if cfg.APIKey != "" {
		// Prepare prompt for the LLM
		prompt := fmt.Sprintf(
			"Analyze the performance of the following %s code. Identify specific bottlenecks, inefficient algorithms, or resource-intensive operations. For each issue, specify the line numbers, severity, and suggest an optimization:\n\n```%s\n%s\n```",
			language,
			language,
			string(content),
		)

		// Get analysis from LLM
		analysis, err := llm.Query(prompt, cfg)
		if err == nil {
			// Parse the LLM response to extract performance issues
			// This is a simplified approach - in a real implementation, you would parse the response more carefully
			llmIssues := parsePerformanceIssues(analysis)
			issues = append(issues, llmIssues...)
		}
	}

	return issues, nil
}

// PerformanceIssue represents a performance issue in the code
type PerformanceIssue struct {
	Description  string
	LineStart    int
	LineEnd      int
	Severity     string  // "Critical", "High", "Medium", "Low", "Info"
	Impact       float64 // Estimated performance impact (0.0-1.0)
	Optimization string  // Suggested optimization
}

// Basic performance analyzers for different languages
// In a real implementation, these would use proper parsing and static analysis

func analyzeGoPerformance(content string) []PerformanceIssue {
	var issues []PerformanceIssue

	// Check for common Go performance issues
	if strings.Contains(content, "for i := range") && !strings.Contains(content, "for i, _ := range") {
		issues = append(issues, PerformanceIssue{
			Description:  "Efficient range loop usage",
			LineStart:    0, // In a real implementation, find the actual line number
			LineEnd:      0,
			Severity:     "Info",
			Impact:       0.1,
			Optimization: "Good use of range loops. Continue using this pattern for better performance.",
		})
	}

	if strings.Contains(content, "make([]") && !strings.Contains(content, "make([], 0)") {
		issues = append(issues, PerformanceIssue{
			Description:  "Pre-allocated slice",
			LineStart:    0,
			LineEnd:      0,
			Severity:     "Info",
			Impact:       0.2,
			Optimization: "Good use of pre-allocated slices. Continue using this pattern to reduce memory allocations.",
		})
	}

	if strings.Contains(content, "append(") && strings.Contains(content, "for") {
		issues = append(issues, PerformanceIssue{
			Description:  "Potential slice growth in loop",
			LineStart:    0,
			LineEnd:      0,
			Severity:     "Medium",
			Impact:       0.4,
			Optimization: "Consider pre-allocating the slice with an estimated capacity to avoid repeated reallocations during append operations in loops.",
		})
	}

	// Check for mutex usage in hot paths
	if strings.Contains(content, "Lock()") && strings.Contains(content, "Unlock()") && strings.Contains(content, "for") {
		issues = append(issues, PerformanceIssue{
			Description:  "Mutex in loop",
			LineStart:    0,
			LineEnd:      0,
			Severity:     "Medium",
			Impact:       0.5,
			Optimization: "Consider moving the mutex outside the loop if possible, or using a more granular locking strategy to reduce contention.",
		})
	}

	return issues
}

func analyzePythonPerformance(content string) []PerformanceIssue {
	var issues []PerformanceIssue

	// Check for common Python performance issues
	if strings.Contains(content, "for i in range") && !strings.Contains(content, "xrange") {
		issues = append(issues, PerformanceIssue{
			Description:  "Using range instead of xrange",
			LineStart:    0,
			LineEnd:      0,
			Severity:     "Low",
			Impact:       0.2,
			Optimization: "For Python 2, consider using xrange instead of range for better memory usage. For Python 3, range is already optimized.",
		})
	}

	if strings.Contains(content, "import numpy") || strings.Contains(content, "import np") {
		issues = append(issues, PerformanceIssue{
			Description:  "NumPy usage",
			LineStart:    0,
			LineEnd:      0,
			Severity:     "Info",
			Impact:       0.1,
			Optimization: "Good use of NumPy for numerical operations. Ensure vectorization is used where possible instead of Python loops.",
		})
	}

	if strings.Contains(content, "+= ") && strings.Contains(content, "for") && !strings.Contains(content, "join(") {
		issues = append(issues, PerformanceIssue{
			Description:  "String concatenation in loop",
			LineStart:    0,
			LineEnd:      0,
			Severity:     "Medium",
			Impact:       0.4,
			Optimization: "Consider using string.join() or a list comprehension instead of += for string concatenation in loops for better performance.",
		})
	}

	if strings.Contains(content, "list(") && strings.Contains(content, "map(") {
		issues = append(issues, PerformanceIssue{
			Description:  "Converting map to list",
			LineStart:    0,
			LineEnd:      0,
			Severity:     "Low",
			Impact:       0.2,
			Optimization: "In Python 3, map returns an iterator. Only convert to list if necessary, as it consumes memory.",
		})
	}

	return issues
}

func analyzeJavaScriptPerformance(content string) []PerformanceIssue {
	var issues []PerformanceIssue

	// Check for common JavaScript performance issues
	if strings.Contains(content, "for (let i = 0;") {
		issues = append(issues, PerformanceIssue{
			Description:  "Standard for loop",
			LineStart:    0,
			LineEnd:      0,
			Severity:     "Info",
			Impact:       0.1,
			Optimization: "Standard for loops are generally fast. For arrays, consider using forEach, map, or filter for cleaner code if performance is not critical.",
		})
	}

	if strings.Contains(content, "document.getElementsBy") && !strings.Contains(content, "document.querySelector") {
		issues = append(issues, PerformanceIssue{
			Description:  "Using older DOM methods",
			LineStart:    0,
			LineEnd:      0,
			Severity:     "Low",
			Impact:       0.2,
			Optimization: "Consider using querySelector or querySelectorAll for more flexible and concise DOM selection.",
		})
	}

	if strings.Contains(content, ".forEach") && strings.Contains(content, "new Array(") {
		issues = append(issues, PerformanceIssue{
			Description:  "Array creation in loop",
			LineStart:    0,
			LineEnd:      0,
			Severity:     "Medium",
			Impact:       0.3,
			Optimization: "Consider using Array.from() or the spread operator for cleaner and potentially more efficient array creation.",
		})
	}

	if strings.Contains(content, "setTimeout") && strings.Contains(content, "for") {
		issues = append(issues, PerformanceIssue{
			Description:  "setTimeout in loop",
			LineStart:    0,
			LineEnd:      0,
			Severity:     "Medium",
			Impact:       0.4,
			Optimization: "Creating many timeouts in a loop can cause performance issues. Consider using a single timeout and managing the timing within your code.",
		})
	}

	return issues
}

// parsePerformanceIssues parses the LLM response to extract performance issues
// This is a simplified implementation - in a real implementation, you would parse the response more carefully
func parsePerformanceIssues(analysis string) []PerformanceIssue {
	var issues []PerformanceIssue

	// Split the analysis into lines
	lines := strings.Split(analysis, "\n")

	currentIssue := PerformanceIssue{}
	inIssue := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Look for patterns that might indicate an issue
		if strings.Contains(line, "Issue") || strings.Contains(line, "Bottleneck") || strings.Contains(line, "Problem") || 
		   strings.Contains(line, "Line") && (strings.Contains(line, "performance") || strings.Contains(line, "slow") || 
		   strings.Contains(line, "inefficient")) {
			
			// If we were already processing an issue, add it to the list
			if inIssue && currentIssue.Description != "" {
				issues = append(issues, currentIssue)
				currentIssue = PerformanceIssue{}
			}
			
			inIssue = true
			
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
					currentIssue.LineStart = start
					currentIssue.LineEnd = end
				} else {
					// Try to parse a single line like "Line 10"
					var lineNum int
					_, err := fmt.Sscanf(lineRange, "Line %d", &lineNum)
					if err == nil {
						currentIssue.LineStart = lineNum
						currentIssue.LineEnd = lineNum
					}
				}
			}
			
			// Try to determine severity
			if strings.Contains(strings.ToLower(line), "critical") {
				currentIssue.Severity = "Critical"
				currentIssue.Impact = 0.9
			} else if strings.Contains(strings.ToLower(line), "high") {
				currentIssue.Severity = "High"
				currentIssue.Impact = 0.7
			} else if strings.Contains(strings.ToLower(line), "medium") {
				currentIssue.Severity = "Medium"
				currentIssue.Impact = 0.5
			} else if strings.Contains(strings.ToLower(line), "low") {
				currentIssue.Severity = "Low"
				currentIssue.Impact = 0.3
			} else {
				currentIssue.Severity = "Medium" // Default
				currentIssue.Impact = 0.5
			}
			
			// Start accumulating the description
			currentIssue.Description = line
		} else if inIssue && (strings.Contains(line, "Suggestion:") || strings.Contains(line, "Optimization:") || 
		                      strings.Contains(line, "Solution:") || strings.Contains(line, "Fix:")) {
			// This line contains the optimization suggestion
			colonIndex := strings.Index(line, ":")
			if colonIndex != -1 && colonIndex+1 < len(line) {
				currentIssue.Optimization = strings.TrimSpace(line[colonIndex+1:])
			} else {
				currentIssue.Optimization = line
			}
		} else if inIssue {
			// Continue accumulating the description
			currentIssue.Description += "\n" + line
		}
	}
	
	// Don't forget to add the last issue if we were processing one
	if inIssue && currentIssue.Description != "" {
		issues = append(issues, currentIssue)
	}
	
	return issues
}

// SuggestRefactoring suggests how to refactor the code for better maintainability
func SuggestRefactoring(filePath string, cfg *config.Config) (string, error) {
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
		"Suggest how to refactor the following %s code to improve maintainability, readability, and following best practices. Focus on code structure, naming, and design patterns. Include code examples where appropriate:\n\n```%s\n%s\n```",
		language,
		language,
		string(content),
	)

	// Get refactoring suggestions from LLM
	refactoring, err := llm.Query(prompt, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to suggest refactoring: %v", err)
	}

	// Format the output
	fileName := filepath.Base(filePath)
	output := fmt.Sprintf("# Refactoring Suggestions for %s\n\n%s", fileName, refactoring)

	// Save the refactoring suggestions to a file if enabled in config
	if cfg.SaveResults {
		outputPath := utils.GetOutputPath(filePath, "refactoring", ".md")
		err = ioutil.WriteFile(outputPath, []byte(output), 0644)
		if err != nil {
			return output, fmt.Errorf("failed to save refactoring suggestions to file: %v", err)
		}
		output += fmt.Sprintf("\n\nRefactoring suggestions saved to: %s", outputPath)
	}

	return output, nil
}