// server/main.go

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// 导入我们刚刚生成的 Go 包
	userpb "grpc-hello/userpb"
)

// 1. 定义一个 struct，它将实现 .proto 文件中定义的 "UserService" 接口。
// 我们必须嵌入 *userpb.UnimplementedUserServiceServer
// 这是 gRPC 强制要求的，以实现“向前兼容”。
type server struct {
	userpb.UnimplementedUserServiceServer
}

// 2. 实现我们的 RPC 方法 "GetUser"
// 它的签名必须和生成的 *pb_grpc.go 文件中的接口完全一致。
func (s *server) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	log.Printf("收到 GetUser 请求，用户 ID: %v", req.GetUserId())

	// 模拟从数据库中获取数据
	// 在真实项目中，你会在这里查询数据库
	if req.GetUserId() == "123" {
		return &userpb.UserResponse{
			UserId:   "123",
			Name:     "Go Web 后端学习者", // 你的名字
			IsActive: true,
		}, nil
	}

	// （在阶段四我们会学习如何返回标准 gRPC 错误）
	return nil, status.Errorf(codes.NotFound, "未找到用户: %s", req.GetUserId())
}

func (s *server) GenerateUserReport(req *userpb.UserReportRequest, stream userpb.UserService_GenerateUserReportServer) error {
	log.Println("收到 GenerateUserReport 请求")

	// 1. 这是一个 "流式" 响应，所以我们用一个循环来模拟
	for i := 1; i <= 10; i++ {
		// 2. 准备 "一行" 报告
		line := &userpb.UserReportLine{
			LineContent: fmt.Sprintf("报告第 %d 行：用户数据...", i),
		}

		// 3. 使用 stream.Send() 发送 "一行" 数据给客户端
		if err := stream.Send(line); err != nil {
			// 如果发送失败 (比如客户端断开了)，记录错误并返回
			log.Printf("发送流数据失败: %v", err)
			return err
		}

		// 4. 模拟耗时操作，比如从数据库查询
		time.Sleep(500 * time.Millisecond)
	}

	// 5. 循环结束，所有数据都发送完毕。
	//    我们只需要 "return nil" 来表示流正常结束。
	log.Println("报告生成完毕")
	return nil
}

func main() {
	// 3. 监听一个 TCP 端口
	lis, err := net.Listen("tcp", ":50051") // 50051 是 gRPC 的常用端口
	if err != nil {
		log.Fatalf("启动监听失败: %v", err)
	}
	log.Printf("服务正在监听: %v", lis.Addr())

	// 4. 创建一个新的 gRPC 服务器实例
	s := grpc.NewServer()

	// 5. 将我们的 "server" 注册到 gRPC 服务器上
	userpb.RegisterUserServiceServer(s, &server{})

	// 6. 启动服务，它会阻塞在这里，等待客户端连接
	if err := s.Serve(lis); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}

}
