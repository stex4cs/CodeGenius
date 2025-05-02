package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/stex4cs/CodeGenius/pkg/config"
)

// Query sends a query to the LLM API and returns the response
func Query(prompt string, cfg *config.Config) (string, error) {
	// Check if API key is set
	if cfg.APIKey == "" {
		// Use local fallback mode if no API key
		return localFallback(prompt), nil
	}

	// Define the request structure based on the configured provider
	var requestBody []byte
	var err error
	switch cfg.LLMProvider {
	case "openai":
		requestBody, err = prepareOpenAIRequest(prompt, cfg)
	case "anthropic":
		requestBody, err = prepareAnthropicRequest(prompt, cfg)
	default:
		return "", fmt.Errorf("unsupported LLM provider: %s", cfg.LLMProvider)
	}

	if err != nil {
		return "", err
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second,
	}

	// Determine API endpoint
	apiEndpoint := getAPIEndpoint(cfg)

	// Create the HTTP request
	req, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.APIKey))

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API returned non-200 status code: %d, response: %s", resp.StatusCode, string(body))
	}

	// Parse the response based on the provider
	var result string
	switch cfg.LLMProvider {
	case "openai":
		result, err = parseOpenAIResponse(body)
	case "anthropic":
		result, err = parseAnthropicResponse(body)
	}

	if err != nil {
		return "", err
	}

	return result, nil
}

// prepareOpenAIRequest prepares the request body for OpenAI API
func prepareOpenAIRequest(prompt string, cfg *config.Config) ([]byte, error) {
	requestBody := map[string]interface{}{
		"model":      cfg.ModelName,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a code assistant. Analyze and improve code with detailed explanations.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": cfg.Temperature,
		"max_tokens":  cfg.MaxTokens,
	}

	return json.Marshal(requestBody)
}

// prepareAnthropicRequest prepares the request body for Anthropic API
func prepareAnthropicRequest(prompt string, cfg *config.Config) ([]byte, error) {
	requestBody := map[string]interface{}{
		"model":       cfg.ModelName,
		"prompt":      fmt.Sprintf("\n\nHuman: %s\n\nAssistant:", prompt),
		"temperature": cfg.Temperature,
		"max_tokens_to_sample": cfg.MaxTokens,
	}

	return json.Marshal(requestBody)
}

// parseOpenAIResponse parses the response from OpenAI API
func parseOpenAIResponse(responseBody []byte) (string, error) {
	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err := json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", fmt.Errorf("failed to parse API response: %v", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in API response")
	}

	return response.Choices[0].Message.Content, nil
}

// parseAnthropicResponse parses the response from Anthropic API
func parseAnthropicResponse(responseBody []byte) (string, error) {
	var response struct {
		Completion string `json:"completion"`
	}

	err := json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", fmt.Errorf("failed to parse API response: %v", err)
	}

	return response.Completion, nil
}

// getAPIEndpoint returns the appropriate API endpoint based on the provider
func getAPIEndpoint(cfg *config.Config) string {
	switch cfg.LLMProvider {
	case "openai":
		return "https://api.openai.com/v1/chat/completions"
	case "anthropic":
		return "https://api.anthropic.com/v1/complete"
	default:
		return ""
	}
}

// localFallback provides simple responses without using an API
// This is used when no API key is configured
func localFallback(prompt string) string {
	return "API key not configured. Please set your API key in the configuration file to use the AI-powered features.\n\n" +
		"You can still use the basic static analysis features without an API key."
}