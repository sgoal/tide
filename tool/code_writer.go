package tool

import (
	"encoding/json"
	"fmt"
	"os"
)

// CodeWriterTool is a tool for writing code to a file.
type CodeWriterTool struct{}

func (t *CodeWriterTool) Name() string {
	return "code_writer"
}

func (t *CodeWriterTool) Description() string {
	return "A tool for writing code to a file. The input should be a JSON object with 'filepath' and 'code' keys."
}

// Execute expects args to be a JSON string with "filepath" and "code"
func (t *CodeWriterTool) Execute(args json.RawMessage) (string, error) {
	var params struct {
		Filepath string `json:"filepath"`
		Code     string `json:"code"`
	}
	err := json.Unmarshal(args, &params)
	if err != nil {
		return "", fmt.Errorf("invalid arguments for code_writer tool: %w", err)
	}

	err = os.WriteFile(params.Filepath, []byte(params.Code), 0644)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Successfully wrote code to %s", params.Filepath), nil
}