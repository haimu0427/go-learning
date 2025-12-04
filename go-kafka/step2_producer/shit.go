// go-kafka 简单使用示例
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	// 1. 配置
	config := sarama.NewConfig()
	// 异步生产者通常设为 WaitForLocal (1) 或 NoResponse (0) 以求快
	// 这里设为 WaitForAll 只是为了演示可靠性，实际异步场景看需求
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// 2. 连接局域网 Kafka
	brokers := []string{"192.168.0.5:9092"}

	// 3. 创建异步生产者
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Fatalln("无法创建异步生产者:", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln("关闭生产者失败:", err)
		}
	}()

	// 4. 设置监听器（关键！如果你不读这两个通道，程序会死锁）
	// 我们用 goroutine 来后台处理成功和失败的回执
	go func() {
		for {
			select {
			case success := <-producer.Successes():
				fmt.Printf("成功发送: partition=%d offset=%d\n", success.Partition, success.Offset)
			case err := <-producer.Errors():
				fmt.Printf("发送失败: err=%v\n", err.Err)
			}
		}
	}()

	// 5. 模拟发送消息
	topic := "web-logs-multi" // 确保你在服务器上创建了这个 topic，或者配置了 auto.create.topics.enable=true

	fmt.Println("开始发送消息 (按 Ctrl+C 停止)...")

	// 捕获退出信号
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

Loop:
	for {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(fmt.Sprintf("Log info at %v", time.Now())),
		}

		// 注意：这里是 input <- msg，非阻塞的（除非 buffer 满了）
		select {
		case producer.Input() <- msg:
			time.Sleep(100 * time.Millisecond) // 模拟每0.1秒发一条
		case <-signals:
			break Loop
		}
	}

	fmt.Println("停止发送.")
}
