package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/wcharczuk/go-chart"
)

type Book struct {
	Title string
	Price float64
}

func main() {
	// 读取CSV文件
	file, err := os.Open("./books.csv")
	if err != nil {
		log.Fatalf("无法打开CSV文件: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("读取CSV文件失败: %v", err)
	}

	var books []Book
	for i, record := range records {
		if i == 0 { // 跳过标题行
			continue
		}
		if len(record) < 2 {
			log.Printf("警告: 第%d行数据格式不正确，跳过", i+1)
			continue
		}

		// 处理价格字符串，移除货币符号和空格
		priceStr := strings.TrimSpace(record[1])
		// 移除£符号和其他非数字字符
		priceStr = strings.ReplaceAll(priceStr, "£", "")
		priceStr = strings.ReplaceAll(priceStr, ",", "")
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			log.Printf("警告: 第%d行价格解析失败: %v，跳过", i+1, err)
			continue
		}

		books = append(books, Book{
			Title: record[0],
			Price: price,
		})
	}

	fmt.Printf("成功读取 %d 本图书数据\n", len(books))

	// 创建价格分布直方图
	continuousSeries := createHistogramSeries(books)

	// 创建图表
	graph := chart.Chart{
		Title: "图书价格分布直方图",
		TitleStyle: chart.Style{
			FontSize:  16,
			FontColor: chart.ColorBlack,
		},
		XAxis: chart.XAxis{
			Name: "价格 (£)",
			NameStyle: chart.Style{
				FontSize: 12,
			},
		},
		YAxis: chart.YAxis{
			Name: "图书数量",
			NameStyle: chart.Style{
				FontSize: 12,
			},
		},
		Series: []chart.Series{
			continuousSeries,
		},
	}

	// 保存图表
	f, err := os.Create("price_distribution_fixed.png")
	if err != nil {
		log.Fatalf("创建输出文件失败: %v", err)
	}
	defer f.Close()

	err = graph.Render(chart.PNG, f)
	if err != nil {
		log.Fatalf("渲染图表失败: %v", err)
	}

	fmt.Println("图表已保存为 price_distribution_fixed.png")
}

// createHistogramSeries 创建直方图序列
func createHistogramSeries(books []Book) chart.ContinuousSeries {
	// 提取价格数据
	var prices []float64
	for _, book := range books {
		prices = append(prices, book.Price)
	}

	// 排序价格
	sort.Float64s(prices)

	// 创建直方图数据点
	minPrice := prices[0]
	maxPrice := prices[len(prices)-1]

	// 定义价格区间
	numBins := 10
	binWidth := (maxPrice - minPrice) / float64(numBins)

	var xValues []float64
	var yValues []float64

	// 统计每个区间的数量
	for i := 0; i < numBins; i++ {
		binStart := minPrice + float64(i)*binWidth
		binEnd := minPrice + float64(i+1)*binWidth

		// 计算区间中点作为X值
		binMid := (binStart + binEnd) / 2
		xValues = append(xValues, binMid)

		// 统计落在该区间的价格数量
		count := 0
		for _, price := range prices {
			if price >= binStart && (price < binEnd || (i == numBins-1 && price <= binEnd)) {
				count++
			}
		}
		yValues = append(yValues, float64(count))
	}

	return chart.ContinuousSeries{
		Name:    "价格分布",
		XValues: xValues,
		YValues: yValues,
	}
}
