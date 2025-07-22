package tool

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileEditorTool is a tool for editing files.
type FileEditorTool struct{}

func (t *FileEditorTool) Name() string {
	return "file_editor"
}

func (t *FileEditorTool) Description() string {
	return "A tool for reading files from a directory, modifying their content, and writing them back. Input should be a JSON object with 'dir_path', 'file_name', 'search_text', and 'replace_text'."
}

func (t *FileEditorTool) Execute(args json.RawMessage) (string, error) {
	var params struct {
		DirPath     string `json:"dir_path"`
		FileName    string `json:"file_name"`
		SearchText  string `json:"search_text"`
		ReplaceText string `json:"replace_text"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("invalid arguments for file_editor tool: %w", err)
	}

	filePath := filepath.Join(params.DirPath, params.FileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) && params.SearchText == "" {
			if err1 := os.WriteFile(filePath, []byte(params.ReplaceText), 0644); err1 != nil {
				return "", fmt.Errorf("failed to create file: %w", err)
			}
			return fmt.Sprintf("Successfully created file: %s", filePath), nil
		}
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	newContent := strings.ReplaceAll(string(content), params.SearchText, params.ReplaceText)
	if string(content) == newContent {
		return fmt.Sprintf("Search text not found in %s. File not modified.", filePath), nil
	}
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fmt.Sprintf("Successfully modified file: %s", filePath), nil
}
