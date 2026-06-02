package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Shiniese/motionphoto-converter/pkg/motionphoto"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: motionphoto-converter <input_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]

	// 读取文件
	data, err := os.ReadFile(inputFile)
	if err != nil {
		slog.Error("Failed to read file", "err", err)
		os.Exit(1)
	}

	// 模拟从文件提取的 XMP 信息（实际应从 JPEG/HEIC 元数据解析）
	xmpInfo := motionphoto.ExtractXMPInfo(data)

	// 处理动态照片
	result := motionphoto.ProcessMotionPhoto(data, xmpInfo)
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
	os.WriteFile(inputFile + ".jpg", result.Data.ImageData, 0644)
	os.WriteFile(inputFile + ".mp4", result.Data.VideoData, 0644)
	fmt.Printf("\n\n  Image saved to '%s.jpg' \n  Video saved to '%s.mp4' \n", inputFile, inputFile)
}
