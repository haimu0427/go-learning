package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

// Movie 结构体用于存储电影信息
type Movie struct {
	Title       string
	URL         string
	Rating      string
	Year        string
	Director    string
	Actors      string
	Description string
}

// 配置常量
const (
	BaseURL    = "https://movie.douban.com/top250" // 针对结构稳定的 Top250 页面
	MaxRetries = 3
	Timeout    = 30 * time.Second
	MaxMovies  = 20
	UserAgent  = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
)

func main() {
	fmt.Println("🎬 豆瓣电影爬虫启动...")
	fmt.Println(strings.Repeat("=", 50))

	// 获取电影数据
	movies, err := scrapeDoubanMovies()
	if err != nil {
		log.Fatalf("❌ 爬虫执行失败: %v", err)
	}

	// 打印结果
	printMovies(movies)

	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("✅ 爬取完成！共获取 %d 部电影信息\n", len(movies))
}

// scrapeDoubanMovies 爬取豆瓣电影信息
func scrapeDoubanMovies() ([]Movie, error) {
	var movies []Movie
	var err error

	// 重试机制
	for i := 0; i < MaxRetries; i++ {
		movies, err = fetchMoviesFromDouban()
		if err == nil && len(movies) > 0 {
			break
		}
		if i < MaxRetries-1 {
			fmt.Printf("⚠️  第 %d 次尝试失败，%d 秒后重试...\n", i+1, 2*(i+1))
			time.Sleep(time.Duration(2*(i+1)) * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("经过 %d 次重试后仍然失败: %w", MaxRetries, err)
	}

	return movies, nil
}

// fetchMoviesFromDouban 从豆瓣获取电影数据
func fetchMoviesFromDouban() ([]Movie, error) {
	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout: Timeout,
	}

	// 创建请求
	req, err := http.NewRequest("GET", BaseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Referer", "https://movie.douban.com/")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("状态码错误: %d %s", resp.StatusCode, resp.Status)
	}

	// 读取响应内容到内存，便于调试
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应内容失败: %w", err)
	}

	// 调试：打印部分响应内容预览
	fmt.Println("==== 响应内容预览（前 500 字符）====")
	preview := bodyBytes
	if len(preview) > 500 {
		preview = preview[:500]
	}
	fmt.Println(string(preview))
	fmt.Println("===================================")

	// 使用 bytes.Reader 交给 goquery 解析
	reader := bytes.NewReader(bodyBytes)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("解析 HTML 失败: %w", err)
	}

	// 简单检查页面是否为 Top250 页面
	pageTitle := strings.TrimSpace(doc.Find("title").Text())
	fmt.Println("页面 title:", pageTitle)
	gridCount := doc.Find(".grid_view").Length()
	fmt.Println(".grid_view 容器数量:", gridCount)
	if gridCount == 0 {
		return nil, fmt.Errorf("未找到 .grid_view，当前页面可能不是 Top250 页面")
	}

	// 提取电影信息
	movies := extractMovies(doc)
	return movies, nil
}

// extractMovies 从文档中提取电影信息
func extractMovies(doc *goquery.Document) []Movie {
	// 只使用 Top250 的解析逻辑，确保结构明确
	movies := extractMoviesFromTop250(doc)
	if len(movies) > MaxMovies {
		return movies[:MaxMovies]
	}
	return movies
}

// extractMoviesFromTop250 针对 https://movie.douban.com/top250 的解析
func extractMoviesFromTop250(doc *goquery.Document) []Movie {
	var movies []Movie

	doc.Find(".grid_view .item").Each(func(i int, s *goquery.Selection) {
		if len(movies) >= MaxMovies {
			return
		}

		var m Movie

		// 标题（Top250 通常有两个 .title，第一个为中文）
		titleSel := s.Find(".info .hd .title").First()
		m.Title = strings.TrimSpace(titleSel.Text())
		if m.Title == "" {
			return
		}

		// 链接（封面图上的链接最稳妥）
		if href, ok := s.Find(".pic a").Attr("href"); ok {
			m.URL = strings.TrimSpace(href)
		}

		// 评分
		m.Rating = strings.TrimSpace(s.Find(".info .bd .rating_num").First().Text())

		// 其它信息（导演 / 主演 / 年份等都在 .info .bd p 的第一行）
		infoP := s.Find(".info .bd p").First()
		infoText := infoP.Text()
		infoText = strings.ReplaceAll(infoText, "\n", " ")
		infoText = strings.Join(strings.Fields(infoText), " ")
		m.Year = extractYearFromText(infoText)
		m.Director = infoText // 简单保留整行文本作为信息

		// 简介
		desc := strings.TrimSpace(s.Find(".info .bd .inq").First().Text())
		m.Description = desc

		movies = append(movies, m)
	})

	return movies
}

// extractMovieInfo 从选择器中提取电影信息
func extractMovieInfo(s *goquery.Selection) Movie {
	var movie Movie

	// 提取标题
	movie.Title = extractTitle(s)
	if movie.Title == "" {
		return movie
	}

	// 提取链接
	if href, exists := s.Find("a").First().Attr("href"); exists {
		movie.URL = href
	}

	// 提取评分
	movie.Rating = extractRating(s)

	// 提取年份
	movie.Year = extractYear(s)

	// 提取导演和演员
	movie.Director = extractDirector(s)
	movie.Actors = extractActors(s)

	// 提取简介
	movie.Description = extractDescription(s)

	// 清理数据
	movie.Title = strings.TrimSpace(movie.Title)
	movie.Description = strings.TrimSpace(movie.Description)

	return movie
}

