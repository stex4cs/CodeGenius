package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	// API settings
	APIKey       string  `json:"api_key"`
	LLMProvider  string  `json:"llm_provider"` // "openai" or "anthropic"
	ModelName    string  `json:"model_name"`
	Temperature  float64 `json:"temperature"`
	MaxTokens    int     `json:"max_tokens"`
	TimeoutSeconds int    `json:"timeout_seconds"`

	// File settings
	FileExtensions []string `json:"file_extensions"`
	IgnorePaths    []string `json:"ignore_paths"`
	SaveResults    bool     `json:"save_results"`
	OutputDir      string   `json:"output_dir"`

	// Application settings
	Verbose       bool     `json:"verbose"`

	// Feature flags
	EnableAnalysis      bool `json:"enable_analysis"`
	EnableDocumentation bool `json:"enable_documentation"`
	EnableExplanation   bool `json:"enable_explanation"`
	EnableOptimization  bool `json:"enable_optimization"`
}

// LoadConfig loads the configuration from a file or creates a default configuration
func LoadConfig(configPath string) (*Config, error) {
	// Create default configuration
	cfg := &Config{
		APIKey:      "",
		LLMProvider: "openai",
		ModelName:   "gpt-4",
		Temperature: 0.7,
		MaxTokens:   2000,
		TimeoutSeconds: 30,
		FileExtensions: []string{
			".go", ".py", ".js", ".ts", ".jsx", ".tsx", 
			".java", ".c", ".cpp", ".cs", ".php", ".rb",
		},
		IgnorePaths: []string{
			"node_modules", "vendor", "dist", "build", "__pycache__",
			".git", ".vscode", ".idea",
		},
		SaveResults: true,
		OutputDir:   "codegenius_output",
		Verbose:     false,
		EnableAnalysis:      true,
		EnableDocumentation: true,
		EnableExplanation:   true,
		EnableOptimization:  true,
	}

	// If config path is provided, load from file
	if configPath != "" {
		// Check if file exists
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			// Create default config file
			err = saveConfig(cfg, configPath)
			if err != nil {
				return nil, fmt.Errorf("failed to create default config file: %v", err)
			}
			fmt.Printf("Created default configuration file at %s\n", configPath)
		} else {
			// Load configuration from file
			data, err := ioutil.ReadFile(configPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read config file: %v", err)
			}

			err = json.Unmarshal(data, cfg)
			if err != nil {
				return nil, fmt.Errorf("failed to parse config file: %v", err)
			}
		}
	} else {
		// Try to find config in standard locations
		homeDir, err := os.UserHomeDir()
		if err == nil {
			// Check for config in home directory
			homeCfgPath := filepath.Join(homeDir, ".codegenius.json")
			if _, err := os.Stat(homeCfgPath); err == nil {
				data, err := ioutil.ReadFile(homeCfgPath)
				if err == nil {
					err = json.Unmarshal(data, cfg)
					if err == nil {
						return cfg, nil
					}
				}
			}
		}

		// Check for config in current directory
		localCfgPath := "codegenius.json"
		if _, err := os.Stat(localCfgPath); err == nil {
			data, err := ioutil.ReadFile(localCfgPath)
			if err == nil {
				err = json.Unmarshal(data, cfg)
				if err == nil {
					return cfg, nil
				}
			}
		}
	}

	// Check for API key in environment variable
	envAPIKey := os.Getenv("CODEGENIUS_API_KEY")
	if envAPIKey != "" {
		cfg.APIKey = envAPIKey
	}

	return cfg, nil
}

// saveConfig saves the configuration to a file
func saveConfig(cfg *Config, path string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Check if output directory exists, create if it doesn't
	if c.SaveResults && c.OutputDir != "" {
		if _, err := os.Stat(c.OutputDir); os.IsNotExist(err) {
			err = os.MkdirAll(c.OutputDir, 0755)
			if err != nil {
				return fmt.Errorf("failed to create output directory: %v", err)
			}
		}
	}

	// Validate LLM provider
	if c.APIKey != "" {
		switch c.LLMProvider {
		case "openai":
			if c.ModelName == "" {
				c.ModelName = "gpt-4"
			}
		case "anthropic":
			if c.ModelName == "" {
				c.ModelName = "claude-2"
			}
		default:
			return fmt.Errorf("invalid LLM provider: %s", c.LLMProvider)
		}
	}

	// Validate timeout
	if c.TimeoutSeconds <= 0 {
		c.TimeoutSeconds = 30
	}

	return nil
}