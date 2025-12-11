package main

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

// 解码图片
func decodeImage(reader io.Reader) (image.Image, string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, "", err
	}

	// 尝试解码webp
	if img, err := webp.Decode(bytes.NewReader(data)); err == nil {
		return img, "webp", nil
	}

	// 尝试解码gif
	if img, err := gif.Decode(bytes.NewReader(data)); err == nil {
		return img, "gif", nil
	}

	// 尝试解码png
	if img, err := png.Decode(bytes.NewReader(data)); err == nil {
		return img, "png", nil
	}

	// 使用标准库解码其他格式
	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, "", err
	}

	return img, format, nil
}

// 转换为 WebP
func convertToWebP(img image.Image, quality int) ([]byte, error) {
	data, err := webp.EncodeRGBA(img, float32(quality))
	if err != nil {
		return nil, fmt.Errorf("failed to encode webp: %v", err)
	}
	return data, nil
}

// 生成 WebP 预览图
func generateWebPPreview(img image.Image, maxWidth, maxHeight, quality int) ([]byte, error) {
	preview := imaging.Fit(img, maxWidth, maxHeight, imaging.Lanczos)
	return convertToWebP(preview, quality)
}

func main() {
	uploadsDir := "/root/oneimg/uploads"
	
	fmt.Println("开始扫描并生成缺失的预览图...")
	
	processed := 0
	skipped := 0
	errors := 0
	
	// 遍历所有日期目录
	err := filepath.Walk(uploadsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}
		
		// 跳过已有的缩略图和预览图
		if strings.Contains(info.Name(), "_thumb") || strings.Contains(info.Name(), "_preview") {
			return nil
		}
		
		// 只处理图片文件
		ext := strings.ToLower(filepath.Ext(info.Name()))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" && ext != ".gif" {
			return nil
		}
		
		// 检查是否已有预览图
		previewPath := strings.TrimSuffix(path, ext) + "_preview.webp"
		if _, err := os.Stat(previewPath); err == nil {
			skipped++
			return nil
		}
		
		fmt.Printf("处理: %s\n", path)
		
		// 读取图片
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("  ✗ 无法打开文件: %v\n", err)
			errors++
			return nil
		}
		defer file.Close()
		
		// 解码图片
		img, format, err := decodeImage(file)
		if err != nil {
			fmt.Printf("  ✗ 无法解码图片: %v\n", err)
			errors++
			return nil
		}
		
		// 生成预览图
		var previewBytes []byte
		if format == "gif" {
			// GIF 取第一帧生成静态预览图
			previewBytes, err = generateWebPPreview(img, 1920, 1080, 75)
			if err != nil {
				fmt.Printf("  ✗ 生成预览图失败: %v\n", err)
				errors++
				return nil
			}
		} else {
			previewBytes, err = generateWebPPreview(img, 1920, 1080, 75)
			if err != nil {
				fmt.Printf("  ✗ 生成预览图失败: %v\n", err)
				errors++
				return nil
			}
		}
		
		// 保存预览图
		err = os.WriteFile(previewPath, previewBytes, 0644)
		if err != nil {
			fmt.Printf("  ✗ 保存预览图失败: %v\n", err)
			errors++
			return nil
		}
		
		fmt.Printf("  ✓ 已生成预览图\n")
		processed++
		
		return nil
	})
	
	if err != nil {
		fmt.Printf("遍历目录出错: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Printf("\n完成！\n")
	fmt.Printf("  已处理: %d\n", processed)
	fmt.Printf("  已跳过: %d\n", skipped)
	fmt.Printf("  错误数: %d\n", errors)
}
