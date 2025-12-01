package main

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// 1. 定义数据结构 (The Mechanics)
// 我们只取我们需要这几个字段
type GitHubUser struct {
	Login       string    `json:"login"`
	Name        string    `json:"name"`
	PublicRepos int       `json:"public_repos"`
	Followers   int       `json:"followers"`
	HtmlUrl     string    `json:"html_url"`
	CreatedAt   time.Time `json:"created_at"` // Resty 甚至能自动解析时间字符串！
}

// 定义错误结构
type GitHubError struct {
	Message          string `json:"message"`
	DocumentationUrl string `json:"documentation_url"`
}

func GitHubCLI() {
	// 2. 工程化配置 (The Engineering)
	client := resty.New()

	client.
		SetTimeout(5*time.Second).                         // 永远设置超时
		SetRetryCount(2).                                  // 简单的重试策略
		SetBaseURL("https://api.github.com").              // 设置 BaseURL，以后只写路径即可
		SetHeader("User-Agent", "Resty-Learning-Bot/v1.0") // GitHub 要求要有 User-Agent

	// 3. 开启调试 (The Mastery)
	// 建议：在开发阶段开启，生产环境可以通过环境变量关闭
	client.SetDebug(true)

	// 要查询的用户名
	targetUser := "haimu0427" // 也可以改成你的 GitHub ID

	var user GitHubUser
	var apiError GitHubError

	// 4. 发起请求 (The Action)
	resp, err := client.R().
		SetResult(&user).                // 成功结果容器
		SetError(&apiError).             // 失败结果容器
		SetPathParams(map[string]string{ // 使用路径参数，防止 URL 注入风险
			"username": targetUser,
		}).
		Get("/users/{username}")

	// 5. 处理结果
	if err != nil {
		// 网络层面的错误（如断网、超时）
		fmt.Println("网络爆炸了:", err)
		return
	}

	if resp.IsError() {
		// 业务层面的错误（如 404 用户不存在，403 API限流）
		fmt.Printf("API 报错 (状态码 %d): %s\n", resp.StatusCode(), apiError.Message)
		return
	}

	// 6. 展示成果
	fmt.Println("\n==============================")
	fmt.Printf("用户: %s (%s)\n", user.Login, user.Name)
	fmt.Printf("仓库数: %d\n", user.PublicRepos)
	fmt.Printf("粉丝数: %d\n", user.Followers)
	fmt.Printf("注册时间: %s\n", user.CreatedAt.Format("2006-01-02"))
	fmt.Println("==============================")
}
