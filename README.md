# CodeGenius

[![Go Report Card](https://goreportcard.com/badge/github.com/stex4cs/CodeGenius)](https://goreportcard.com/report/github.com/stex4cs/CodeGenius)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub stars](https://img.shields.io/github/stars/stex4cs/CodeGenius.svg)](https://github.com/stex4cs/CodeGenius/stargazers)

CodeGenius is an AI-powered code assistant tool designed to help developers analyze, document, explain, and optimize their code. By leveraging advanced language models, CodeGenius can automatically provide valuable insights and improvements for your codebase.

## Features

- **Code Analysis**: Identify code issues, bugs, and improvement opportunities
- **Documentation Generation**: Automatically generate comprehensive documentation
- **Code Explanation**: Get detailed explanations of complex code segments
- **Optimization Suggestions**: Receive suggestions for performance improvements
- **Multi-language Support**: Works with multiple programming languages (Go, Python, JavaScript, TypeScript, Java, and more)
- **Local Operation Mode**: Basic functionality works without API keys

## Installation

### Using Go Install

```bash
go install github.com/stex4cs/CodeGenius/cmd/codegenius@latest
From Source
bashgit clone https://github.com/stex4cs/CodeGenius.git
cd CodeGenius
go build -o codegenius ./cmd/codegenius
Quick Start

Set up your API key (optional but recommended for full functionality)
bashexport CODEGENIUS_API_KEY="your-api-key"

Run a code analysis
bashcodegenius -analyze -file path/to/your/file.go

Generate documentation
bashcodegenius -doc -file path/to/your/file.go

Explain complex code
bashcodegenius -explain -file path/to/your/file.go

Get optimization suggestions
bashcodegenius -optimize -file path/to/your/file.go


Configuration
CodeGenius looks for configuration in these locations (in order):

Path specified with -config flag
.codegenius.json in your home directory
codegenius.json in the current directory

You can create a default configuration file:
bashcodegenius -config ~/.codegenius.json
Example configuration:
json{
  "api_key": "your-api-key",
  "llm_provider": "openai",
  "model_name": "gpt-4",
  "temperature": 0.7,
  "max_tokens": 2000,
  "timeout_seconds": 30,
  "file_extensions": [".go", ".py", ".js", ".ts", ".java"],
  "ignore_paths": ["node_modules", "vendor", "dist"],
  "save_results": true,
  "output_dir": "codegenius_output",
  "verbose": false
}
Supported Languages
CodeGenius supports a wide range of programming languages, including:

Go
Python
JavaScript/TypeScript
Java
C/C++
C#
PHP
Ruby
And many more!

How It Works
CodeGenius uses a combination of static analysis and AI assistance to provide insights about your code:

Static Analysis: Basic checks are performed locally without requiring API access
AI Analysis: When an API key is configured, advanced language models are used to provide deeper insights

Local Mode
Without an API key, CodeGenius can still perform:

Basic static analysis
Language detection
Code structure analysis
Simple performance checks

For advanced features like detailed explanations and optimization suggestions, an API key is recommended.
Contributing
Contributions are welcome! Please feel free to submit a Pull Request.

Fork the repository
Create your feature branch (git checkout -b feature/amazing-feature)
Commit your changes (git commit -m 'Add some amazing feature')
Push to the branch (git push origin feature/amazing-feature)
Open a Pull Request

License
This project is licensed under the MIT License - see the LICENSE file for details.
Acknowledgments

Thanks to all the contributors who have helped shape CodeGenius
Inspired by the need for better tools to understand and improve code