// extractMovieFromLink 从链接中提取电影信息
func extractMovieFromLink(s *goquery.Selection) Movie {
	var movie Movie

	// 提取标题
	movie.Title = strings.TrimSpace(s.Text())
	if movie.Title == "" || utf8.RuneCountInString(movie.Title) < 2 {
		return movie
	}

	// 提取链接
	if href, exists := s.Attr("href"); exists {
		movie.URL = href
	}

	// 尝试从父元素获取更多信息的容器
	parent := s.Parent()
	if parent.Length() > 0 {
		movie.Rating = extractRating(parent)
		movie.Year = extractYear(parent)
		movie.Description = extractDescription(parent)
	}

	return movie
}

// extractTitle 提取标题
func extractTitle(s *goquery.Selection) string {
	// 尝试多种方式提取标题
	titleSelectors := []string{
		".title",
		".name",
		"h3",
		"h2",
		"h1",
		"a",
		"span",
	}

	for _, selector := range titleSelectors {
		title := strings.TrimSpace(s.Find(selector).First().Text())
		if title != "" && utf8.RuneCountInString(title) > 1 {
			return title
		}
	}

	return ""
}

// extractRating 提取评分
func extractRating(s *goquery.Selection) string {
	ratingSelectors := []string{
		".rating",
		".rate",
		".score",
		"[class*='rating']",
		"[class*='score']",
	}

	for _, selector := range ratingSelectors {
		rating := strings.TrimSpace(s.Find(selector).First().Text())
		if rating != "" && isValidRating(rating) {
			return rating
		}
	}

	return ""
}

// extractYear 提取年份
func extractYear(s *goquery.Selection) string {
	yearSelectors := []string{
		".year",
		".date",
		"[class*='year']",
		"[class*='date']",
	}

	for _, selector := range yearSelectors {
		year := strings.TrimSpace(s.Find(selector).First().Text())
		if year != "" && isValidYear(year) {
			return year
		}
	}

	// 尝试从文本中提取年份
	text := s.Text()
	return extractYearFromText(text)
}

// extractDirector 提取导演
func extractDirector(s *goquery.Selection) string {
	directorSelectors := []string{
		".director",
		".dir",
		"[class*='director']",
		"[class*='dir']",
	}

	for _, selector := range directorSelectors {
		director := strings.TrimSpace(s.Find(selector).First().Text())
		if director != "" && len(director) < 50 {
			return director
		}
	}

	return ""
}

// extractActors 提取演员
func extractActors(s *goquery.Selection) string {
	actorSelectors := []string{
		".actor",
		".cast",
		"[class*='actor']",
		"[class*='cast']",
	}

	for _, selector := range actorSelectors {
		actors := strings.TrimSpace(s.Find(selector).First().Text())
		if actors != "" && len(actors) < 100 {
			return actors
		}
	}

	return ""
}

// extractDescription 提取简介
func extractDescription(s *goquery.Selection) string {
	descSelectors := []string{
		".desc",
		".description",
		".summary",
		".intro",
		"[class*='desc']",
		"[class*='summary']",
	}

	for _, selector := range descSelectors {
		desc := strings.TrimSpace(s.Find(selector).First().Text())
		if desc != "" && len(desc) > 10 && len(desc) < 200 {
			return desc
		}
	}

	return ""
}

// extractYearFromText 从文本中提取年份
func extractYearFromText(text string) string {
	// 简单的年份提取逻辑
	for i := 1900; i <= 2030; i++ {
		yearStr := fmt.Sprintf("%d", i)
		if strings.Contains(text, yearStr) {
			return yearStr
		}
	}
	return ""
}

// isValidRating 检查评分是否有效
func isValidRating(rating string) bool {
	// 移除空格和常见符号
	rating = strings.TrimSpace(rating)
	rating = strings.Trim(rating, "分")
	rating = strings.Trim(rating, "/10")
	rating = strings.Trim(rating, "/5")

	// 检查是否为数字
	var score float64
	_, err := fmt.Sscanf(rating, "%f", &score)
	if err != nil {
		return false
	}

	// 检查范围是否合理 (0-10分制)
	return score >= 0 && score <= 10
}

// isValidYear 检查年份是否有效
func isValidYear(year string) bool {
	year = strings.TrimSpace(year)
	year = strings.Trim(year, "()")
	year = strings.Trim(year, "年")

	var y int
	_, err := fmt.Sscanf(year, "%d", &y)
	if err != nil {
		return false
	}

	return y >= 1900 && y <= 2030
}

// printMovies 打印电影信息
func printMovies(movies []Movie) {
	if len(movies) == 0 {
		fmt.Println("⚠️  未找到任何电影信息")
		return
	}

	for i, movie := range movies {
		fmt.Printf("\n🎬 [%d] %s\n", i+1, movie.Title)

		if movie.Year != "" {
			fmt.Printf("   📅 年份: %s\n", movie.Year)
		}

		if movie.Rating != "" {
			fmt.Printf("   ⭐ 评分: %s\n", movie.Rating)
		}

		if movie.Director != "" {
			fmt.Printf("   🎭 导演: %s\n", movie.Director)
		}

		if movie.Actors != "" {
			fmt.Printf("   👥 演员: %s\n", movie.Actors)
		}

		if movie.Description != "" {
			// 限制简介长度
			desc := movie.Description
			if len(desc) > 100 {
				desc = desc[:100] + "..."
			}
			fmt.Printf("   📝 简介: %s\n", desc)
		}

		if movie.URL != "" {
			fmt.Printf("   🔗 链接: %s\n", movie.URL)
		}

		fmt.Println("   " + strings.Repeat("-", 40))
	}
}
