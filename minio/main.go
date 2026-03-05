package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	ctx := context.Background()

	// 1. 配置连接信息
	endpoint := "127.0.0.1:9000"     // OrbStack 映射的 API 端口
	accessKeyID := "admin"           // 你启动 Docker 时设置的账号
	secretAccessKey := "password123" // 你设置的密码
	useSSL := false                  // 本地开发通常不带 SSL

	// 2. 初始化 MinIO Client
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln("初始化失败:", err)
	}

	log.Printf("成功连接到 MinIO: %s\n", endpoint)

	// 3. 创建一个存储桶 (Bucket)
	bucketName := "my-test-bucket"
	location := "us-east-1" // 默认区域即可

	// 检查桶是否已存在
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err == nil && exists {
		log.Printf("桶 %s 已经存在了\n", bucketName)
	} else {
		// 不存在则创建
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
		if err != nil {
			log.Fatalln("创建桶失败:", err)
		}
		log.Printf("成功创建桶: %s\n", bucketName)
	}

	objectName := "3ing.jpeg"                       // 在云端显示的名字
	filePath := "/Users/agiuser/Pictures/3ing.jpeg" // 本地文件路径
	contentType := "image/jpeg"

	// 检查对象是否已存在
	stat, err := minioClient.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err == nil {
		// 对象已存在，询问用户是否覆盖
		fmt.Printf("⚠️  对象 '%s' 已存在（大小: %d bytes，修改时间: %s）\n",
			objectName, stat.Size, stat.LastModified.Format("2006-01-02 15:04:05"))
		fmt.Print("是否覆盖？(y/n): ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input != "y" && input != "yes" {
			fmt.Println("❌ 已取消上传")
			return
		}
		fmt.Println("✅ 继续覆盖上传...")
	}

	// 上传文件
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln("上传失败:", err)
	}
	fmt.Printf("✅ 上传成功，大小: %d bytes\n", info.Size)

	// 4. 生成预签名 URL
	generatePresignedURL(minioClient, ctx, bucketName, objectName)
}

// generatePresignedURL 生成预签名访问链接
func generatePresignedURL(client *minio.Client, ctx context.Context, bucketName, objectName string) {
	// 设置链接过期时间（24小时）
	expiry := time.Second * 24 * 60 * 60

	// 设置请求参数（可选）
	reqParams := make(url.Values)
	// 强制浏览器下载而不是预览，如需要可取消注释：
	// reqParams.Set("response-content-disposition", "attachment; filename=\"download.jpg\"")

	// 生成预签名 URL
	presignedURL, err := client.PresignedGetObject(ctx, bucketName, objectName, expiry, reqParams)
	if err != nil {
		log.Printf("生成预签名链接失败: %v\n", err)
		return
	}

	fmt.Println("\n===========================================")
	fmt.Printf("🔗 预签名访问链接（24小时内有效）:\n\n%s\n", presignedURL)
	fmt.Println("===========================================")
}
