package solo

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sgoal/tide/tool"
)

// ProjectType represents different types of projects
const (
	ProjectTypeStatic   = "static"
	ProjectTypeReact    = "react"
	ProjectTypeNextJS   = "nextjs"
	ProjectTypeExpress  = "express"
	ProjectTypeFullStack = "fullstack"
)

// DeployPlatform represents deployment platforms
const (
	PlatformGitHubPages = "github-pages"
	PlatformVercel      = "vercel"
	PlatformNetlify     = "netlify"
	PlatformRailway     = "railway"
)

// SoloManager manages the SOLO mode functionality
type SoloManager struct {
	tools     map[string]tool.Tool
	logWriter io.Writer
	workspace string
}

// ProjectConfig holds project configuration
type ProjectConfig struct {
	Name        string            `json:"name"`
	Type        string            `json:"type"`
	Platform    string            `json:"platform"`
	Database    string            `json:"database,omitempty"`
	Description string            `json:"description"`
	EnvVars     map[string]string `json:"env_vars,omitempty"`
}

// NewSoloManager creates a new SOLO mode manager
func NewSoloManager(logWriter io.Writer) *SoloManager {
	if logWriter == nil {
		logWriter = os.Stdout
	}

	return &SoloManager{
		tools: map[string]tool.Tool{
			"search":      &tool.SearchTool{},
			"code_writer": &tool.CodeWriterTool{},
			"file_editor": &tool.FileEditorTool{},
			"terminal":    &tool.TerminalTool{},
		},
		logWriter: logWriter,
		workspace: getWorkspacePath(),
	}
}

// StartSoloMode starts the SOLO mode with a given requirement
func (sm *SoloManager) StartSoloMode(requirement string) error {
	fmt.Fprintf(sm.logWriter, "🚀 启动SOLO模式: %s\n", requirement)

	// Step 1: Parse requirement and detect project type
	config, err := sm.parseRequirement(requirement)
	if err != nil {
		return fmt.Errorf("解析需求失败: %w", err)
	}

	// Step 2: Create project structure
	projectPath := filepath.Join(sm.workspace, config.Name)
	if err := sm.createProjectStructure(projectPath, config); err != nil {
		return fmt.Errorf("创建项目结构失败: %w", err)
	}

	// Step 3: Generate project files
	if err := sm.generateProjectFiles(projectPath, config); err != nil {
		return fmt.Errorf("生成项目文件失败: %w", err)
	}

	// Step 4: Install dependencies
	if err := sm.installDependencies(projectPath, config); err != nil {
		return fmt.Errorf("安装依赖失败: %w", err)
	}

	// Step 5: Deploy to platform
	if err := sm.deployProject(projectPath, config); err != nil {
		return fmt.Errorf("部署项目失败: %w", err)
	}

	fmt.Fprintf(sm.logWriter, "✅ SOLO模式完成! 项目已部署到: %s\n", config.Platform)
	return nil
}

// parseRequirement analyzes the requirement and creates project configuration
func (sm *SoloManager) parseRequirement(requirement string) (*ProjectConfig, error) {
	config := &ProjectConfig{
		Name:        sm.generateProjectName(requirement),
		EnvVars:     make(map[string]string),
	}

	// Simple keyword-based project type detection
	requirementLower := strings.ToLower(requirement)
	
	if strings.Contains(requirementLower, "react") {
		config.Type = ProjectTypeReact
		config.Platform = PlatformVercel
	} else if strings.Contains(requirementLower, "nextjs") || strings.Contains(requirementLower, "next.js") {
		config.Type = ProjectTypeNextJS
		config.Platform = PlatformVercel
	} else if strings.Contains(requirementLower, "express") || strings.Contains(requirementLower, "api") {
		config.Type = ProjectTypeExpress
		config.Platform = PlatformRailway
	} else if strings.Contains(requirementLower, "blog") || strings.Contains(requirementLower, "static") {
		config.Type = ProjectTypeStatic
		config.Platform = PlatformGitHubPages
	} else {
		// Default to fullstack for complex requirements
		config.Type = ProjectTypeFullStack
		config.Platform = PlatformNetlify
	}

	// Database detection
	if strings.Contains(requirementLower, "database") || strings.Contains(requirementLower, "db") {
		config.Database = "supabase"
	}

	config.Description = requirement
	return config, nil
}

