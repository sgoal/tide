package tool

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// CodeWriterTool is a tool for writing code to a file.
type CodeWriterTool struct{}

func (t *CodeWriterTool) Name() string {
	return "code_writer"
}

func (t *CodeWriterTool) Description() string {
	return "A tool for writing code to a file. The input should be a JSON object with 'dir_path', 'file_name', and 'code' keys."
}

// Execute expects args to be a JSON string with "filepath" and "code"
func (t *CodeWriterTool) Execute(args json.RawMessage) (string, error) {
	var params struct {
		DirPath  string `json:"dir_path"`
		FileName string `json:"file_name"`
		Code     string `json:"code"`
	}
	err := json.Unmarshal(args, &params)
	if err != nil {
		return "", fmt.Errorf("invalid arguments for code_writer tool: %w", err)
	}

	filePath := filepath.Join(params.DirPath, params.FileName)
	err = os.WriteFile(filePath, []byte(params.Code), 0644)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Successfully wrote code to %s", filePath), nil
}
