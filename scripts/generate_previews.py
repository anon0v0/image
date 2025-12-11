#!/usr/bin/env python3
"""
批量生成缺失的预览图脚本
遍历 uploads 目录，为所有缺少 _preview.webp 的图片生成预览图
"""

import os
import sys
from pathlib import Path
from PIL import Image

def generate_preview(image_path, quality=75):
    """
    生成预览图
    """
    try:
        # 打开原始图片
        with Image.open(image_path) as img:
            # 转换为 RGB 模式（WebP 需要）
            if img.mode in ('RGBA', 'LA', 'P'):
                # 创建白色背景
                background = Image.new('RGB', img.size, (255, 255, 255))
                if img.mode == 'P':
                    img = img.convert('RGBA')
                background.paste(img, mask=img.split()[-1] if img.mode in ('RGBA', 'LA') else None)
                img = background
            elif img.mode != 'RGB':
                img = img.convert('RGB')
            
            # 计算缩放后的尺寸（保持宽高比，最大 1920x1080）
            max_width = 1920
            max_height = 1080
            img.thumbnail((max_width, max_height), Image.Resampling.LANCZOS)
            
            # 生成预览图路径
            preview_path = str(image_path).rsplit('.', 1)[0] + '_preview.webp'
            
            # 保存为 WebP
            img.save(preview_path, 'WEBP', quality=quality)
            return True, preview_path
    except Exception as e:
        return False, str(e)

def main():
    uploads_dir = Path('/root/oneimg/uploads')
    
    if not uploads_dir.exists():
        print(f"错误: uploads 目录不存在: {uploads_dir}")
        sys.exit(1)
    
    print("开始扫描并生成缺失的预览图...")
    print(f"扫描目录: {uploads_dir}")
    print()
    
    processed = 0
    skipped = 0
    errors = 0
    
    # 支持的图片扩展名
    image_extensions = {'.jpg', '.jpeg', '.png', '.webp', '.gif'}
    
    # 遍历所有文件
    for image_path in uploads_dir.rglob('*'):
        if not image_path.is_file():
            continue
        
        # 跳过已有的缩略图和预览图
        if '_thumb' in image_path.stem or '_preview' in image_path.stem:
            continue
        
        # 只处理图片文件
        if image_path.suffix.lower() not in image_extensions:
            continue
        
        # 检查是否已有预览图
        preview_path = image_path.parent / (image_path.stem + '_preview.webp')
        if preview_path.exists():
            skipped += 1
            continue
        
        print(f"处理: {image_path.relative_to(uploads_dir)}")
        
        # 生成预览图
        success, result = generate_preview(image_path)
        
        if success:
            print(f"  ✓ 已生成预览图")
            processed += 1
        else:
            print(f"  ✗ 失败: {result}")
            errors += 1
    
    print()
    print("完成！")
    print(f"  已处理: {processed}")
    print(f"  已跳过: {skipped}")
    print(f"  错误数: {errors}")

if __name__ == '__main__':
    main()
