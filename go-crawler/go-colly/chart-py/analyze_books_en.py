import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import re
from matplotlib import rcParams

# Set font for better display
rcParams['font.family'] = 'DejaVu Sans'
rcParams['axes.unicode_minus'] = False

# Set chart style
try:
    plt.style.use('seaborn-v0_8')
except:
    plt.style.use('seaborn')
sns.set_palette("husl")

def clean_price(price_str):
    """Clean price string and extract numeric value"""
    # Use regex to extract numbers
    match = re.search(r'£(\d+\.?\d*)', price_str)
    if match:
        return float(match.group(1))
    return 0.0

def analyze_books_csv(file_path):
    """Analyze book CSV file"""
    print("Starting book data analysis...")
    
    # Read CSV file
    try:
        df = pd.read_csv(file_path, encoding='utf-8')
    except UnicodeDecodeError:
        df = pd.read_csv(file_path, encoding='gbk')
    
    print(f"Successfully read data, total {len(df)} records")
    print(f"Data columns: {list(df.columns)}")
    
    # Clean price data
    df['price_numeric'] = df['价格'].apply(clean_price)
    
    # Basic statistics
    print("\nBasic Statistics:")
    print(f"Average price: £{df['price_numeric'].mean():.2f}")
    print(f"Highest price: £{df['price_numeric'].max():.2f}")
    print(f"Lowest price: £{df['price_numeric'].min():.2f}")
    print(f"Median price: £{df['price_numeric'].median():.2f}")
    
    # Create charts
    fig, axes = plt.subplots(2, 2, figsize=(15, 12))
    fig.suptitle('Book Data Analysis Report', fontsize=16, fontweight='bold')
    
    # 1. Price distribution histogram
    axes[0, 0].hist(df['price_numeric'], bins=30, alpha=0.7, color='skyblue', edgecolor='black')
    axes[0, 0].set_title('Price Distribution Histogram')
    axes[0, 0].set_xlabel('Price (£)')
    axes[0, 0].set_ylabel('Number of Books')
    axes[0, 0].grid(True, alpha=0.3)
    
    # 2. Price range distribution
    price_bins = [0, 15, 25, 35, 45, 55, 70]
    price_labels = ['£0-15', '£15-25', '£25-35', '£35-45', '£45-55', '£55-70']
    df['price_range'] = pd.cut(df['price_numeric'], bins=price_bins, labels=price_labels, right=False)
    
    price_counts = df['price_range'].value_counts().sort_index()
    axes[0, 1].bar(price_counts.index, price_counts.values, color='lightcoral', alpha=0.8)
    axes[0, 1].set_title('Price Range Distribution')
    axes[0, 1].set_xlabel('Price Range')
    axes[0, 1].set_ylabel('Number of Books')
    axes[0, 1].tick_params(axis='x', rotation=45)
    axes[0, 1].grid(True, alpha=0.3)
    
    # 3. Box plot
    axes[1, 0].boxplot(df['price_numeric'], patch_artist=True, 
                      boxprops=dict(facecolor='lightgreen', alpha=0.7))
    axes[1, 0].set_title('Price Distribution Box Plot')
    axes[1, 0].set_ylabel('Price (£)')
    axes[1, 0].grid(True, alpha=0.3)
    
    # 4. Cumulative distribution
    sorted_prices = sorted(df['price_numeric'])
    cumulative = [i/len(sorted_prices) for i in range(len(sorted_prices))]
    axes[1, 1].plot(sorted_prices, cumulative, color='purple', linewidth=2)
    axes[1, 1].set_title('Price Cumulative Distribution')
    axes[1, 1].set_xlabel('Price (£)')
    axes[1, 1].set_ylabel('Cumulative Proportion')
    axes[1, 1].grid(True, alpha=0.3)
    
    # Adjust layout
    plt.tight_layout()
    
    # Save chart
    plt.savefig('books_analysis_report_en.png', dpi=300, bbox_inches='tight')
    print("Chart saved as 'books_analysis_report_en.png'")
    
    # Show chart
    plt.show()
    
    # Detailed statistical report
    print("\nDetailed Statistical Report:")
    print("=" * 50)
    
    # Price range statistics
    print("\nPrice Range Statistics:")
    for range_name, count in price_counts.items():
        percentage = (count / len(df)) * 100
        print(f"{range_name}: {count} books ({percentage:.1f}%)")
    
    # Top 5 most expensive books
    print("\nTop 5 Most Expensive Books:")
    top_5_expensive = df.nlargest(5, 'price_numeric')[['书名', '价格']]
    for idx, row in top_5_expensive.iterrows():
        print(f"{row['书名'][:50]}... - {row['价格']}")
    
    print("\nTop 5 Cheapest Books:")
    top_5_cheap = df.nsmallest(5, 'price_numeric')[['书名', '价格']]
    for idx, row in top_5_cheap.iterrows():
        print(f"{row['书名'][:50]}... - {row['价格']}")
    
    return df

if __name__ == "__main__":
    # Analyze CSV file
    csv_file = "colly/books.csv"
    df = analyze_books_csv(csv_file)
    
    print("\nAnalysis Complete!")
    print("Tip: Chart has been saved, you can view 'books_analysis_report_en.png' file")