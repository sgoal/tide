# Tide SOLO模式集成指南

本指南说明如何将SOLO模式集成到Tide主程序中，使其成为Tide的一个功能模块。

## 🔗 集成方式

### 方式1: 作为Tide Agent的新工具

将SOLO模式作为Tide Agent的一个新工具，用户可以通过自然语言触发：

```go
// 在agent/agent.go中添加SOLO工具
func registerSoloTool(agent *ReActAgent) {
    agent.Tools["solo_mode"] = &tool.Tool{
        Name:        "solo_mode",
        Description: "根据需求描述自动生成完整项目并部署",
        Func:        soloToolHandler,
    }
}
```

### 方式2: 作为独立CLI命令

在Tide主程序中添加CLI命令：

```go
// 在main.go中添加
import "github.com/sgoal/tide/solo"

func handleSoloCommand(args []string) {
    if len(args) < 1 {
        fmt.Println("Usage: tide solo 'your requirement'")
        return
    }
    
    manager := solo.NewSoloManager(os.Stdout)
    requirement := strings.Join(args, " ")
    
    if err := manager.StartSoloMode(requirement); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## 🏗️ 代码集成示例

### 1. 扩展agent.go

```go
// 在agent/agent.go中添加SOLO模式支持

// 添加新的工具类型
const (
    ToolTypeSolo = "solo_mode"
)

// 在NewReActAgent中添加SOLO工具
func NewReActAgent(openaiKey string) *ReActAgent {
    // ... 现有代码 ...
    
    // 添加SOLO模式工具
    agent.Tools[ToolTypeSolo] = &tool.Tool{
        Name:        "solo_mode",
        Description: "根据自然语言需求自动生成完整项目并部署到云平台",
        Func: func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
            requirement, ok := args["requirement"].(string)
            if !ok {
                return nil, fmt.Errorf("requirement参数缺失")
            }
            
            manager := solo.NewSoloManager(os.Stdout)
            if err := manager.StartSoloMode(requirement); err != nil {
                return nil, err
            }
            
            return map[string]interface{}{
                "status": "completed",
                "message": "项目已生成并准备部署",
            }, nil
        },
    }
    
    return agent
}

// 添加SOLO模式触发关键词
func (a *ReActAgent) shouldUseSoloMode(command string) bool {
    keywords := []string{
        "创建项目", "生成网站", "构建应用", "开发应用",
        "制作网站", "搭建服务", "新建项目", "一键部署",
        "solo模式", "快速开发", "自动生成",
    }
    
    command = strings.ToLower(command)
    for _, keyword := range keywords {
        if strings.Contains(command, keyword) {
            return true
        }
    }
    return false
}
```

### 2. 更新命令处理

```go
// 在ProcessCommand方法中添加SOLO模式检测
func (a *ReActAgent) ProcessCommand(ctx context.Context, command string) (string, error) {
    // 检查是否触发SOLO模式
    if a.shouldUseSoloMode(command) {
        result, err := a.Tools[ToolTypeSolo].Func(ctx, map[string]interface{}{
            "requirement": command,
        })
        if err != nil {
            return "", err
        }
        return fmt.Sprintf("SOLO模式已启动: %v", result), nil
    }
    
    // ... 现有逻辑 ...
}
```

## 🎯 使用场景

### 场景1: 用户直接请求

```
用户: "创建一个React博客网站"
Tide: 启动SOLO模式 → 生成React项目 → 配置部署 → 提供部署指南
```

### 场景2: 复杂项目需求

```
用户: "我需要开发一个全栈任务管理应用，有React前端和Express后端"
Tide: 启动SOLO模式 → 分析需求 → 生成全栈项目 → 配置数据库 → 部署准备
```

### 场景3: 快速原型

```
用户: "快速做一个产品展示页面"
Tide: 启动SOLO模式 → 生成静态网站 → 配置GitHub Pages → 5分钟内完成部署
```

## 🔧 配置集成

### 1. 配置文件

创建配置文件 `solo-config.yaml`：

```yaml
# SOLO模式配置
solo:
  default_platform: "vercel"
  supported_platforms:
    - "github-pages"
    - "vercel"
    - "netlify"
    - "railway"
  default_database: "supabase"
  supported_databases:
    - "supabase"
    - "mongodb"
    - "sqlite"
  templates_path: "./solo/templates"
  output_path: "./generated"
```

### 2. 环境变量

```bash
# .env.example
SOLO_DEFAULT_PLATFORM=vercel
SOLO_DEFAULT_DATABASE=supabase
SOLO_ENABLE_ADVANCED=true
```

## 📊 监控和日志

### 1. 集成日志系统

```go
// 在solo/manager.go中添加日志集成
import "github.com/sgoal/tide/logger"

type SoloManager struct {
    logger *logger.Logger
    // ... 其他字段 ...
}

func NewSoloManager(logger *logger.Logger) *SoloManager {
    return &SoloManager{
        logger: logger,
        // ... 其他初始化 ...
    }
}
```

### 2. 性能监控

```go
// 添加性能指标收集
func (m *SoloManager) collectMetrics() {
    metrics := map[string]interface{}{
        "generation_time": time.Since(startTime),
        "project_type": m.config.ProjectType,
        "platform": m.config.Platform,
        "success": true,
    }
    
    m.logger.Info("SOLO模式执行完成", metrics)
}
```

## 🚀 快速测试

### 1. 测试命令

```bash
# 测试SOLO模式集成
go run main.go "创建一个React博客"

# 测试特定场景
go run main.go "构建Express API"

# 测试部署流程
go run main.go "制作静态网站部署到GitHub Pages"
```

### 2. 集成测试

```bash
# 运行演示脚本
./solo/demo.sh

# 检查生成的项目
ls -la solo-demo-*
```

## 📋 检查清单

### 集成完成检查项

- [ ] SOLO模式工具已注册到Agent
- [ ] CLI命令已添加
- [ ] 触发关键词已配置
- [ ] 日志集成已完成
- [ ] 配置文件已创建
- [ ] 测试用例已验证
- [ ] 文档已更新

### 验证步骤

1. **功能验证**:
   ```bash
   go run main.go "创建一个React应用"
   ```

2. **部署验证**:
   - 检查生成的项目结构
   - 验证配置文件正确性
   - 测试部署流程

3. **集成验证**:
   - 确认Tide Agent能正确调用SOLO模式
   - 验证错误处理机制
   - 测试边界情况

## 🎉 使用示例

### 完整使用流程

```bash
# 1. 启动Tide
./tide

# 2. 用户输入
> 帮我创建一个React博客网站

# 3. Tide响应
🌊 检测到项目创建需求，启动SOLO模式...
🎯 项目类型: React应用
🚀 部署平台: Vercel
📊 数据库: Supabase

# 4. 自动生成项目...
✅ 项目结构已生成
✅ 依赖已安装
✅ 配置文件已创建
✅ 部署准备完成

# 5. 提供部署指南
下一步: 
1. cd react-blog-2024
2. git init && git add . && git commit -m "Initial commit"
3. vercel --prod
```

## 🔗 下一步计划

1. **高级集成**: 与Tide的AI决策系统深度集成
2. **模板扩展**: 增加更多项目模板
3. **智能优化**: 基于用户反馈优化项目生成
4. **社区模板**: 支持用户自定义模板
5. **监控仪表板**: 提供项目部署状态监控

---

通过以上集成，SOLO模式将成为Tide的核心功能之一，为用户提供真正的一键项目创建和部署体验。