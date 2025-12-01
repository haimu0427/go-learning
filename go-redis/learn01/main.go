package main

import (
	"context"
	"encoding/json"
	"fmt"
	"learn01/contact"

	"github.com/redis/go-redis/v9"
)

func main() {
	//加载环境变量
	if err := contact.LoadEnvlog(); err != nil {
		panic(err)
	}
	//设置空白文本流和Redis客户端
	//用来监控超时或者传递数据, 也可以监控redis性能
	ctx := context.Background()

	option, err := redis.ParseURL(contact.BuildClientURL())
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(option)
	defer rdb.Close()

	// 测试连接
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	// 获取 stus key
	val, err := rdb.Get(ctx, "stus").Result()
	switch {
	case err == redis.Nil:
		fmt.Println("stus key does not exist")
	case err != nil:
		panic(err)
	case val == "":
		fmt.Println("stus key is empty")
	default:
		fmt.Println("stus key:", val)
	}

	// string 操作
	err = rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err = rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key value:", val)

	// hash
	hashFields := []string{
		"model", "demos",
		"brand", "Ergnom",
		"year", "2023",
		"type", "Enduro bikes",
		"price", "4972",
	}
	res1, err := rdb.HSet(ctx, "bike:1", hashFields).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("HSet", res1)

	res2, err := rdb.HGet(ctx, "bike:1", "model").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("HGet", res2)
	res2, err = rdb.HGet(ctx, "bike:1", "brand").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("HGet", res2)
	res2, err = rdb.HGet(ctx, "bike:1", "year").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("HGet", res2)

	res4, err := rdb.HGetAll(ctx, "bike:1").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("HGetAll", res4)

}
func GetUser(userId int) string {
	_, err : = rdb.Get(ctx, fmt.Sprintf("user:%d", userId)).Result()
	if err == nil {
		var user User 
		err = json.Unmarshal([]byte(val), &user)
		if err != nil {
			return user
		}
	}
	user, err := db.QueryUserByID(userId)
	if err != nil {
		return nil
	}
	data, _ := json.Marshal(user)
	rdb.Set(ctx, fmt.Sprintf("user:%d", userId), data, time.Hour)
	return user
}
