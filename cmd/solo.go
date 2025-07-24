package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/sgoal/tide/solo"
)

func main() {
	fmt.Println("🌊 Tide SOLO模式 - 一键项目生成与部署")
	fmt.Println("====================================")
	fmt.Println()
	
	if len(os.Args) < 2 {
		fmt.Println("使用方法:")
		fmt.Println("  go run cmd/solo.go \"你的项目需求描述\"")
		fmt.Println("")
		fmt.Println("示例:")
		fmt.Println("  go run cmd/solo.go \"创建一个React博客网站\"")
		fmt.Println("  go run cmd/solo.go \"构建一个Express REST API\"")
		fmt.Println("  go run cmd/solo.go \"制作一个响应式静态网站\"")
		os.Exit(1)
	}

	// 合并所有参数作为需求描述
	requirement := strings.Join(os.Args[1:], " ")
	
	fmt.Printf("🚀 开始处理需求: %s\n", requirement)
	fmt.Println()

	// 创建SOLO管理器
	manager := solo.NewSoloManager(os.Stdout)
	
	// 启动SOLO模式
	if err := manager.StartSoloMode(requirement); err != nil {
		fmt.Printf("❌ SOLO模式执行失败: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println()
	fmt.Println("🎉 SOLO模式执行完成！")
	fmt.Println("💡 按照上面的步骤操作即可完成部署")
}