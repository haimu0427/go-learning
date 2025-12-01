import os
from PIL import Image

def remove_black_background(input_dir="output_slices", output_dir="output_transparent"):
    """
    去除文件夹中所有图片的黑色背景，并将结果保存到新文件夹中。
    
    :param input_dir: 输入图片文件夹路径
    :param output_dir: 输出图片文件夹路径
    """
    
    # 创建输出目录
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)
        print(f"创建输出目录: {output_dir}")
    
    # 获取输入目录中的所有PNG文件
    if not os.path.exists(input_dir):
        print(f"错误: 输入目录 {input_dir} 不存在")
        return
    
    png_files = [f for f in os.listdir(input_dir) if f.lower().endswith('.png')]
    
    if not png_files:
        print(f"错误: 在 {input_dir} 中没有找到PNG文件")
        return
    
    print(f"找到 {len(png_files)} 个PNG文件需要处理")
    
    processed_count = 0
    
    for filename in png_files:
        input_path = os.path.join(input_dir, filename)
        output_path = os.path.join(output_dir, filename)
        
        try:
            # 打开图片
            img = Image.open(input_path)
            
            # 如果图片已经有透明通道，直接复制
            if img.mode in ('RGBA', 'LA'):
                background = img.convert('RGBA')
            else:
                # 转换为RGBA模式，添加透明通道
                background = img.convert('RGBA')
            
            # 创建一个新的透明图片
            transparent_img = Image.new('RGBA', background.size, (0, 0, 0, 0))
            
            # 获取图片的像素数据
            pixels = background.load()
            transparent_pixels = transparent_img.load()
            
            # 阈值设置 - 判断黑色像素的阈值
            # 可以根据需要调整这个值 (0-255)
            black_threshold = 30  # RGB值都小于30被认为是黑色
            
            width, height = background.size
            
            # 遍历每个像素
            for x in range(width):
                for y in range(height):
                    r, g, b, a = pixels[x, y]
                    
                    # 如果像素接近黑色，设置为透明
                    if r < black_threshold and g < black_threshold and b < black_threshold:
                        transparent_pixels[x, y] = (0, 0, 0, 0)  # 完全透明
                    else:
                        # 保留原始像素
                        transparent_pixels[x, y] = (r, g, b, a)
            
            # 保存处理后的图片
            transparent_img.save(output_path, 'PNG')
            
            print(f"处理完成: {filename}")
            processed_count += 1
            
        except Exception as e:
            print(f"处理 {filename} 时出错: {e}")
    
    print(f"\n成功处理了 {processed_count} 个图片！")
    print(f"处理后的图片保存在: {output_dir}")

def batch_remove_black_background():
    """
    批量处理多个文件夹中的图片
    """
    # 默认处理 output_slices 文件夹
    remove_black_background()

# --- 如何使用 ---
if __name__ == "__main__":
    print("开始去除图片黑色背景...")
    
    # 处理默认文件夹
    batch_remove_black_background()
    
    print("\n处理完成！")