package solo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Deployer handles deployment to various platforms
type Deployer struct {
	projectPath string
	config      *ProjectConfig
}

// NewDeployer creates a new deployer instance
func NewDeployer(projectPath string, config *ProjectConfig) *Deployer {
	return &Deployer{
		projectPath: projectPath,
		config:      config,
	}
}

// Deploy handles the deployment process based on platform
func (d *Deployer) Deploy() error {
	switch d.config.Platform {
	case PlatformGitHubPages:
		return d.deployToGitHubPages()
	case PlatformVercel:
		return d.deployToVercel()
	case PlatformNetlify:
		return d.deployToNetlify()
	case PlatformRailway:
		return d.deployToRailway()
	default:
		return fmt.Errorf("unsupported deployment platform: %s", d.config.Platform)
	}
}

// deployToGitHubPages deploys to GitHub Pages
func (d *Deployer) deployToGitHubPages() error {
	repoName := d.config.Name
	if !strings.HasSuffix(repoName, ".github.io") {
		repoName = repoName + ".github.io"
	}

	// Initialize git repository
	if err := d.runCommand("git", "init"); err != nil {
		return fmt.Errorf("git init failed: %w", err)
	}

	// Create .gitignore
	gitignore := `# Dependencies
node_modules/
.env
.env.local
.env.production

# Build outputs
build/
dist/

# IDE
.vscode/
.idea/

# OS
.DS_Store
Thumbs.db`

	if err := os.WriteFile(filepath.Join(d.projectPath, ".gitignore"), []byte(gitignore), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	// Create GitHub Actions workflow for automatic deployment
	workflowDir := filepath.Join(d.projectPath, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return fmt.Errorf("failed to create workflows directory: %w", err)
	}

	workflow := `name: Deploy to GitHub Pages

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'
        cache: 'npm'
    
    - name: Install dependencies
      run: npm ci
    
    - name: Build
      run: npm run build
    
    - name: Deploy to GitHub Pages
      uses: peaceiris/actions-gh-pages@v3
      with:
        github_token: ${{ secrets.GITHUB_TOKEN }}
        publish_dir: ./build`

	// For static sites, use current directory
	if d.config.Type == ProjectTypeStatic {
		workflow = strings.Replace(workflow, "npm ci", "# Static site - no build needed", 1)
		workflow = strings.Replace(workflow, "npm run build", "# Static site - no build needed", 1)
		workflow = strings.Replace(workflow, "./build", ".", 1)
	}

	if err := os.WriteFile(filepath.Join(workflowDir, "deploy.yml"), []byte(workflow), 0644); err != nil {
		return fmt.Errorf("failed to create workflow: %w", err)
	}

	// Add and commit files
	commands := [][]string{
		{"git", "add", "."},
		{"git", "commit", "-m", "Initial commit - SOLO mode generated"},
		{"git", "branch", "-M", "main"},
	}

	for _, cmd := range commands {
		if err := d.runCommand(cmd[0], cmd[1:]...); err != nil {
			return fmt.Errorf("git command failed: %w", err)
		}
	}

	fmt.Printf("✅ GitHub Pages项目已准备就绪！\n")
	fmt.Printf("📁 项目路径: %s\n", d.projectPath)
	fmt.Printf("🚀 下一步: 将代码推送到GitHub仓库\n")
	fmt.Printf("   1. 在GitHub上创建仓库: %s\n", repoName)
	fmt.Printf("   2. 运行: git remote add origin <your-repo-url>\n")
	fmt.Printf("   3. 运行: git push -u origin main\n")

	return nil
}

// deployToVercel deploys to Vercel
func (d *Deployer) deployToVercel() error {
	vercelJSON := `{
  "version": 2,
  "builds": [
    {
      "src": "package.json",
      "use": "@vercel/static-build",
      "config": {
        "distDir": "build"
      }
    }
  ],
  "routes": [
    {
      "src": "/(.*)",
      "dest": "/$1"
    }
  ]
}`

	// Adjust for different project types
	if d.config.Type == ProjectTypeReact || d.config.Type == ProjectTypeNextJS {
		vercelJSON = strings.Replace(vercelJSON, "build", "build", 1)
	} else if d.config.Type == ProjectTypeExpress {
		vercelJSON = `{
  "version": 2,
  "builds": [
    {
      "src": "src/index.js",
      "use": "@vercel/node"
    }
  ],
  "routes": [
    {
      "src": "/(.*)",
      "dest": "src/index.js"
    }
  ]
}`
	}

	if err := os.WriteFile(filepath.Join(d.projectPath, "vercel.json"), []byte(vercelJSON), 0644); err != nil {
		return fmt.Errorf("failed to create vercel.json: %w", err)
	}

	// Create build script if not exists
	if d.config.Type == ProjectTypeStatic {
		packagePath := filepath.Join(d.projectPath, "package.json")
		if _, err := os.Stat(packagePath); os.IsNotExist(err) {
			packageJSON := `{
  "name": "` + d.config.Name + `",
  "version": "1.0.0",
  "scripts": {
    "build": "echo 'Static site - no build needed'"
  }
}`
			if err := os.WriteFile(packagePath, []byte(packageJSON), 0644); err != nil {
				return fmt.Errorf("failed to create package.json: %w", err)
			}
		}
	}

	fmt.Printf("✅ Vercel项目已准备就绪！\n")
	fmt.Printf("📁 项目路径: %s\n", d.projectPath)
	fmt.Printf("🚀 部署命令: vercel --prod\n")

	return nil
}

// deployToNetlify deploys to Netlify
func (d *Deployer) deployToNetlify() error {
	netlifyTOML := `[build]
  publish = "build"
  command = "npm run build"

[build.environment]
  NODE_VERSION = "18"

[[redirects]]
  from = "/*"
  to = "/index.html"
  status = 200`

	// Adjust for different project types
	if d.config.Type == ProjectTypeStatic {
		netlifyTOML = strings.Replace(netlifyTOML, "build", ".", 1)
		netlifyTOML = strings.Replace(netlifyTOML, "npm run build", "echo 'Static site - no build needed'", 1)
	} else if d.config.Type == ProjectTypeFullStack {
		netlifyTOML = `[build]
  command = "npm run build"
  publish = "client/build"

[dev]
  command = "npm run dev"
  port = 3000
  publish = "client/build"

[[redirects]]
  from = "/api/*"
  to = "/.netlify/functions/server/:splat"
  status = 200

[[redirects]]
  from = "/*"
  to = "/index.html"
  status = 200`
	}

	if err := os.WriteFile(filepath.Join(d.projectPath, "netlify.toml"), []byte(netlifyTOML), 0644); err != nil {
		return fmt.Errorf("failed to create netlify.toml: %w", err)
	}

	fmt.Printf("✅ Netlify项目已准备就绪！\n")
	fmt.Printf("📁 项目路径: %s\n", d.projectPath)
	fmt.Printf("🚀 部署方式:\n")
	fmt.Printf("   1. 推送到GitHub\n")
	fmt.Printf("   2. 在Netlify上导入GitHub仓库\n")

	return nil
}

// deployToRailway deploys to Railway
func (d *Deployer) deployToRailway() error {
	railwayTOML := `[build]
builder = "NIXPACKS"

[deploy]
startCommand = "npm start"
restartPolicyType = "ON_FAILURE"
restartPolicyMaxRetries = 3`

	if err := os.WriteFile(filepath.Join(d.projectPath, "railway.toml"), []byte(railwayTOML), 0644); err != nil {
		return fmt.Errorf("failed to create railway.toml: %w", err)
	}

	// Create Procfile for Railway
	procfile := `web: npm start`
	if err := os.WriteFile(filepath.Join(d.projectPath, "Procfile"), []byte(procfile), 0644); err != nil {
		return fmt.Errorf("failed to create Procfile: %w", err)
	}

	// Create environment variables template
	envVars := `# Railway环境变量
PORT=3000
NODE_ENV=production

# 数据库配置
DATABASE_URL=your-database-url

# 其他配置
API_URL=https://your-app.railway.app`

	if err := os.WriteFile(filepath.Join(d.projectPath, ".env.railway"), []byte(envVars), 0644); err != nil {
		return fmt.Errorf("failed to create .env.railway: %w", err)
	}

	fmt.Printf("✅ Railway项目已准备就绪！\n")
	fmt.Printf("📁 项目路径: %s\n", d.projectPath)
	fmt.Printf("🚀 部署命令: railway up\n")

	return nil
}

// runCommand executes a command in the project directory
func (d *Deployer) runCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = d.projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// GetDeploymentURL returns the expected deployment URL
func (d *Deployer) GetDeploymentURL() string {
	repoName := d.config.Name
	if !strings.HasSuffix(repoName, ".github.io") {
		repoName = repoName + ".github.io"
	}

	switch d.config.Platform {
	case PlatformGitHubPages:
		return fmt.Sprintf("https://%s.github.io/%s", "your-username", d.config.Name)
	case PlatformVercel:
		return fmt.Sprintf("https://%s.vercel.app", d.config.Name)
	case PlatformNetlify:
		return fmt.Sprintf("https://%s.netlify.app", d.config.Name)
	case PlatformRailway:
		return fmt.Sprintf("https://%s.up.railway.app", d.config.Name)
	default:
		return ""
	}
}

// GetNextSteps returns the next steps for deployment
func (d *Deployer) GetNextSteps() []string {
	switch d.config.Platform {
	case PlatformGitHubPages:
		return []string{
			"在GitHub上创建新仓库",
			"git remote add origin <your-repo-url>",
			"git push -u origin main",
			"在仓库设置中启用GitHub Pages",
		}
	case PlatformVercel:
		return []string{
			"安装Vercel CLI: npm i -g vercel",
			"运行: vercel --prod",
		}
	case PlatformNetlify:
		return []string{
			"将代码推送到GitHub",
			"在Netlify.com导入GitHub仓库",
			"自动部署将开始",
		}
	case PlatformRailway:
		return []string{
			"安装Railway CLI: npm i -g @railway/cli",
			"运行: railway login",
			"运行: railway up",
		}
	default:
		return []string{"请联系支持获取部署帮助"}
	}
}