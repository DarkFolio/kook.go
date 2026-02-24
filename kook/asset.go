package kook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// AssetService 媒体资源相关API服务
type AssetService struct {
	client *Client
}

// UploadFile 上传文件
func (s *AssetService) UploadFile(ctx context.Context, filePath string) (*Asset, error) {
	if filePath == "" {
		return nil, fmt.Errorf("文件路径不能为空")
	}

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 读取文件内容
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	fileName := filepath.Base(filePath)
	return s.UploadFileContent(ctx, fileName, fileContent)
}

// UploadFileContent 上传文件内容
func (s *AssetService) UploadFileContent(ctx context.Context, fileName string, content []byte) (*Asset, error) {
	if fileName == "" {
		return nil, fmt.Errorf("文件名不能为空")
	}
	if len(content) == 0 {
		return nil, fmt.Errorf("文件内容不能为空")
	}

	// 创建multipart表单
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加文件字段
	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, fmt.Errorf("创建表单文件失败: %w", err)
	}

	_, err = part.Write(content)
	if err != nil {
		return nil, fmt.Errorf("写入文件内容失败: %w", err)
	}

	writer.Close()

	// 构建请求
	url := s.client.buildURL("asset/create")
	req, err := http.NewRequestWithContext(ctx, "POST", url, &buf)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", s.client.tokenType, s.client.token))
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	s.client.logger.Debugf("上传文件: %s", fileName)

	// 执行请求
	resp, err := s.client.httpClient.Do(req)
	if err != nil {
		s.client.logger.WithError(err).Errorf("上传文件失败")
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		s.client.logger.WithError(err).Errorf("读取上传响应失败")
		return nil, fmt.Errorf("读取上传响应失败: %w", err)
	}

	s.client.logger.Debugf("文件上传响应: %s", string(respBody))

	// 解析响应
	var response Response
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("解析上传响应失败: %w", err)
	}

	// 检查API错误
	if response.Code != 0 {
		err := &APIError{
			Code:    response.Code,
			Message: response.Message,
		}
		s.client.logger.WithError(err).Errorf("文件上传API错误")
		return nil, err
	}

	var asset Asset
	if err := json.Unmarshal(response.Data, &asset); err != nil {
		return nil, fmt.Errorf("解析资源信息失败: %w", err)
	}

	s.client.logger.Infof("文件上传成功: %s -> %s", fileName, asset.URL)
	return &asset, nil
}

// 数据结构定义

// Asset 媒体资源信息
type Asset struct {
	URL  string `json:"url"`  // 资源URL
	Type string `json:"type"` // 资源类型
	Name string `json:"name"` // 文件名
	Size int64  `json:"size"` // 文件大小
}
