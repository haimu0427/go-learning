package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

// Consumer 是我们实现的 sarama.ConsumerGroupHandler 结构体
type Consumer struct {
	ready chan bool // 可选：用于通知主程序消费者已就绪
}

// 1. Setup：上岗准备
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// 标记为 ready（如果在通道没关闭的情况下）
	// 注意：实际项目中通常只打印日志
	fmt.Println("Sarama: Consumer Group 会话开始 (Setup)")
	return nil
}

// 2. Cleanup：下班清场
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	fmt.Println("Sarama: Consumer Group 会话结束 (Cleanup - 发生 Rebalance 或关闭)")
	return nil
}

// 3. ConsumeClaim：核心工作区
// 提示：这个方法是在一个新的 Goroutine 中被调用的
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Claim 代表被分配给你的一个或多个分区
	// 注意：不要在循环里打印这句话，否则消息多了会刷屏
	fmt.Printf("开始消费分区 Topic:%s Partition:%d InitialOffset:%d\n",
		claim.Topic(), claim.Partition(), claim.InitialOffset())

	// **核心循环：遍历分配给该 Goroutine 的消息通道**
	// 只要 session.Context() 没有被取消，且 claim.Messages() 没有关闭，就一直读
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				fmt.Printf("通道关闭，退出 Partition %d\n", claim.Partition())
				return nil
			}

			// --- 业务逻辑开始 ---
			fmt.Printf("[消费者] 分区:%d 收到: value=%s, offset=%d\n",
				message.Partition, string(message.Value), message.Offset)

			// 模拟处理耗时（比如写数据库）
			// time.Sleep(10 * time.Millisecond)
			// --- 业务逻辑结束 ---

			// **提交位移（Offset Commitment）**
			// 必须调用 MarkMessage 来标记消息已被成功处理。
			// Sarama 会在后台定期（默认1秒）把这些标记的位移提交给 Kafka。
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			// 收到 Rebalance 或 外部关闭信号，必须立即退出循环，否则会造成死锁或重复消费
			fmt.Printf("收到取消信号，停止消费 Partition %d\n", claim.Partition())
			return nil
		}
	}
}

func main() {
	// 1. 配置
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0                                            // 建议显式指定版本
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange() // 分区分配策略：Range, Sticky, RoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest                       // 如果是新组，从最旧的消息开始读

	// 2. 连接
	brokers := []string{"192.168.0.5:9092"} // 记得替换！
	groupID := "my-learning-group-v1"       // 消费者组 ID

	// 创建消费者组客户端
	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalln("创建消费者组失败:", err)
	}
	defer client.Close()

	// 3. 准备 Context 用于优雅退出
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 4. 消费循环 (这也是新手最容易写错的地方)
	consumer := &Consumer{
		ready: make(chan bool),
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			// **关键点**：client.Consume 是一个阻塞方法。
			// 它可以运行很久。但如果发生了 Rebalance（比如有新消费者加入），
			// 它会返回 nil，你需要再次调用它来重新加入组。
			// 只有当 ctx 被取消，或者发生严重错误时，我们才退出这个死循环。
			if err := client.Consume(ctx, []string{"web-logs-multi"}, consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}

			// 检查是否是用户主动取消
			if ctx.Err() != nil {
				return
			}

			fmt.Println("Rebalance 完成，重新加入消费者组...")
		}
	}()

	fmt.Println("消费者组已启动，按 Ctrl+C 退出...")

	// 5. 监听退出信号
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm // 阻塞直到收到信号

	fmt.Println("接收到退出信号，正在关闭 Context...")
	cancel() // 取消 context，通知 ConsumeClaim 退出循环

	wg.Wait() // 等待后台消费者真正清理完毕
	fmt.Println("消费者优雅退出完毕。")
}
