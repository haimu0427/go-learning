// package main

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"github.com/Shopify/sarama"
// )

// // 消费者组处理器
// type ConsumerGroupHandler struct{}

// // Setup 在新的会话开始前被调用
// func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
// 	return nil
// }

// // Cleanup 在会话结束时被调用
// func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
// 	return nil
// }

// // ConsumeClaim 实际处理消息的地方
// func (h ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
// 	// 遍历分区中的消息
// 	for message := range claim.Messages() {
// 		fmt.Printf("收到消息 -> Topic: %s, Partition: %d, Offset: %d, Value: %s\n",
// 			message.Topic, message.Partition, message.Offset, string(message.Value))

// 		// 提交位移 (Offset)，标记消息已被处理
// 		session.MarkMessage(message, "")
// 	}
// 	return nil
// }

// func main() {
// 	brokerList := []string{"127.0.0.1:9092"} // 请替换为你的 Kafka 地址
// 	topic := "test_topic"                    // 你的主题名称
// 	groupID := "go-consumer-group-v1"        // 消费者组ID

// 	// 1. 配置消费者
// 	config := sarama.NewConfig()
// 	// 从最旧的位移开始消费 (如果是新组)
// 	config.Consumer.Offsets.Initial = sarama.OffsetOldest

// 	// 2. 创建消费者组
// 	client, err := sarama.NewConsumerGroup(brokerList, groupID, config)
// 	if err != nil {
// 		fmt.Printf("创建消费者组失败, 错误: %v\n", err)
// 		return
// 	}
// 	defer client.Close()

// 	// 3. 监听中断信号，用于优雅退出
// 	signals := make(chan os.Signal, 1)
// 	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	// 4. 开始消费
// 	handler := ConsumerGroupHandler{}
// 	fmt.Println("开始消费...")
// 	go func() {
// 		for {
// 			// Consume() 会一直阻塞，直到会话结束 (Cleanup 被调用) 或 ctx 被取消
// 			err := client.Consume(ctx, []string{topic}, handler)
// 			if err != nil {
// 				fmt.Printf("消费错误: %v\n", err)
// 				return
// 			}
// 			// 检查 context 是否被取消
// 			if ctx.Err() != nil {
// 				return
// 			}
// 		}
// 	}()

// 	// 5. 等待信号退出
// 	<-signals
// 	fmt.Println("\n接收到退出信号，正在停止消费者...")
// }
