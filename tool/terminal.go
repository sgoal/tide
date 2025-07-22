package tool

import (
	"encoding/json"
	"os/exec"
)

// TerminalTool is a tool for executing terminal commands.
type TerminalTool struct{}

// TerminalToolArgs represents the arguments for the TerminalTool.
type TerminalToolArgs struct {
	Command string `json:"command"`
}

func (t *TerminalTool) Name() string {
	return "terminal"
}

// Description returns the description of the tool.
func (t *TerminalTool) Description() string {
	return "A tool for executing terminal commands."
}

// Execute executes a terminal command and returns its output.
func (t *TerminalTool) Execute(args json.RawMessage) (string, error) {
	var toolArgs TerminalToolArgs
	if err := json.Unmarshal(args, &toolArgs); err != nil {
		return "", err
	}

	cmd := exec.Command("sh", "-c", toolArgs.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}