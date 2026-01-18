package initialize

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed templates/llms.txt
var llmsTxt string

//go:embed templates/CLAUDE.md
var claudeMd string

//go:embed templates/cursorrules
var cursorrules string

type AgentType string

const (
	AgentLLMs   AgentType = "llms"
	AgentClaude AgentType = "claude"
	AgentCursor AgentType = "cursor"
	AgentAll    AgentType = "all"
)

var ValidAgentTypes = []AgentType{AgentLLMs, AgentClaude, AgentCursor, AgentAll}

func IsValidAgentType(t string) bool {
	for _, valid := range ValidAgentTypes {
		if string(valid) == t {
			return true
		}
	}
	return false
}

type WriteMode int

const (
	ModeSkip WriteMode = iota
	ModeForce
	ModeAppend
)

type InitResult struct {
	File     string
	Created  bool
	Appended bool
	Skipped  bool
	Error    error
}

func Run(agentType string, mode WriteMode) ([]InitResult, error) {
	if !IsValidAgentType(agentType) {
		return nil, fmt.Errorf("invalid agent type: %s (valid: llms, claude, cursor, all)", agentType)
	}

	var results []InitResult

	switch AgentType(agentType) {
	case AgentLLMs:
		results = append(results, writeFile("llms.txt", llmsTxt, mode))
	case AgentClaude:
		results = append(results, writeFile("CLAUDE.md", claudeMd, mode))
	case AgentCursor:
		results = append(results, writeFile(".cursorrules", cursorrules, mode))
	case AgentAll:
		results = append(results, writeFile("llms.txt", llmsTxt, mode))
		results = append(results, writeFile("CLAUDE.md", claudeMd, mode))
		results = append(results, writeFile(".cursorrules", cursorrules, mode))
	}

	return results, nil
}

const appendSeparator = "\n\n# --- todo-rank agent guide (appended) ---\n\n"

func writeFile(filename, content string, mode WriteMode) InitResult {
	result := InitResult{File: filename}

	path := filepath.Join(".", filename)

	// Check if file exists
	existingContent, err := os.ReadFile(path)
	fileExists := err == nil

	if fileExists {
		switch mode {
		case ModeSkip:
			result.Skipped = true
			return result
		case ModeAppend:
			// Append new content to existing
			newContent := string(existingContent) + appendSeparator + content
			if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
				result.Error = err
				return result
			}
			result.Appended = true
			return result
		case ModeForce:
			// Fall through to overwrite
		}
	}

	// Create or overwrite file
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		result.Error = err
		return result
	}

	result.Created = true
	return result
}
