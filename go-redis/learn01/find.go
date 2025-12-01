package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// 假设的用户结构体
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var rdb *redis.Client // 你的 Redis 客户端
var db *sql.DB        // 你的数据库客户端
var ctx = context.Background()

// 你的 Gin/Echo/Chi 里的 Handler
func GetUser(userID string) (*User, error) {

	// 1. 定义 key
	key := "user:" + userID

	// 2. 先查你的小笔记本 (Redis)
	val, err := rdb.Get(ctx, key).Result()

	// 3. 检查查询结果
	if err == nil {
		// 3a. 命中 (Hit)!
		// 从 Redis 拿到的数据是 JSON 字符串，需要反序列化
		var user User
		if err := json.Unmarshal([]byte(val), &user); err == nil {
			log.Println("Cache Hit for", key)
			return &user, nil // 立刻返回
		}
	}

	// 4. 未命中 (Miss) - (err == redis.Nil)
	log.Println("Cache Miss for", key)

	// 5. 亲自跑去档案室 (MySQL) 查
	var user User
	dbErr := db.QueryRowContext(ctx,
		"SELECT id, name, age FROM users WHERE id = ?", userID,
	).Scan(&user.ID, &user.Name, &user.Age)

	if dbErr != nil {
		// 数据库里也没查到 (比如用户不存在)
		// 注意：这里可以“缓存空值”防止缓存穿透，但我们先简化
		return nil, dbErr
	}

	// 6. 顺手在小笔记本 (Redis) 上抄一份
	// 序列化成 JSON
	jsonData, _ := json.Marshal(user)

	// SETEX = SET + EXPIRE
	// 设置 10 分钟过期
	rdb.SetEX(ctx, key, jsonData, 10*time.Minute)

	// 7. 把从档案室拿到的新资料返回
	return &user, nil
}
