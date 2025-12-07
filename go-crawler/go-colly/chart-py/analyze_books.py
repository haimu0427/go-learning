import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import re
import sys
from matplotlib import rcParams

# 设置控制台编码
if sys.platform == 'win32':
    import codecs
    sys.stdout = codecs.getwriter('utf-8')(sys.stdout.detach())

# 设置中文字体支持
rcParams['font.sans-serif'] = ['SimHei', 'DejaVu Sans']
rcParams['axes.unicode_minus'] = False

# 设置图表样式
try:
    plt.style.use('seaborn-v0_8')
except:
    plt.style.use('seaborn')
sns.set_palette("husl")

def clean_price(price_str):
    """清理价格字符串，提取数字"""
    # 使用正则表达式提取数字
    match = re.search(r'£(\d+\.?\d*)', price_str)
    if match:
        return float(match.group(1))
    return 0.0

def analyze_books_csv(file_path):
    """分析图书CSV文件"""
    print("开始分析图书数据...")
    
    # 读取CSV文件
    try:
        df = pd.read_csv(file_path, encoding='utf-8')
    except UnicodeDecodeError:
        df = pd.read_csv(file_path, encoding='gbk')
    
    print(f"成功读取数据，共 {len(df)} 条记录")
    print(f"数据列: {list(df.columns)}")
    
    # 清理价格数据
    df['价格_数值'] = df['价格'].apply(clean_price)
    
    # 基本统计信息
    print("\n基本统计信息:")
    print(f"平均价格: £{df['价格_数值'].mean():.2f}")
    print(f"最高价格: £{df['价格_数值'].max():.2f}")
    print(f"最低价格: £{df['价格_数值'].min():.2f}")
    print(f"价格中位数: £{df['价格_数值'].median():.2f}")
    
    # 创建图表
    fig, axes = plt.subplots(2, 2, figsize=(15, 12))
    fig.suptitle('图书数据分析报告', fontsize=16, fontweight='bold')
    
    # 1. 价格分布直方图
    axes[0, 0].hist(df['价格_数值'], bins=30, alpha=0.7, color='skyblue', edgecolor='black')
    axes[0, 0].set_title('价格分布直方图')
    axes[0, 0].set_xlabel('价格 (£)')
    axes[0, 0].set_ylabel('图书数量')
    axes[0, 0].grid(True, alpha=0.3)
    
    # 2. 价格区间分布
    price_bins = [0, 15, 25, 35, 45, 55, 70]
    price_labels = ['£0-15', '£15-25', '£25-35', '£35-45', '£45-55', '£55-70']
    df['价格区间'] = pd.cut(df['价格_数值'], bins=price_bins, labels=price_labels, right=False)
    
    price_counts = df['价格区间'].value_counts().sort_index()
    axes[0, 1].bar(price_counts.index, price_counts.values, color='lightcoral', alpha=0.8)
    axes[0, 1].set_title('价格区间分布')
    axes[0, 1].set_xlabel('价格区间')
    axes[0, 1].set_ylabel('图书数量')
    axes[0, 1].tick_params(axis='x', rotation=45)
    axes[0, 1].grid(True, alpha=0.3)
    
    # 3. 箱线图
    axes[1, 0].boxplot(df['价格_数值'], patch_artist=True, 
                      boxprops=dict(facecolor='lightgreen', alpha=0.7))
    axes[1, 0].set_title('价格分布箱线图')
    axes[1, 0].set_ylabel('价格 (£)')
    axes[1, 0].grid(True, alpha=0.3)
    
    # 4. 累积分布
    sorted_prices = sorted(df['价格_数值'])
    cumulative = [i/len(sorted_prices) for i in range(len(sorted_prices))]
    axes[1, 1].plot(sorted_prices, cumulative, color='purple', linewidth=2)
    axes[1, 1].set_title('价格累积分布')
    axes[1, 1].set_xlabel('价格 (£)')
    axes[1, 1].set_ylabel('累积比例')
    axes[1, 1].grid(True, alpha=0.3)
    
    # 调整布局
    plt.tight_layout()
    
    # 保存图表
    plt.savefig('books_analysis_report.png', dpi=300, bbox_inches='tight')
    print("图表已保存为 'books_analysis_report.png'")
    
    # 显示图表
    plt.show()
    
    # 详细统计报告
    print("\n详细统计报告:")
    print("=" * 50)
    
    # 价格区间统计
    print("\n价格区间统计:")
    for interval, count in price_counts.items():
        percentage = (count / len(df)) * 100
        print(f"{interval}: {count} 本 ({percentage:.1f}%)")
    
    # 最贵和最便宜的前5本书
    print("\n最贵的5本书:")
    top_5_expensive = df.nlargest(5, '价格_数值')[['书名', '价格']]
    for idx, row in top_5_expensive.iterrows():
        print(f"{row['书名'][:50]}... - {row['价格']}")
    
    print("\n最便宜的5本书:")
    top_5_cheap = df.nsmallest(5, '价格_数值')[['书名', '价格']]
    for idx, row in top_5_cheap.iterrows():
        print(f"{row['书名'][:50]}... - {row['价格']}")
    
    return df

if __name__ == "__main__":
    # 分析CSV文件
    csv_file = "colly/books.csv"
    df = analyze_books_csv(csv_file)
    
    print("\n分析完成！")
    print("提示：图表已保存，你可以查看 'books_analysis_report.png' 文件")