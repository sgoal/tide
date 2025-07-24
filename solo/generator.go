package solo

import (
	"fmt"
	"os"
	"path/filepath"
)

// ProjectGenerator handles file generation for different project types
type ProjectGenerator struct{}

// NewProjectGenerator creates a new project generator
func NewProjectGenerator() *ProjectGenerator {
	return &ProjectGenerator{}
}

// GenerateFiles generates all necessary files for the project
func (pg *ProjectGenerator) GenerateFiles(projectPath string, config *ProjectConfig) error {
	switch config.Type {
	case ProjectTypeStatic:
		return pg.generateStaticFiles(projectPath, config)
	case ProjectTypeReact:
		return pg.generateReactFiles(projectPath, config)
	case ProjectTypeNextJS:
		return pg.generateNextJSFiles(projectPath, config)
	case ProjectTypeExpress:
		return pg.generateExpressFiles(projectPath, config)
	case ProjectTypeFullStack:
		return pg.generateFullStackFiles(projectPath, config)
	default:
		return fmt.Errorf("unsupported project type: %s", config.Type)
	}
}

func (pg *ProjectGenerator) generateStaticFiles(projectPath string, config *ProjectConfig) error {
	// Generate index.html
	indexHTML := fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <link rel="stylesheet" href="css/style.css">
</head>
<body>
    <header>
        <h1>%s</h1>
        <p>%s</p>
    </header>
    
    <main>
        <section class="hero">
            <h2>欢迎使用SOLO模式生成的项目</h2>
            <p>这是一个由AI驱动的静态网站项目。</p>
        </section>
        
        <section class="features">
            <h3>项目特性</h3>
            <ul>
                <li>响应式设计</li>
                <li>现代CSS</li>
                <li>优化的性能</li>
                <li>一键部署</li>
            </ul>
        </section>
    </main>
    
    <footer>
        <p>&copy; 2024 由Tide SOLO模式生成</p>
    </footer>
    
    <script src="js/main.js"></script>
</body>
</html>`, config.Name, config.Name, config.Description)

	// Generate CSS
	cssContent := `/* 现代CSS重置 */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
    line-height: 1.6;
    color: #333;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    min-height: 100vh;
}

header {
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(10px);
    padding: 2rem;
    text-align: center;
    color: white;
}

header h1 {
    font-size: 2.5rem;
    margin-bottom: 0.5rem;
}

header p {
    font-size: 1.2rem;
    opacity: 0.8;
}

main {
    max-width: 800px;
    margin: 2rem auto;
    padding: 0 1rem;
}

.hero {
    background: white;
    padding: 2rem;
    border-radius: 10px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
    margin-bottom: 2rem;
    text-align: center;
}

.features {
    background: white;
    padding: 2rem;
    border-radius: 10px;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.features ul {
    list-style: none;
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
    margin-top: 1rem;
}

.features li {
    padding: 1rem;
    background: #f8f9fa;
    border-radius: 5px;
    text-align: center;
}

footer {
    text-align: center;
    padding: 2rem;
    color: white;
    opacity: 0.7;
}

@media (max-width: 600px) {
    header h1 {
        font-size: 2rem;
    }
    
    .hero, .features {
        padding: 1.5rem;
    }
}`

	// Generate JavaScript
	jsContent := `// 简单的交互功能
console.log('SOLO模式生成的静态网站已加载');

// 添加一些交互效果
document.addEventListener('DOMContentLoaded', function() {
    const features = document.querySelectorAll('.features li');
    
    features.forEach(feature => {
        feature.addEventListener('click', function() {
            this.style.transform = 'scale(1.05)';
            setTimeout(() => {
                this.style.transform = 'scale(1)';
            }, 200);
        });
    });
});

// 添加滚动动画
window.addEventListener('scroll', function() {
    const elements = document.querySelectorAll('.hero, .features');
    
    elements.forEach(element => {
        const elementTop = element.getBoundingClientRect().top;
        const elementVisible = 150;
        
        if (elementTop < window.innerHeight - elementVisible) {
            element.classList.add('visible');
        }
    });
});`

	// Write files
	files := map[string]string{
		"index.html":    indexHTML,
		"css/style.css": cssContent,
		"js/main.js":    jsContent,
	}

	return pg.writeFiles(projectPath, files)
}

func (pg *ProjectGenerator) generateReactFiles(projectPath string, config *ProjectConfig) error {
	packageJSON := fmt.Sprintf(`{
  "name": "%s",
  "version": "1.0.0",
  "private": true,
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-scripts": "5.0.1"
  },
  "scripts": {
    "start": "react-scripts start",
    "build": "react-scripts build",
    "test": "react-scripts test",
    "eject": "react-scripts eject"
  },
  "eslintConfig": {
    "extends": [
      "react-app",
      "react-app/jest"
    ]
  },
  "browserslist": {
    "production": [
      ">0.2%%",
      "not dead",
      "not op_mini all"
    ],
    "development": [
      "last 1 chrome version",
      "last 1 firefox version",
      "last 1 safari version"
    ]
  }
}`, config.Name)

	appJS := `import React from 'react';
