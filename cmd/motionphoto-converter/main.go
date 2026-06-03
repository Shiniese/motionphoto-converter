package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Shiniese/motionphoto-converter/pkg/motionphoto"
)

// 这些变量会在编译时由Makefile注入值
var (
	AppName   = "unknown" // 应用名称
	Version   = "unknown" // 版本号
	GitCommit = "unknown" // Git提交哈希
	BuildTime = "unknown" // 编译时间
	Author    = "unknown" // 作者
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: motionphoto-converter <input_file>")
		os.Exit(1)
	}

	// 加个 --version 参数打印版本信息
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Printf("%s %s\n", AppName, Version)
		fmt.Printf("Git commit: %s\n", GitCommit)
		fmt.Printf("Build time: %s\n", BuildTime)
		fmt.Printf("Author: %s\n", Author)
		return
	}

	inputFile := os.Args[1]

	// 读取文件
	data, err := os.ReadFile(inputFile)
	if err != nil {
		slog.Error("Failed to read file", "err", err)
		os.Exit(1)
	}

	// 处理动态照片
	result := motionphoto.ProcessMotionPhoto(data)
	if !result.Success {
		slog.Error("Failed to process motion photo", "err", *result.ErrorMessage)
		os.Exit(1)
	}

	// 输出结果
	fmt.Printf("✓ Processed successfully!\n")
	fmt.Printf("  Brand: %s\n", result.Data.Brand.DisplayName())
	fmt.Printf("  Image size: %d bytes\n", len(result.Data.ImageData))
	fmt.Printf("  Video size: %d bytes\n", len(result.Data.VideoData))
	if result.Data.VideoOffset != nil {
		fmt.Printf("  Video offset: %d\n", *result.Data.VideoOffset)
	}

	// 保存提取的图片和视频（示例）
	os.WriteFile(inputFile+".jpg", result.Data.ImageData, 0644)
	os.WriteFile(inputFile+".mp4", result.Data.VideoData, 0644)
	fmt.Printf("\n\n  Image saved to '%s.jpg' \n  Video saved to '%s.mp4' \n", inputFile, inputFile)
}
