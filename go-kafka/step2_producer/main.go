// package main

// import (
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/IBM/sarama"
// )

// func main() {
// 	brokers := []string{"192.168.0.5:9092"}
// 	topic := "my-topic"

// 	// 测试同步发送
// 	syncDuration := testSyncProducer(brokers, topic, 100)
// 	fmt.Printf("同步发送 100 条耗时: %v\n", syncDuration)

// 	// 测试异步发送
// 	asyncDuration := testAsyncProducer(brokers, topic, 100)
// 	fmt.Printf("异步发送 100 条耗时: %v\n", asyncDuration)

// 	fmt.Printf("\n速度提升: %.2f 倍\n", float64(syncDuration)/float64(asyncDuration))
// }

// // 同步生产者
// func testSyncProducer(brokers []string, topic string, count int) time.Duration {
// 	config := sarama.NewConfig()
// 	config.Producer.RequiredAcks = sarama.WaitForAll
// 	config.Producer.Return.Successes = true

// 	producer, err := sarama.NewSyncProducer(brokers, config)
// 	if err != nil {
// 		log.Fatalln("创建同步生产者失败:", err)
// 	}
// 	defer producer.Close()

// 	start := time.Now()

// 	for i := 0; i < count; i++ {
// 		msg := &sarama.ProducerMessage{
// 			Topic: topic,
// 			Value: sarama.StringEncoder(fmt.Sprintf("sync msg %d", i)),
// 		}
// 		// 同步发送：等待 Kafka 确认后才返回
// 		_, _, err := producer.SendMessage(msg)
// 		if err != nil {
// 			log.Printf("同步发送失败: %v\n", err)
// 		}
// 	}

// 	return time.Since(start)
// }

// // 异步生产者
// func testAsyncProducer(brokers []string, topic string, count int) time.Duration {
// 	config := sarama.NewConfig()
// 	config.Producer.RequiredAcks = sarama.WaitForAll
// 	config.Producer.Return.Successes = true
// 	config.Producer.Return.Errors = true

// 	producer, err := sarama.NewAsyncProducer(brokers, config)
// 	if err != nil {
// 		log.Fatalln("创建异步生产者失败:", err)
// 	}
// 	defer producer.Close()

// 	// 统计完成数量
// 	successCount := 0
// 	errorCount := 0
// 	done := make(chan bool)

// 	// 后台消费回执
// 	go func() {
// 		for successCount+errorCount < count {
// 			select {
// 			case <-producer.Successes():
// 				successCount++
// 			case <-producer.Errors():
// 				errorCount++
// 			}
// 		}
// 		done <- true
// 	}()

// 	start := time.Now()

// 	// 异步发送：放入通道后立即返回
// 	for i := 0; i < count; i++ {
// 		msg := &sarama.ProducerMessage{
// 			Topic: topic,
// 			Value: sarama.StringEncoder(fmt.Sprintf("async msg %d", i)),
// 		}
// 		producer.Input() <- msg
// 	}

// 	// 等待所有消息确认完成
// 	<-done

// 	return time.Since(start)
// }
