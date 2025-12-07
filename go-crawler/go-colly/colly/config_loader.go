package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v3"
)

// ConfigYAML 用于解析YAML的中间结构体
type ConfigYAML struct {
	UserAgent      string   `yaml:"UserAgent"`
	AllowedDomains []string `yaml:"AllowedDomains"`
	Parallelism    int      `yaml:"Parallelism"`
	Delay          string   `yaml:"Delay"` // 延迟时间作为字符串解析
}

// LoadConfig 从YAML文件加载配置
func LoadConfig(filename string) (*SpiderConfig, error) {
	// 读取YAML文件
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 创建中间配置结构体实例
	yamlConfig := &ConfigYAML{}

	// 解析YAML内容到中间结构体
	err = yaml.Unmarshal(yamlFile, yamlConfig)
	if err != nil {
		return nil, fmt.Errorf("解析YAML配置失败: %v", err)
	}


	// 解析延迟时间
	delay, err := time.ParseDuration(yamlConfig.Delay)
	if err != nil {
		log.Printf("警告: 无法解析延迟时间 '%s'，将使用默认值 1s。错误: %v", yamlConfig.Delay, err)
		delay = 1 * time.Second
	}

	// 创建最终配置结构体
	config := &SpiderConfig{
		UserAgent:      yamlConfig.UserAgent,
		AllowedDomains: yamlConfig.AllowedDomains,
		Parallelism:    yamlConfig.Parallelism,
		Delay:          delay,
	}

	// 验证配置
	if config.Parallelism <= 0 {
		log.Printf("警告: 并发数设置为 %d，将使用默认值 5", config.Parallelism)
		config.Parallelism = 5
	}

	if len(config.AllowedDomains) == 0 {
		return nil, fmt.Errorf("错误: 必须至少指定一个允许的域名")
	}

	return config, nil
}

// PrintConfig 打印当前配置（用于调试）
func PrintConfig(config *SpiderConfig) {
	fmt.Println("🔧 ===== 爬虫配置信息 =====")
	fmt.Printf("🌐 User-Agent: %s\n", config.UserAgent)
	fmt.Printf("🎯 允许的域名: %v\n", config.AllowedDomains)
	fmt.Printf("⚡ 并发数: %d\n", config.Parallelism)
	fmt.Printf("⏱️  请求延迟: %v\n", config.Delay)
	fmt.Println("===========================")
}