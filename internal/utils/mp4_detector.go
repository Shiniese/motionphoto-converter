package utils

import (
	"bytes"
	"encoding/binary"
)

// MP4VideoInfo 存储检测到的 MP4 视频信息
type MP4VideoInfo struct {
	Offset int
	Length int
}

// FindMP4VideoByFileTypeBox 通过搜索 File Type Box (ftyp) 查找 MP4 视频数据
func FindMP4VideoByFileTypeBox(data []byte) *MP4VideoInfo {
	ftypSignature := []byte{0x66, 0x74, 0x79, 0x70} // "ftyp"
	mp4Brand := []byte("mp4")
	
	dataCount := len(data)
	searchIndex := 0
	
	for searchIndex < dataCount-8 {
		ftypRange := bytes.Index(data[searchIndex:], ftypSignature)
		if ftypRange == -1 {
			break
		}
		
		ftypStart := searchIndex + ftypRange
		
		if ftypStart >= 4 {
			boxSizeStart := ftypStart - 4
			if boxSizeStart+4 <= len(data) {
				boxSizeData := data[boxSizeStart:ftypStart]
				boxSize := int(binary.BigEndian.Uint32(boxSizeData))
				
				if boxSize >= 16 && boxSize <= 1024 {
					brandStart := ftypStart + 4
					brandEnd := brandStart + boxSize - 8
					if brandEnd > len(data) {
						brandEnd = len(data)
					}
					
					if brandEnd > brandStart {
						brandData := data[brandStart:brandEnd]
						if bytes.Contains(brandData, mp4Brand) {
							// 找到 MP4 ftyp box，返回从 box 开始到文件末尾的视频数据
							mp4Length := dataCount - boxSizeStart
							return &MP4VideoInfo{
								Offset: boxSizeStart,
								Length: mp4Length,
							}
						}
					}
				}
			}
		}
		
		searchIndex = ftypStart + 4
	}
	
	return nil
}