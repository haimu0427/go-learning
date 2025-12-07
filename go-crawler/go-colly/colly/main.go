package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

// 模拟的配置结构 (对应 yaml 配置)
type SpiderConfig struct {
	UserAgent      string        `yaml:"UserAgent"`
	AllowedDomains []string      `yaml:"AllowedDomains"`
	Parallelism    int           `yaml:"Parallelism"` // 并发数
	Delay          time.Duration `yaml:"Delay"`       // 每次请求间隔
}

// BookData 书籍数据结构
type BookData struct {
	Title string `json:"title"` // 书名
	Price string `json:"price"` // 价格
}

func main() {
	fmt.Println("🚀 ===== Go 爬虫启动 =====")

	// 1. 从YAML文件加载配置
	config, err := LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("❌ 加载配置文件失败: %v", err)
	}

	// 打印当前配置（用于调试）
	PrintConfig(config)
	fmt.Println("🎯 开始爬取数据...")

	// 创建书籍数据切片和互斥锁
	var books []BookData
	var mutex sync.Mutex

	// 2. 创建 Collector (爬虫实例)
	c := colly.NewCollector(
		colly.AllowedDomains(config.AllowedDomains...),
		// 开启异步模式 (必须开启，否则 Limit 规则无效)
		colly.Async(true),
	)

	// 3. 加载扩展 (Batteries Included)
	// 爱woc 这是反扒关键
	// 自动随机切换 User-Agent (模拟不同浏览器)
	extensions.RandomUserAgent(c)
	// 自动处理 Referer (有些网站检查图片防盗链需要这个)
	extensions.Referer(c)

	// 4. 【专家级配置】限制并发与延迟
	// 这是防止 IP 被封的最重要一步
	err = c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: config.Parallelism, // 同时最多 5 个 Goroutine 在跑
		RandomDelay: config.Delay,       // 随机延迟，模拟人类行为
	})
	if err != nil {
		log.Fatal(err)
	}

	// 5. 注册回调函数 (事件驱动逻辑)

	// [事件] 请求发出前
	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("🌐 正在访问: %s\n", r.URL.String())
	})

	// [事件] 发生错误
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("❌ 访问失败 %s: %v", r.Request.URL, err)
		// 这里可以加入重试逻辑 r.Request.Retry()
	})

	// [事件] 发现 HTML 中的书籍条目 (类似 Phase 2 的 goquery)
	// 这里的 "article.product_pod" 是目标网站的书籍容器
	c.OnHTML("article.product_pod", func(e *colly.HTMLElement) {
		// 提取数据
		title := e.ChildAttr("h3 a", "title")
		price := e.ChildText("p.price_color")

		// 创建书籍数据
		book := BookData{
			Title: title,
			Price: price,
		}

		// 使用互斥锁保护共享数据
		mutex.Lock()
		books = append(books, book)
		mutex.Unlock()

		// 打印进度信息
		fmt.Printf("📚 书名: %-50s 💰 价格: %s\n", title, price)
	})

	// [事件] 自动翻页 (递归爬取)
	// 找到 "next" 按钮的链接，并访问它
	c.OnHTML("li.next a", func(e *colly.HTMLElement) {
		nextLink := e.Attr("href")
		fmt.Printf("📄 发现下一页，准备跳转: %s\n", nextLink)
		// Visit 会自动处理相对路径拼接
		if err := e.Request.Visit(nextLink); err != nil {
			log.Printf("❌ 访问下一页失败: %v", err)
		}
	})

	// 6. 启动爬虫
	startURL := "http://books.toscrape.com/"
	if err := c.Visit(startURL); err != nil {
		log.Fatalf("❌ 访问起始URL失败: %v", err)
	}

	// 7. 【关键】等待结束
	// 因为开启了 Async，主线程(main)不会自动等待爬虫结束。
	// 必须显式调用 Wait()，否则程序一启动就退出了。
	c.Wait()

	fmt.Printf("✅ 爬虫完成！共爬取 %d 本书籍\n", len(books))

	// 8. 保存数据到CSV文件
	err = saveToCSV(books, "books.csv")
	if err != nil {
		log.Printf("❌ 保存CSV文件失败: %v", err)
	} else {
		fmt.Println("📄 数据已保存到 books.csv 文件")
	}

	fmt.Println("✅ ===== 爬虫任务完成 =====")
}

// saveToCSV 将书籍数据保存到CSV文件
func saveToCSV(books []BookData, filename string) error {
	// 创建CSV文件
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 创建CSV写入器
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入表头
	headers := []string{"书名", "价格"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("写入表头失败: %v", err)
	}

	// 写入数据
	for _, book := range books {
		record := []string{book.Title, book.Price}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("写入数据失败: %v", err)
		}
	}

	return nil
}
