package main

import (
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	brokers := []string{"192.168.0.5:9092"}
	topic := "my-topic"

	// 测试slowproducer
	slowDuration := testSlowProducer(brokers, topic, 10000)
	fmt.Printf("slowproducer 发送 100 条耗时: %v\n", slowDuration)

	// 测试fastproducer
	fastDuration := testFastProducer(brokers, topic, 10000)
	fmt.Printf("fastproducer 发送 100 条耗时: %v\n", fastDuration)
	fmt.Printf("\n速度提升: %.2f 倍\n", float64(slowDuration)/float64(fastDuration))
}

// slowproducer
func testSlowProducer(brokers []string, topic string, count int) time.Duration {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Flush.Frequency = 0                  // 关闭批量发送
	config.Producer.Compression = sarama.CompressionNone // 关闭压缩

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Fatalln("创建slow生产者失败:", err)
	}
	defer producer.Close()

	start := time.Now()

	// 统计完成数量
	successCount := 0
	errorCount := 0
	done := make(chan bool)

	// 后台消费回执
	go func() {
		for successCount+errorCount < count {
			select {
			case <-producer.Successes():
				successCount++
			case <-producer.Errors():
				errorCount++
			}
		}
		done <- true
	}()

	for i := 0; i < count; i++ {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(fmt.Sprintf("sync msg %d", i)),
		}
		// 异步发送：放入通道后立即返回
		producer.Input() <- msg
	}

	// 等待所有消息确认完成
	<-done

	return time.Since(start)
}

// fastproducer
func testFastProducer(brokers []string, topic string, count int) time.Duration {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = false
	config.Producer.Flush.Frequency = 1 * time.Second    // 每秒发送一次批量
	config.Producer.Flush.Bytes = 16 * 1024              // 达到16KB时发送批量
	config.Producer.Compression = sarama.CompressionGZIP // 关闭压缩

	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		log.Fatalln("创建fast生产者失败:", err)
	}
	defer producer.Close()

	start := time.Now()

	// 异步发送：放入通道后立即返回
	for i := range count {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(fmt.Sprintf("async msg %d", i)),
		}
		producer.Input() <- msg
	}

	return time.Since(start)
}