// generateProjectName creates a project name from requirement
func (sm *SoloManager) generateProjectName(requirement string) string {
	// Simple name generation - take first few words and sanitize
	words := strings.Fields(requirement)
	if len(words) == 0 {
		return "solo-project"
	}

	name := strings.ToLower(words[0])
	for _, word := range words[1:] {
		if len(name) > 20 {
			break
		}
		name += "-" + strings.ToLower(word)
	}

	// Sanitize name
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return '-'
	}, name)

	return strings.Trim(name, "-")
}

// createProjectStructure creates the basic project directory structure
func (sm *SoloManager) createProjectStructure(projectPath string, config *ProjectConfig) error {
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return err
	}

	// Create subdirectories based on project type
	dirs := sm.getProjectDirectories(config.Type)
	for _, dir := range dirs {
		fullPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return err
		}
	}

	return nil
}

// getProjectDirectories returns the directory structure for a project type
func (sm *SoloManager) getProjectDirectories(projectType string) []string {
	switch projectType {
	case ProjectTypeStatic:
		return []string{"css", "js", "images"}
	case ProjectTypeReact, ProjectTypeNextJS:
		return []string{"src", "public", "src/components", "src/styles"}
	case ProjectTypeExpress:
		return []string{"src", "src/routes", "src/models", "src/middleware"}
	case ProjectTypeFullStack:
		return []string{"client", "server", "client/src", "server/src", "server/routes"}
	default:
		return []string{"src"}
	}
}

// Helper function to get workspace path
func getWorkspacePath() string {
	workspace := os.Getenv("TIDE_WORKSPACE")
	if workspace == "" {
		workspace = filepath.Join(os.Getenv("HOME"), "tide-workspace")
	}
	return workspace
}

// generateProjectFiles generates all project files based on configuration
func (sm *SoloManager) generateProjectFiles(projectPath string, config *ProjectConfig) error {
	fmt.Fprintf(sm.logWriter, "📁 生成项目文件...\n")
	
	generator := NewProjectGenerator()
	if err := generator.GenerateFiles(projectPath, config); err != nil {
		return fmt.Errorf("生成项目文件失败: %w", err)
	}

	// Save configuration
	if err := sm.SaveConfig(projectPath, config); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	// Generate platform-specific configuration files
	if err := sm.generatePlatformConfig(projectPath, config); err != nil {
		return fmt.Errorf("生成平台配置失败: %w", err)
	}

	return nil
}

// installDependencies installs project dependencies
func (sm *SoloManager) installDependencies(projectPath string, config *ProjectConfig) error {
	fmt.Fprintf(sm.logWriter, "📦 安装依赖...\n")
	
	var cmd *exec.Cmd
	
	switch config.Type {
	case ProjectTypeReact, ProjectTypeNextJS, ProjectTypeExpress, ProjectTypeFullStack:
		// Check if package.json exists
		packagePath := filepath.Join(projectPath, "package.json")
		if _, err := os.Stat(packagePath); err == nil {
			cmd = exec.Command("npm", "install")
			cmd.Dir = projectPath
			cmd.Stdout = sm.logWriter
			cmd.Stderr = sm.logWriter
			
			if err := cmd.Run(); err != nil {
				// Try with yarn as fallback
				cmd = exec.Command("yarn", "install")
				cmd.Dir = projectPath
				cmd.Stdout = sm.logWriter
				cmd.Stderr = sm.logWriter
				
				if err := cmd.Run(); err != nil {
					fmt.Fprintf(sm.logWriter, "⚠️  依赖安装失败，但项目已创建完成\n")
					return nil // Continue with deployment
				}
			}
			fmt.Fprintf(sm.logWriter, "✅ 依赖安装完成\n")
		}
	case ProjectTypeStatic:
		// No dependencies for static sites
		fmt.Fprintf(sm.logWriter, "✅ 静态站点无需安装依赖\n")
	}
	
	return nil
}

