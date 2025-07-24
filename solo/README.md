# Tide SOLO模式

Tide SOLO模式是一个AI驱动的项目生成与部署工具，能够根据自然语言需求描述，自动生成完整的项目结构并配置部署到主流平台。

## ✨ 功能特性

- **🎯 智能需求解析**: 自动识别项目类型和技术栈
- **🚀 一键生成**: 15分钟内完成从需求到部署
- **🌐 多平台支持**: GitHub Pages、Vercel、Netlify、Railway
- **📊 数据库集成**: Supabase、MongoDB、SQLite支持
- **🎨 现代模板**: 响应式设计，最佳实践
- **⚡ 零配置**: 自动生成所有必要配置文件

## 🚀 快速开始

### 使用方法

```bash
# 进入项目目录
cd /Users/bytedance/go/src/github.com/sgoal/tide

# 运行SOLO模式
go run solo/main.go "你的项目需求描述"
```

### 使用示例

#### 1. 静态网站
```bash
go run solo/main.go "创建一个响应式的个人作品集网站"
```

#### 2. React应用
```bash
go run solo/main.go "构建一个React博客，需要文章列表和详情页"
```

#### 3. Express API
```bash
go run solo/main.go "创建一个RESTful API，支持用户注册登录"
```

#### 4. Next.js应用
```bash
go run solo/main.go "制作一个Next.js电商网站，需要商品展示和购物车"
```

#### 5. 全栈应用
```bash
go run solo/main.go "构建一个任务管理应用，React前端和Express后端"
```

## 📋 项目类型支持

| 项目类型 | 描述 | 部署平台 | 技术栈 |
|----------|------|----------|--------|
| **static** | 静态网站 | GitHub Pages | HTML/CSS/JS |
| **react** | React应用 | Vercel | React 18, CSS Modules |
| **nextjs** | Next.js应用 | Vercel | Next.js 13, TypeScript |
| **express** | Express API | Railway | Express.js, Node.js |
| **fullstack** | 全栈应用 | Netlify | React + Express |

## 🗂️ 数据库支持

| 数据库 | 描述 | 使用场景 |
|--------|------|----------|
| **supabase** | PostgreSQL + 实时功能 | 全栈应用 |
| **mongodb** | NoSQL文档数据库 | API服务 |
| **sqlite** | 轻量级文件数据库 | 本地开发 |

## 🎯 智能检测

SOLO模式会根据需求描述自动选择最佳配置：

- **关键词匹配**: React、Next.js、Express、API、博客、静态等
- **复杂度评估**: 简单需求→静态站点，复杂需求→全栈应用
- **平台选择**: 静态内容→GitHub Pages，动态应用→Vercel/Netlify

## 📁 项目结构

生成的项目包含：

```
your-project/
├── solo-config.json      # SOLO模式配置
├── README.md            # 项目文档
├── .gitignore          # Git忽略文件
└── [平台特定配置文件]   # vercel.json, netlify.toml等
```

## 🌍 部署平台

### GitHub Pages
- **适用**: 静态网站
- **成本**: 完全免费
- **域名**: `username.github.io/project-name`
- **HTTPS**: 自动启用

### Vercel
- **适用**: React、Next.js应用
- **成本**: 免费额度充足
- **域名**: `project-name.vercel.app`
- **功能**: 自动HTTPS、CDN、预览部署

### Netlify
- **适用**: 全栈应用、静态站点
- **成本**: 免费额度
- **域名**: `project-name.netlify.app`
- **功能**: 表单处理、无服务器函数

### Railway
- **适用**: Express API、后端服务
- **成本**: 500小时免费/月
- **域名**: `project-name.up.railway.app`
- **功能**: 自动SSL、数据库集成

## 🔧 环境要求

- **Go**: 1.19+
- **Node.js**: 16+
- **npm/yarn**: 最新版本
- **Git**: 2.x+
- **平台CLI**: Vercel CLI、Railway CLI（可选）

## 🚦 使用流程

1. **需求输入**: 用自然语言描述项目需求
2. **智能解析**: AI分析需求并选择最佳方案
3. **项目生成**: 自动创建项目结构和文件
4. **依赖安装**: 自动安装所需依赖
5. **配置生成**: 生成平台特定的部署配置
6. **部署准备**: 提供详细的部署步骤

## 📊 性能指标

- **生成时间**: 30秒-2分钟
- **部署时间**: 5-15分钟
- **项目大小**: 最小50KB，最大10MB
- **依赖数量**: 根据项目类型自动优化

## 🛠️ 高级用法

### 自定义配置

在需求描述中添加特定要求：

```bash
# 指定技术栈
go run solo/main.go "使用Next.js和TypeScript创建博客，部署到Vercel"

# 指定数据库
go run solo/main.go "Express API需要MongoDB数据库支持"

# 指定平台
go run solo/main.go "静态网站部署到GitHub Pages"
```

### 环境变量

生成的项目会自动包含环境变量模板：

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑配置
vim .env
```

## 🐛 故障排除

### 常见问题

1. **依赖安装失败**
   - 检查Node.js版本
   - 清理npm缓存: `npm cache clean --force`

2. **部署失败**
   - 检查平台CLI是否安装
   - 验证配置文件格式
   - 查看部署日志

3. **项目生成错误**
   - 检查需求描述是否清晰
   - 确保有足够的磁盘空间
   - 验证文件权限

### 调试模式

```bash
# 启用详细日志
export DEBUG=1
go run solo/main.go "你的需求描述"
```

## 📈 路线图

- [ ] 更多项目模板
- [ ] Docker支持
- [ ] 数据库迁移工具
- [ ] 性能优化
- [ ] 国际化支持
- [ ] 插件系统

## 🤝 贡献

欢迎提交Issue和Pull Request来改进SOLO模式！

## 📄 许可证

MIT License - 详见项目根目录LICENSE文件