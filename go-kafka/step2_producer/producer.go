// package main

// import (
// 	"fmt"
// 	"time"

// 	"github.com/Shopify/sarama"
// )

// func main() {
// 	// 1. 配置生产者
// 	config := sarama.NewConfig()
// 	// 要求所有副本都确认消息，确保消息不丢失
// 	config.Producer.RequiredAcks = sarama.WaitForAll
// 	// 随机选择一个分区发送消息
// 	config.Producer.Partitioner = sarama.NewRandomPartitioner
// 	// 启用成功交付的回调
// 	config.Producer.Return.Successes = true

// 	// Kafka 集群地址
// 	brokerList := []string{"127.0.0.1:9092"} // 请替换为你的 Kafka 地址

// 	// 2. 创建同步生产者
// 	producer, err := sarama.NewSyncProducer(brokerList, config)
// 	if err != nil {
// 		fmt.Printf("创建生产者失败, 错误: %v\n", err)
// 		return
// 	}
// 	defer producer.Close()

// 	topic := "test_topic" // 你的主题名称

// 	// 3. 构造并发送消息
// 	msg := &sarama.ProducerMessage{
// 		Topic: topic,
// 		Key:   sarama.StringEncoder("key"), // 可选
// 		Value: sarama.StringEncoder("Hello, Kafka from Go! " + time.Now().String()),
// 	}

// 	// 4. 发送消息
// 	partition, offset, err := producer.SendMessage(msg)
// 	if err != nil {
// 		fmt.Printf("发送消息失败, 错误: %v\n", err)
// 		return
// 	}

// 	fmt.Printf("消息发送成功到 Topic: %s, Partition: %d, Offset: %d\n", topic, partition, offset)
// }
