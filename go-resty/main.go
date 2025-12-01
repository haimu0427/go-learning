package main

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2" // 引入 resty
	"github.com/google/uuid"       // 引入 uuid 库
)

// 定义与 JSON 结构对应的 Go 结构体
// 这一步是后端开发的基本功
type SlideshowResponse struct {
	Slideshow struct {
		Author string `json:"author"`
		Date   string `json:"date"`
		Slides []struct {
			Title string `json:"title"`
			Type  string `json:"type"`
		} `json:"slides"`
		Title string `json:"title"`
	} `json:"slideshow"`
}

// 假设这是 API 返回的错误结构
type APIError struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

func main() {
	// 1. 创建一个 Client (就像雇佣了一个秘书)
	client := resty.New()
	client.
		EnableTrace().
		SetDebug(true).
		SetTimeout(5 * time.Second). // 设置超时时间为 5 秒
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second) // 设置重试策略
	client.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			return r.StatusCode() >= 500
		},
	)

	var result SlideshowResponse

	// 2. 链式调用 (Chainable API)
	// 像写句子一样描述你的需求：
	// "Client, 帮我创建一个 Request, 设置 Header, 设置 Body, 然后 Post 到这个 URL"
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]string{ // 优雅点 1: 直接传 Map 或 Struct，自动序列化
			"name":  "Resty User",
			"skill": "Efficiency",
		}).
		SetResult(&result).
		SetError((&APIError{})).              // 设置错误结果的结构体
		Get("https://httpbin.org/status/404") // 这里用 httpbin.org 模拟 API

	if err != nil {
		panic(err)
	}

	client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		// 假设 getToken() 是你获取动态 Token 的逻辑
		token := "my-secret-token-123"
		// 自动把 Token 塞进 Header
		r.SetAuthToken(token)
		requestID := uuid.New().String()
		r.SetHeader("X-Request-ID", requestID)
		// 甚至可以打印日志
		fmt.Printf(">>> [%s]正在向 %s 发起请求...\n", requestID, r.URL)

		return nil
	})
	// 现在的调用变得极其清爽：
	resp, err = client.R().Get("https://httpbin.org/bearer")
	if err != nil {
		panic(err)
	}

	// 3. 直接获取结果
	// 优雅点 2: 自动处理了 Body 关闭和读取
	fmt.Println("Resty Response:", resp.String())
	// 3. 验证结果
	// 注意：我们这里没有做任何 json.Unmarshal 的操作！
	fmt.Println("HTTP 状态码:", resp.StatusCode())
	fmt.Println("自动解析的标题:", result.Slideshow.Title)
	fmt.Println("自动解析的作者:", result.Slideshow.Author)
	if resp.IsError() {
		apiErr := resp.Error().(*APIError)
		fmt.Println("API 错误代码:", apiErr.ErrorCode)
		fmt.Println("API 错误信息:", apiErr.Message)
	} else {
		fmt.Println("请求成功，无错误。")
		fmt.Println("响应内容:", resp.String())
		//fmt.Println("解析后的标题:", resp.)
	}

	//message of trace
	ti := resp.Request.TraceInfo()
	fmt.Println("DNS耗时", ti.DNSLookup)
	fmt.Println("TCP连接耗时", ti.ConnTime)
	fmt.Println("服务器处理耗时", ti.ServerTime)
	GitHubCLI()
}
