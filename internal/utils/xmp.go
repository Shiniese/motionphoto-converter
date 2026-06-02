package utils

import (
	"bytes"
	"encoding/xml"
	"io"
	"strings"
)

// extractXMPData 从文件数据中提取 XMP 数据块
func extractXMPData(data []byte) []byte {
	// 1. 优先查找标准的 <x:xmpmeta> 标签
	startTag := []byte("<x:xmpmeta")
	endTag := []byte("</x:xmpmeta>")

	startIdx := bytes.Index(data, startTag)
	if startIdx != -1 {
		endIdx := bytes.Index(data[startIdx:], endTag)
		if endIdx != -1 {
			return data[startIdx : startIdx+endIdx+len(endTag)]
		}
	}

	// 2. Fallback: 尝试查找 <?xpacket ... ?> 标签
	startTagPacket := []byte("<?xpacket begin=")
	endTagPacket := []byte("<?xpacket end=")

	startIdx = bytes.Index(data, startTagPacket)
	if startIdx != -1 {
		endIdx := bytes.Index(data[startIdx:], endTagPacket)
		if endIdx != -1 {
			rest := data[startIdx+endIdx:]
			// 确保包含结尾的 "?>"，形成合法的闭合 XML
			endOfPacket := bytes.Index(rest, []byte("?>"))
			if endOfPacket != -1 {
				return data[startIdx : startIdx+endIdx+endOfPacket+2]
			}
		}
	}
	return nil
}

// ExtractXMPInfo 从文件数据中提取并解析 XMP 信息
func ExtractXMPInfo(data []byte) map[string]string {
	xmpData := extractXMPData(data)
	if xmpData == nil {
		return nil
	}

	// 清理数据：找到第一个 '<' 以剔除可能的 BOM 或无效前缀
	firstBracket := bytes.IndexByte(xmpData, '<')
	if firstBracket != -1 {
		xmpData = xmpData[firstBracket:]
	}

	result := make(map[string]string)
	var itemLengths []string
	var itemPaddings []string
	var itemMimes []string
	var itemSemantics []string

	// 核心：记录 URI 到前缀的映射，用于还原带前缀的完整键名 (如 "GCamera:MicroVideoVersion")
	uriToPrefix := make(map[string]string)
	var currentElement string

	decoder := xml.NewDecoder(bytes.NewReader(xmpData))
	decoder.Strict = false // 容错模式，兼容部分相机生成的不规范 XMP

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			break // 遇到解析错误时，中断并返回已成功解析的部分数据
		}

		switch t := tok.(type) {
		case xml.StartElement:
			// 1. 重建带前缀的元素名 (Element Name)
			elemName := t.Name.Local
			if t.Name.Space != "" && t.Name.Space != "xmlns" && t.Name.Space != "xml" {
				if prefix, ok := uriToPrefix[t.Name.Space]; ok {
					elemName = prefix + ":" + t.Name.Local
				} else if !strings.Contains(t.Name.Space, "/") && !strings.Contains(t.Name.Space, ":") {
					// 降级处理：如果 Space 看起来像前缀而不是 URI (某些非标准 XML)
					elemName = t.Name.Space + ":" + t.Name.Local
				}
			}
			currentElement = elemName

			// 2. 处理属性 (Attributes)
			for _, attr := range t.Attr {
				// 记录 xmlns 声明，建立 URI 到前缀的映射 (例如: "http://..." -> "GCamera")
				if attr.Name.Space == "xmlns" {
					uriToPrefix[attr.Value] = attr.Name.Local
					continue
				}

				// 重建带前缀的属性键名 (Key)
				key := attr.Name.Local
				if attr.Name.Space != "" && attr.Name.Space != "xmlns" && attr.Name.Space != "xml" {
					if prefix, ok := uriToPrefix[attr.Name.Space]; ok {
						key = prefix + ":" + attr.Name.Local
					} else if !strings.Contains(attr.Name.Space, "/") && !strings.Contains(attr.Name.Space, ":") {
						// 降级处理
						key = attr.Name.Space + ":" + attr.Name.Local
					}
				}
				// 注意：如果 XML 中根本没有 xmlns 声明，Go 会直接将 "GCamera:MicroVideoVersion" 整个放入 Local，
				// 此时 key 已经是完整字符串，上述逻辑会安全跳过，保持原样。

				if strings.Contains(key, "MicroVideoOffset") ||
					strings.Contains(key, "ItemLength") ||
					strings.Contains(key, "PresentationTimestampUs") ||
					strings.Contains(key, "MotionPhoto") ||
					strings.Contains(key, "MicroVideo") ||
					strings.Contains(key, "ItemPadding") ||
					strings.Contains(key, "ItemMime") ||
					strings.Contains(key, "ItemSemantic") {
					result[key] = attr.Value
				}

				// 特殊处理 Container Item 属性 (同时支持 "Item:Length" 和 "GContainer:Length")
				switch key {
				case "Item:Length", "GContainer:Length":
					itemLengths = append(itemLengths, attr.Value)
				case "Item:Padding", "GContainer:Padding":
					itemPaddings = append(itemPaddings, attr.Value)
				case "Item:Mime", "GContainer:Mime":
					itemMimes = append(itemMimes, attr.Value)
				case "Item:Semantic", "GContainer:Semantic":
					itemSemantics = append(itemSemantics, attr.Value)
				}
			}

		case xml.CharData:
			trimmed := strings.TrimSpace(string(t))
			if trimmed != "" && currentElement != "" {
				// 检查 currentElement 是否包含我们关心的关键词
				if strings.Contains(currentElement, "MicroVideoOffset") ||
					strings.Contains(currentElement, "ItemLength") ||
					strings.Contains(currentElement, "PresentationTimestampUs") ||
					strings.Contains(currentElement, "MotionPhoto") ||
					strings.Contains(currentElement, "MicroVideo") ||
					strings.Contains(currentElement, "ItemPadding") ||
					strings.Contains(currentElement, "ItemMime") ||
					strings.Contains(currentElement, "ItemSemantic") ||
					currentElement == "Motion Photo" ||
					currentElement == "Motion Photo Version" ||
					currentElement == "Motion Photo Presentation Timestamp Us" ||
					currentElement == "Directory Item Length" ||
					currentElement == "Directory Item Padding" ||
					currentElement == "Directory Item Mime" ||
					currentElement == "Directory Item Semantic" {
					result[currentElement] = trimmed
				}
			}

		case xml.EndElement:
			currentElement = ""
		}
	}

	// 解析完成后，将收集到的 Item 属性组合为标准格式 (保持与 GContainer 格式兼容)
	if len(itemLengths) > 0 {
		joinedLengths := strings.Join(itemLengths, ", ")
		result["Directory Item Length"] = joinedLengths
		result["GContainer:ItemLength"] = joinedLengths
	}
	if len(itemPaddings) > 0 {
		result["Directory Item Padding"] = strings.Join(itemPaddings, ", ")
	}
	if len(itemMimes) > 0 {
		result["Directory Item Mime"] = strings.Join(itemMimes, ", ")
	}
	if len(itemSemantics) > 0 {
		result["Directory Item Semantic"] = strings.Join(itemSemantics, ", ")
	}

	return result
}
