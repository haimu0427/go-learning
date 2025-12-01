package cacheaside

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var baseTTL = 10 * time.Minute
var jitter = time.Duration(300 * time.Second)
var rdb *redis.Client
var db *sql.DB
var ctx = context.Background()

func GetUser(userID int) (*User, error) {
	key := "user:" + string(rune(userID))
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		var user User
		if err := json.Unmarshal([]byte(val), &user); err == nil {
			log.Println("Cache hit")
			return &user, nil
		}
		return nil, err
	}
	log.Printf("Cache miss %v", key)
	lockKey := "lock:" + key
	isLocker, err := rdb.SetNX(ctx, lockKey, "1", 30*time.Second).Result()
	if err != nil {
		return nil, err
	}
	if isLocker {
		log.Println("Acquired lock for key:", key)
		defer rdb.Del(ctx, lockKey)
	} else {
		log.Println("Lock failed, sleep and retry:", key)
		time.Sleep(100 * time.Millisecond)
		return GetUser(userID)
	}
	var user User
	dbErr := db.QueryRowContext(ctx,
		"SELECT id, name, age FROM users WHERE id = ?", userID).
		Scan(&user.ID, &user.Name, &user.Age)

	if dbErr != nil {
		if dbErr == sql.ErrNoRows {
			log.Println("caching null value for key:", key)
			rdb.Set(ctx, key, "", 1*time.Minute)
		}
		return nil, dbErr
	}
	jsonData, _ := json.Marshal(user)
	rdb.Set(ctx, key, jsonData, 10*time.Minute)
	return &user, nil
}
func UpdateUserName(userId string, newName string) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, exeErr := tx.ExecContext(ctx,
		"UPDATE users SET name = ? WHERE id = ?",
		newName, userId)
	if exeErr != nil {
		tx.Rollback()
		return exeErr
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	key := "user:" + userId
	log.Println("Invalidating cache for key:", key)
	rdb.Del(ctx, key)
	return nil
}
