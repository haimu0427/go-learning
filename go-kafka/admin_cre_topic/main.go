package main

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func main() {
	// 1. 配置
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0

	// 替换为你的公网 IP
	brokers := []string{"192.168.0.5:9092"}

	// 2. 创建 ClusterAdmin (管理员对象)
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		log.Fatalln("无法创建集群管理员:", err)
	}
	defer admin.Close()

	// 3. 定义 Topic 详情
	topicName := "web-logs-multi" // 起个新名字
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     3, // 关键点：创建 3 个分区！
		ReplicationFactor: 1, // 因为是单机 Kafka，副本只能设为 1
	}

	// 4. 创建 Topic
	err = admin.CreateTopic(topicName, topicDetail, false)
	if err != nil {
		// 如果报错 Topic 已存在，可以忽略
		log.Printf("创建 Topic 结果: %v\n", err)
	} else {
		fmt.Println("Topic 创建成功！名称:", topicName, "分区数:", 3)
	}
}
