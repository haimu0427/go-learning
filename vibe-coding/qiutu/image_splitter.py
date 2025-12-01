import os
from PIL import Image

def split_image_4x3(image_path, output_dir="output_slices"):
    """
    将一张图片切割成 4 列 x 3 行（共 12 块）。

    :param image_path: 待切割的图片路径。
    :param output_dir: 切片输出的文件夹名称。
    """
    
    # 设定目标网格
    COLS = 4  # 列数 (水平方向)
    ROWS = 3  # 行数 (垂直方向)

    try:
        # 1. 打开图片
        img = Image.open(image_path)
        W, H = img.size
        print(f"✅ 原始图片尺寸: {W} x {H} 像素")

        # 2. 计算基础切片尺寸和剩余像素
        base_w = W // COLS
        base_h = H // ROWS
        
        extra_w = W % COLS # 水平方向剩余像素
        extra_h = H % ROWS # 垂直方向剩余像素

        # 3. 创建输出目录
        if not os.path.exists(output_dir):
            os.makedirs(output_dir)
            
        print(f"📁 切片将保存到: {output_dir}")

        # 4. 遍历并切割图片
        current_y = 0
        slice_count = 0
        
        for r in range(ROWS):
            # 计算当前行的高度 (前 few 行会吸收 extra_h)
            h = base_h + (1 if r < extra_h else 0)
            
            current_x = 0
            for c in range(COLS):
                # 计算当前列的宽度 (前 few 列会吸收 extra_w)
                w = base_w + (1 if c < extra_w else 0)
                
                # 定义裁剪区域 (左上角 x, 左上角 y, 右下角 x, 右下角 y)
                box = (current_x, current_y, current_x + w, current_y + h)
                
                # 裁剪图片
                slice_img = img.crop(box)
                
                # 生成文件名: 原始文件名_r[行号]_c[列号].png
                base_name = os.path.splitext(os.path.basename(image_path))[0]
                output_filename = f"{base_name}_r{r+1}_c{c+1}.png"
                output_path = os.path.join(output_dir, output_filename)
                
                # 保存切片 (使用 PNG 格式以保留质量)
                slice_img.save(output_path)
                
                print(f"   - 保存切片: {output_filename} ({w}x{h})")
                
                current_x += w
                slice_count += 1
            
            current_y += h
            
        print(f"\n🎉 成功切割并保存了 {slice_count} 个切片！")

    except FileNotFoundError:
        print(f"❌ 错误: 找不到文件 {image_path}")
    except Exception as e:
        print(f"❌ 发生错误: {e}")

# --- 如何使用 ---
if __name__ == "__main__":
    # 将这里的 'your_image.jpg' 替换为您要切割的图片路径
    # 确保图片文件与脚本在同一目录下，或提供完整的路径
    IMAGE_TO_SPLIT = "qiutu.jpeg" 
    
    split_image_4x3(IMAGE_TO_SPLIT)