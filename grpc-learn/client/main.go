// client/main.go

package main

import (
	"context"
	"io"
	"log"
	"time"

	userpb "grpc-hello/userpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // 导入 "insecure"
)

func main() {
	// 1. 连接到服务器
	// gRPC 默认使用 TLS (SSL) 是加密安全的。
	// 在这个 "Hello, World" 例子中，我们使用 "insecure" 选项来跳过 TLS 验证。
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	// 确保在函数退出时关闭连接
	defer conn.Close()

	// 2. 创建一个 "UserService" 的客户端 "存根 (Stub)"
	// 这是我们与服务器对话的“替身”
	c := userpb.NewUserServiceClient(conn)

	// 准备一个 1 秒超时的上下文 (Context)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 3. 调用远程函数 "GetUser"
	// 就像调用本地函数一样！
	r, err := c.GetUser(ctx, &userpb.GetUserRequest{UserId: "123"})
	if err != nil {
		log.Fatalf("调用 GetUser 失败: %v", err)
	}

	// 4. 打印服务端的响应
	log.Printf("来自服务端的回应: %s (活跃: %v)", r.GetName(), r.GetIsActive())

	// 尝试一个不存在的用户
	r, err = c.GetUser(ctx, &userpb.GetUserRequest{UserId: "999"})
	if err != nil {
		log.Printf("调用 GetUser (999) 失败: %v", err) // 我们预期这里会失败
	} else {
		log.Printf("来自服务端的回应: %s", r.GetName())
	}

	// 阶段三：调用 Server Streaming RPC
	// ===================================
	log.Println("\n--- 开始调用 Server Streaming RPC (GenerateUserReport) ---")

	// 1. 准备请求
	reportReq := &userpb.UserReportRequest{}

	// 2. 调用 RPC。注意！这里返回的是一个 "stream" 和一个 error
	stream, err := c.GenerateUserReport(context.Background(), reportReq)
	if err != nil {
		log.Fatalf("调用 GenerateUserReport 失败: %v", err)
	}

	// 3. 我们必须用一个 "for" 循环来不断地从 "stream" 中接收数据
	for {
		// 4. stream.Recv() 会 "阻塞" 在这里，直到
		//    a) 服务器发来一条新消息
		//    b) 服务器关闭了流 (返回 io.EOF)
		//    c) 发生了一个错误
		line, err := stream.Recv()

		// 5. 判断是否是 "流结束" 的信号
		if err == io.EOF {
			// 服务器说："我发完了"，客户端正常退出循环
			log.Println("--- 服务端数据流结束 ---")
			break
		}

		// 6. 判断是否发生了其他错误
		if err != nil {
			log.Fatalf("接收流数据失败: %v", err)
		}

		// 7. 打印收到的 "一行" 数据
		log.Printf("收到报告行: %s", line.GetLineContent())
	}
}
