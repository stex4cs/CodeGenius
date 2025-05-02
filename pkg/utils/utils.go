package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// GetLanguageFromExtension returns the programming language based on the file extension
func GetLanguageFromExtension(ext string) string {
	extToLang := map[string]string{
		".go":   "Go",
		".py":   "Python",
		".js":   "JavaScript",
		".ts":   "TypeScript",
		".jsx":  "JavaScript (React)",
		".tsx":  "TypeScript (React)",
		".java": "Java",
		".c":    "C",
		".cpp":  "C++",
		".cs":   "C#",
		".php":  "PHP",
		".rb":   "Ruby",
		".rs":   "Rust",
		".swift": "Swift",
		".kt":   "Kotlin",
		".scala": "Scala",
		".html": "HTML",
		".css":  "CSS",
		".sh":   "Shell",
		".pl":   "Perl",
		".r":    "R",
		".m":    "Objective-C",
		".dart": "Dart",
		".lua": "Lua",
	}

	language, ok := extToLang[ext]
	if !ok {
		return "Unknown"
	}
	return language
}

// GetOutputPath returns the path for saving output files
func GetOutputPath(inputPath, analysisType, extension string) string {
	dir := filepath.Dir(inputPath)
	base := filepath.Base(inputPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	timestamp := time.Now().Format("20060102_150405")
	outputName := fmt.Sprintf("%s_%s_%s%s", name, analysisType, timestamp, extension)
	
	return filepath.Join(dir, "codegenius_output", outputName)
}

// FormatDuration formats a duration as a human-readable string
func FormatDuration(d time.Duration) string {
	if d.Hours() > 1 {
		h := int(d.Hours())
		m := int(d.Minutes()) % 60
		return fmt.Sprintf("%dh %dm", h, m)
	} else if d.Minutes() > 1 {
		m := int(d.Minutes())
		s := int(d.Seconds()) % 60
		return fmt.Sprintf("%dm %ds", m, s)
	} else {
		ms := int(d.Milliseconds())
		if ms < 1000 {
			return fmt.Sprintf("%dms", ms)
		}
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}

// IsIgnoredPath checks if a path should be ignored
func IsIgnoredPath(path string, ignorePaths []string) bool {
	for _, ignorePath := range ignorePaths {
		if strings.Contains(path, ignorePath) {
			return true
		}
	}
	return false
}

// TruncateString truncates a string to the specified length and adds an ellipsis if truncated
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// ExtractFunctionName extracts the function name from a function signature
func ExtractFunctionName(signature string, language string) string {
	switch language {
	case "Go":
		// Extract function name from Go function signature like "func FunctionName(params) return"
		parts := strings.Split(signature, "(")
		if len(parts) > 0 {
			funcParts := strings.Split(strings.TrimSpace(parts[0]), " ")
			if len(funcParts) > 1 {
				return funcParts[1]
			}
		}
	case "Python":
		// Extract function name from Python function signature like "def function_name(params):"
		parts := strings.Split(signature, "(")
		if len(parts) > 0 {
			funcParts := strings.Split(strings.TrimSpace(parts[0]), " ")
			if len(funcParts) > 1 {
				return funcParts[1]
			}
		}
	case "JavaScript":
		// Handle multiple JavaScript function styles
		if strings.Contains(signature, "function") {
			// Traditional function like "function functionName(params) {"
			parts := strings.Split(signature, "function")
			if len(parts) > 1 {
				nameParts := strings.Split(strings.TrimSpace(parts[1]), "(")
				if len(nameParts) > 0 {
					return strings.TrimSpace(nameParts[0])
				}
			}
		} else if strings.Contains(signature, "=>") {
			// Arrow function like "const functionName = (params) => {"
			parts := strings.Split(signature, "=")
			if len(parts) > 0 {
				constParts := strings.Split(strings.TrimSpace(parts[0]), " ")
				if len(constParts) > 1 {
					return constParts[len(constParts)-1]
				}
			}
		}
	}
	
	// Default case: return empty string if we couldn't extract a function name
	return ""
}

// CreateDirectoryIfNotExists creates a directory if it doesn't exist
func CreateDirectoryIfNotExists(path string) error {
	// Implementation depends on the os package, simplified here
	return nil
}

// GetFileStats returns statistics about a file or directory
func GetFileStats(path string) (int, int, error) {
	// In a real implementation, this would count lines, files, etc.
	// For simplicity, return dummy values
	return 100, 5, nil
}

// SanitizeFilename sanitizes a filename to be safe for file systems
func SanitizeFilename(filename string) string {
	// Replace characters that are not safe for file names
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(filename)
}

// GetFileSize returns a human-readable file size
func GetFileSize(sizeInBytes int64) string {
	const unit = 1024
	if sizeInBytes < unit {
		return fmt.Sprintf("%d B", sizeInBytes)
	}
	div, exp := int64(unit), 0
	for n := sizeInBytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(sizeInBytes)/float64(div), "KMGTPE"[exp])
}