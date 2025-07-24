package tool

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// DeployerTool is a tool for deploying projects.
type DeployerTool struct{}

func (t *DeployerTool) Name() string {
	return "deployer"
}

func (t *DeployerTool) Description() string {
	return "A tool for deploying projects using Vercel."
}

// Execute executes the deployer tool.
func (t *DeployerTool) Execute(args json.RawMessage) (string, error) {
	var params struct {
		ProjectPath string `json:"project_path"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("invalid arguments: %w", err)
	}

	if params.ProjectPath == "" {
		return "", fmt.Errorf("usage: deployer '{\"project_path\": \"/path/to/project\"}'")
	}

	cmd := exec.Command("vercel", "--prod")
	cmd.Dir = params.ProjectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to deploy project: %v\n%s", err, string(output))
	}

	return fmt.Sprintf("Project deployed successfully:\n%s", string(output)), nil
}