package solve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// ParseInput 解析 operate 输入的 JSON 字节流。
//
// 未知字段会被拒绝（DisallowUnknownFields），字段类型不符、
// 语法错误与文件缺失都返回 *ParseError，供 CLI 层直接写 stderr。
func ParseInput(data []byte) (*Input, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var in Input
	if err := dec.Decode(&in); err != nil {
		return nil, &ParseError{Reason: fmt.Sprintf("invalid JSON: %v", err)}
	}
	// 拒绝尾随的非空白内容，避免把拼凑的多段 JSON 当一份输入。
	if err := ensureEOF(dec); err != nil {
		return nil, err
	}
	if err := in.Validate(); err != nil {
		return nil, err
	}
	return &in, nil
}

// ParseInputFile 读取并解析 JSON 文件。
func ParseInputFile(path string) (*Input, error) {
	data, err := readFile(path)
	if err != nil {
		return nil, &ParseError{Reason: err.Error()}
	}
	return ParseInput(data)
}

// ensureEOF 确认解码器已经读到流末尾。
func ensureEOF(dec *json.Decoder) error {
	var extra any
	err := dec.Decode(&extra)
	if err == io.EOF {
		return nil
	}
	if err == nil {
		return &ParseError{Reason: "unexpected extra JSON content after the first object"}
	}
	return &ParseError{Reason: fmt.Sprintf("invalid trailing content: %v", err)}
}

// readFile 读取文件并返回字节，依赖 os.ReadFile 以便测试可注入。
var readFile = osReadFile
