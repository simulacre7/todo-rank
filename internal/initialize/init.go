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

type InitResult struct {
	File    string
	Created bool
	Skipped bool
	Error   error
}

func Run(agentType string, force bool) ([]InitResult, error) {
	if !IsValidAgentType(agentType) {
		return nil, fmt.Errorf("invalid agent type: %s (valid: llms, claude, cursor, all)", agentType)
	}

	var results []InitResult

	switch AgentType(agentType) {
	case AgentLLMs:
		results = append(results, writeFile("llms.txt", llmsTxt, force))
	case AgentClaude:
		results = append(results, writeFile("CLAUDE.md", claudeMd, force))
	case AgentCursor:
		results = append(results, writeFile(".cursorrules", cursorrules, force))
	case AgentAll:
		results = append(results, writeFile("llms.txt", llmsTxt, force))
		results = append(results, writeFile("CLAUDE.md", claudeMd, force))
		results = append(results, writeFile(".cursorrules", cursorrules, force))
	}

	return results, nil
}

func writeFile(filename, content string, force bool) InitResult {
	result := InitResult{File: filename}

	path := filepath.Join(".", filename)

	if !force {
		if _, err := os.Stat(path); err == nil {
			result.Skipped = true
			return result
		}
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		result.Error = err
		return result
	}

	result.Created = true
	return result
}
