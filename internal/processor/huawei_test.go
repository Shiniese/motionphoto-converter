package processor

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Shiniese/motionphoto-converter/internal/models"
	"github.com/stretchr/testify/assert"
)

// getTestFilePath 帮助获取测试文件的绝对路径，假设文件存放在当前目录的 testdata 文件夹下
func getTestFilePath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(currentFile), "testdata", filename)
}

func TestHuaweiMotionPhotoProcessor_CanProcess(t *testing.T) {
	processor := NewHuaweiMotionPhotoProcessor()

	tests := []struct {
		name     string
		xmpInfo  map[string]string
		expected bool
	}{
		{
			name:     "Make is HUAWEI",
			xmpInfo:  map[string]string{"Make": "HUAWEI"},
			expected: true,
		},
		{
			name:     "Manufacturer is HUAWEI",
			xmpInfo:  map[string]string{"Manufacturer": "HUAWEI"},
			expected: true,
		},
		{
			name:     "Neither Make nor Manufacturer is HUAWEI",
			xmpInfo:  map[string]string{"Make": "Apple", "Manufacturer": "Apple"},
			expected: false,
		},
		{
			name:     "Empty XMP info",
			xmpInfo:  map[string]string{},
			expected: false,
		},
		{
			name:     "Case sensitive check (should be false)",
			xmpInfo:  map[string]string{"Make": "huawei"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.CanProcess(tt.xmpInfo)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHuaweiMotionPhotoProcessor_ProcessMotionPhoto_Success(t *testing.T) {
	// 1. 初始化处理器
	processor := NewHuaweiMotionPhotoProcessor()

	// 2. 读取测试文件 (请确保 testdata/HUAWEI.jpeg 存在)
	filePath := getTestFilePath("HUAWEI.jpeg")
	data, err := os.ReadFile(filePath)

	// 如果文件不存在，跳过测试或给出明确提示（在 CI 环境中非常有用）
	if err != nil {
		t.Skipf("Skipping test: test file not found at %s (%v)", filePath, err)
	}

	// 3. 构造符合华为特征的 XMP 信息
	xmpInfo := map[string]string{
		"Make":         "HUAWEI",
		"Model":        "Some Huawei Model",
		"Manufacturer": "HUAWEI",
	}

	// 4. 执行处理
	result := processor.ProcessMotionPhoto(data, xmpInfo)

	// 5. 断言结果
	assert.NotNil(t, result, "Result should not be nil")

	// 假设 NewSuccessResult 返回的结构体中，错误信息为空或 nil，或者有 IsSuccess() 方法
	// 请根据您 models.MotionPhotoProcessingResult 的实际字段进行调整
	assert.Nil(t, result.Error(), "Expected no error, but got: %v", result.Error())
	assert.NotNil(t, result.Data, "MotionPhotoData should not be nil on success")

	// 验证提取的数据
	mpData := result.Data
	assert.Equal(t, models.BrandHuawei, mpData.Brand, "Brand should be Huawei")
	assert.NotNil(t, mpData.VideoOffset, "VideoStartOffset should not be nil for Huawei")
	assert.Nil(t, mpData.PresentationTimestamp, "Timestamp should be nil for Huawei format")

	// 验证数据长度大于 0 (这同时也间接证明了 utils.FindMP4VideoByFileTypeBox 成功找到了视频)
	assert.Greater(t, len(mpData.ImageData), 0, "ImageData length should be greater than 0")
	assert.Greater(t, len(mpData.VideoData), 0, "VideoData length should be greater than 0")

	// 验证截取的总长度是否与原文件一致 (可选，取决于您的业务逻辑是否包含尾部多余数据)
	assert.Equal(t, len(data), len(mpData.ImageData)+len(mpData.VideoData), "Total extracted data length should match original")
}

func TestHuaweiMotionPhotoProcessor_ProcessMotionPhoto_Fail_NoVideo(t *testing.T) {
	processor := NewHuaweiMotionPhotoProcessor()

	// 构造一个不包含 MP4 File Type Box 的假数据
	fakeData := []byte("This is just a regular JPEG without any MP4 video payload appended.")
	xmpInfo := map[string]string{"Make": "HUAWEI"}

	result := processor.ProcessMotionPhoto(fakeData, xmpInfo)

	assert.NotNil(t, result)
	assert.NotNil(t, result.Error(), "Expected an error when no MP4 video is found")
	assert.Contains(t, *result.Error(), "No MP4 video found", "Error message should indicate missing MP4 video")
	assert.Nil(t, result.Data, "MotionPhotoData should be nil on failure")
}