// deployProject handles the deployment process
func (sm *SoloManager) deployProject(projectPath string, config *ProjectConfig) error {
	fmt.Fprintf(sm.logWriter, "🚀 开始部署到 %s...\n", config.Platform)
	
	deployer := NewDeployer(projectPath, config)
	
	// Generate deployment configuration
	if err := deployer.Deploy(); err != nil {
		return fmt.Errorf("部署失败: %w", err)
	}

	// Print deployment information
	fmt.Fprintf(sm.logWriter, "\n🎉 部署准备完成！\n")
	fmt.Fprintf(sm.logWriter, "📋 项目信息:\n")
	fmt.Fprintf(sm.logWriter, "   名称: %s\n", config.Name)
	fmt.Fprintf(sm.logWriter, "   类型: %s\n", config.Type)
	fmt.Fprintf(sm.logWriter, "   平台: %s\n", config.Platform)
	fmt.Fprintf(sm.logWriter, "   路径: %s\n", projectPath)
	
	// Print next steps
	fmt.Fprintf(sm.logWriter, "\n📋 下一步操作:\n")
	steps := deployer.GetNextSteps()
	for i, step := range steps {
		fmt.Fprintf(sm.logWriter, "   %d. %s\n", i+1, step)
	}
	
	// Print expected URL
	url := deployer.GetDeploymentURL()
	if url != "" {
		fmt.Fprintf(sm.logWriter, "\n🌐 预期访问地址: %s\n", url)
	}
	
	return nil
}

// generatePlatformConfig generates platform-specific configuration files
func (sm *SoloManager) generatePlatformConfig(projectPath string, config *ProjectConfig) error {
	// Create README.md for all platforms
	readme := fmt.Sprintf(`# %s

由Tide SOLO模式生成的项目。

## 项目信息
- **名称**: %s
- **类型**: %s
- **平台**: %s
- **描述**: %s

## 快速开始

### 开发环境
`+"```"+`bash
cd %s
`, config.Name, config.Name, config.Type, config.Platform, config.Description, config.Name)

	switch config.Type {
	case ProjectTypeReact, ProjectTypeNextJS:
		readme += `npm install
npm start
`+"```"+`

### 部署
`+"```"+`bash
npm run build
vercel --prod
`+"```"+`
`
	case ProjectTypeExpress:
		readme += `npm install
npm run dev
`+"```"+`

### 部署
`+"```"+`bash
railway up
`+"```"+`
`
	case ProjectTypeStatic:
		readme += `# 静态站点无需构建

### 部署到GitHub Pages
1. 推送到GitHub仓库
2. 在仓库设置中启用GitHub Pages
`
	case ProjectTypeFullStack:
		readme += `npm install
npm run dev
`+"```"+`

### 部署到Netlify
1. 推送到GitHub仓库
2. 在Netlify.com导入仓库
`
	}

	readme += fmt.Sprintf(`

## 项目结构
`+"```"+`
%s/
├── solo-config.json    # SOLO模式配置
├── README.md          # 项目文档
`, config.Name)

	return os.WriteFile(filepath.Join(projectPath, "README.md"), []byte(readme), 0644)
}

// SaveConfig saves project configuration to file
func (sm *SoloManager) SaveConfig(projectPath string, config *ProjectConfig) error {
	configFile := filepath.Join(projectPath, "solo-config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0644)
}

// LoadConfig loads project configuration from file
func (sm *SoloManager) LoadConfig(projectPath string) (*ProjectConfig, error) {
	configFile := filepath.Join(projectPath, "solo-config.json")
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config ProjectConfig
	err = json.Unmarshal(data, &config)
	return &config, err
}