import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>SOLO模式生成的React应用</h1>
        <p>
          这是一个由AI驱动的React项目，已准备好部署到Vercel。
        </p>
        <div className="features">
          <h2>项目特性</h2>
          <ul>
            <li>⚡ 现代React 18</li>
            <li>🎨 CSS Modules</li>
            <li>📱 响应式设计</li>
            <li>🚀 一键部署</li>
          </ul>
        </div>
      </header>
    </div>
  );
}

export default App;`

	appCSS := `.App {
  text-align: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.App-header {
  padding: 2rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.App-header h1 {
  font-size: 2.5rem;
  margin-bottom: 1rem;
}

.App-header p {
  font-size: 1.2rem;
  margin-bottom: 2rem;
  opacity: 0.9;
}

.features {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  padding: 2rem;
  border-radius: 10px;
  margin: 2rem;
}

.features ul {
  list-style: none;
  padding: 0;
}

.features li {
  padding: 0.5rem;
  font-size: 1.1rem;
}`

	publicFiles := map[string]string{
		"index.html": `<!DOCTYPE html>
<html lang="zh-CN">
  <head>
    <meta charset="utf-8" />
    <link rel="icon" href="%PUBLIC_URL%/favicon.ico" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="theme-color" content="#000000" />
    <meta name="description" content="由SOLO模式生成的React应用" />
    <title>SOLO React App</title>
  </head>
  <body>
    <noscript>You need to enable JavaScript to run this app.</noscript>
    <div id="root"></div>
  </body>
</html>`,
		"manifest.json": `{
  "short_name": "SOLO App",
  "name": "SOLO模式生成的React应用",
  "icons": [],
  "start_url": ".",
  "display": "standalone",
  "theme_color": "#000000",
  "background_color": "#ffffff"
}`,
	}

	files := map[string]string{
		"package.json": packageJSON,
		"src/App.js":   appJS,
		"src/App.css":  appCSS,
		"src/index.js": `import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);`,
		"src/index.css": `body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

code {
  font-family: source-code-pro, Menlo, Monaco, Consolas, 'Courier New',
    monospace;
}`,
	}

	// Add public files
	for name, content := range publicFiles {
		files[filepath.Join("public", name)] = content
	}

	return pg.writeFiles(projectPath, files)
}

func (pg *ProjectGenerator) generateNextJSFiles(projectPath string, config *ProjectConfig) error {
	packageJSON := fmt.Sprintf(`{
  "name": "%s",
  "version": "0.1.0",
  "private": true,
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start"
  },
  "dependencies": {
    "react": "18.2.0",
    "react-dom": "18.2.0",
    "next": "13.4.12"
  }
}`, config.Name)

	indexJS := `import Head from 'next/head'

export default function Home() {
  return (
    <div>
      <Head>
        <title>Create Next App</title>
        <meta name="description" content="Generated by create next app" />
        <link rel="icon" href="/favicon.ico" />
      </Head>

      <main>
        <h1>
          Welcome to <a href="https://nextjs.org">Next.js!</a>
        </h1>
      </main>
    </div>
  )
}`

	files := map[string]string{
		"package.json":   packageJSON,
		"pages/index.js": indexJS,
	}

	return pg.writeFiles(projectPath, files)
}

func (pg *ProjectGenerator) generateExpressFiles(projectPath string, config *ProjectConfig) error {
	packageJSON := fmt.Sprintf(`{
  "name": "%s",
  "version": "1.0.0",
  "description": "%s",
  "main": "src/index.js",
  "scripts": {
    "start": "node src/index.js",
    "dev": "nodemon src/index.js",
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "dependencies": {
    "express": "^4.18.2",
    "cors": "^2.8.5",
    "helmet": "^7.0.0",
    "morgan": "^1.10.0"%s
  },
  "devDependencies": {
    "nodemon": "^3.0.1"
  },
  "keywords": ["express", "api", "solo"],
  "author": "SOLO模式",
  "license": "MIT"
}`, config.Name, config.Description, pg.getDatabaseDeps(config.Database))

	indexJS := fmt.Sprintf(`const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const morgan = require('morgan');
%s

const app = express();
const PORT = process.env.PORT || 3000;

// 中间件
app.use(helmet());
app.use(cors());
app.use(morgan('combined'));
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// 健康检查端点
app.get('/health', (req, res) => {
  res.json({ status: 'OK', message: 'SOLO模式生成的Express API运行正常' });
});

// API路由
app.get('/api', (req, res) => {
  res.json({
    message: '欢迎使用SOLO模式生成的Express API',
    version: '1.0.0',
    endpoints: [
      'GET /health - 健康检查',
      'GET /api - API信息',
      'GET /api/users - 用户列表'
    ]
  });
});

// 示例用户路由
const users = [
  { id: 1, name: '用户1', email: 'user1@example.com' },
  { id: 2, name: '用户2', email: 'user2@example.com' }
];

app.get('/api/users', (req, res) => {
  res.json(users);
});

// 错误处理中间件
app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).json({ error: '服务器内部错误' });
});

// 404处理
app.use((req, res) => {
  res.status(404).json({ error: '路由未找到' });
});

app.listen(PORT, () => {
  console.log('服务器运行在端口 ' + PORT);
  console.log('访问 http://localhost:' + PORT + '/health 检查服务状态');
});

` + pg.getDatabaseSetup(config.Database))

	files := map[string]string{
		"package.json": packageJSON,
		"src/index.js": indexJS,
		".env.example": `PORT=3000
# 数据库配置
DATABASE_URL=your-database-url
# 其他环境变量
NODE_ENV=development`,
		"README.md": fmt.Sprintf(`# %s

由SOLO模式生成的Express.js API项目。

## 快速开始

### 安装依赖
`+"```"+`bash
npm install
`+"```"+`

### 开发模式
`+"```"+`bash
npm run dev
`+"```"+`

### 生产模式
`+"```"+`bash
npm start
`+"```"+`

## API端点

- GET /health - 健康检查
- GET /api - API信息
- GET /api/users - 用户列表

## 部署

本项目已配置为可部署到Railway平台。`, config.Name),
	}

	return pg.writeFiles(projectPath, files)
}

func (pg *ProjectGenerator) generateFullStackFiles(projectPath string, config *ProjectConfig) error {
	// 生成客户端文件（类似React）
	clientPackageJSON := `{
  "name": "solo-client",
  "version": "1.0.0",
  "private": true,
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-scripts": "5.0.1",
    "axios": "^1.4.0"
  },
  "scripts": {
    "start": "react-scripts start",
    "build": "react-scripts build",
    "test": "react-scripts test"
  },
  "proxy": "http://localhost:5000"
}`

	// 生成服务端文件
	serverPackageJSON := fmt.Sprintf(`{
  "name": "solo-server",
  "version": "1.0.0",
  "description": "%s",
  "main": "index.js",
  "scripts": {
    "start": "node index.js",
    "dev": "nodemon index.js"
  },
  "dependencies": {
    "express": "^4.18.2",
    "cors": "^2.8.5",
    "helmet": "^7.0.0",
    "morgan": "^1.10.0"%s
  },
  "devDependencies": {
    "nodemon": "^3.0.1"
  }
}`, config.Description, pg.getDatabaseDeps(config.Database))

	netlifyTOML := `[build]
  command = "npm run build"
  publish = "client/build"
  functions = "server/netlify/functions"

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

	files := map[string]string{
		"client/package.json": clientPackageJSON,
		"server/package.json": serverPackageJSON,
		"netlify.toml":        netlifyTOML,
		"package.json": fmt.Sprintf(`{
  "name": "%s",
  "version": "1.0.0",
  "description": "%s",
  "scripts": {
    "dev": "concurrently \"npm run dev --prefix server\" \"npm start --prefix client\"",
    "build": "npm run build --prefix client",
    "start": "npm start --prefix server"
  },
  "devDependencies": {
    "concurrently": "^8.2.0"
  }
}`, config.Name, config.Description),
	}

	return pg.writeFiles(projectPath, files)
}

// Helper functions
func (pg *ProjectGenerator) getDatabaseDeps(database string) string {
	switch database {
	case "supabase":
		return `,
    "@supabase/supabase-js": "^2.26.0"`
	case "mongodb":
		return `,
    "mongoose": "^7.4.0"`
	case "sqlite":
		return `,
    "sqlite3": "^5.1.6"`
	default:
		return ""
	}
}

func (pg *ProjectGenerator) getDatabaseSetup(database string) string {
	switch database {
	case "supabase":
		return `const { createClient } = require('@supabase/supabase-js');
const supabase = createClient(process.env.SUPABASE_URL, process.env.SUPABASE_KEY);`
	case "mongodb":
		return `const mongoose = require('mongoose');
mongoose.connect(process.env.MONGODB_URI || 'mongodb://localhost:27017/solo-app');`
	default:
		return ""
	}
}

func (pg *ProjectGenerator) writeFiles(projectPath string, files map[string]string) error {
	for relativePath, content := range files {
		fullPath := filepath.Join(projectPath, relativePath)

		// Create directory if it doesn't exist
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("创建目录失败 %s: %w", dir, err)
		}

		// Write file
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("写入文件失败 %s: %w", fullPath, err)
		}
	}
	return nil
}
