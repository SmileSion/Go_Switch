package model

import (
	"context"
	"edulimitrate/config"
	"errors"
	"os"
	"strings"
)

// 使用 Redis 缓存提升性能
var ctx = context.Background()

const promptRedisKey = "PromptWords"
const allowedKey = "ratelimit_regions:allowed"
const deniedKey = "ratelimit_regions:denied"

func InsertRegionCode(code string) error {
	

	// 优先检查 Redis 是否已存在
	exists, err := config.RedisClient.SIsMember(ctx, allowedKey, code).Result()
	if err == nil && exists {
		return errors.New("region code already exists")
	}

	// 插入数据库
	_, err = config.DB.Exec("INSERT IGNORE INTO ratelimit_regions (region_code) VALUES (?)", code)
	if err != nil {
		return err
	}

	// 插入成功后同步缓存
	_ = config.RedisClient.SAdd(ctx, allowedKey, code).Err()
	_ = config.RedisClient.SRem(ctx, deniedKey, code).Err() // 移除原来的“不存在”缓存（如果有）

	return nil
}

func DeleteRegionCode(code string) error {

	// 删除数据库记录
	_, err := config.DB.Exec("DELETE FROM ratelimit_regions WHERE region_code = ?", code)
	if err != nil {
		return err
	}
	// 即使 RowsAffected 是 0，我们也认为删除成功

	// 更新 Redis 缓存
	_ = config.RedisClient.SRem(ctx, allowedKey, code).Err()
	_ = config.RedisClient.SAdd(ctx, deniedKey, code).Err() // 加入“已删除”的列表，避免误查

	return nil
}


func ExistsRegionCode(code string) (bool, error) {

	// 优先检查 allowed 列表
	allowed, err := config.RedisClient.SIsMember(ctx, allowedKey, code).Result()
	if err == nil && allowed {
		return true, nil
	}

	// 再检查 denied 列表
	denied, err := config.RedisClient.SIsMember(ctx, deniedKey, code).Result()
	if err == nil && denied {
		return false, nil
	}

	// Redis 不存在或 miss，回退查 MySQL
	var count int
	err = config.DB.QueryRow("SELECT COUNT(*) FROM ratelimit_regions WHERE region_code = ?", code).Scan(&count)
	if err != nil {
		return false, err
	}

	// 分类写入 Redis（避免下次再查数据库）
	if count > 0 {
		_ = config.RedisClient.SAdd(ctx, allowedKey, code).Err()
		return true, nil
	} else {
		_ = config.RedisClient.SAdd(ctx, deniedKey, code).Err()
		return false, nil
	}
}


func GetPromptWords() (string, error) {
	prompt, err := config.RedisClient.Get(ctx, promptRedisKey).Result()
	if err == nil && prompt != "" {
		return prompt, nil
	}
	// 读本地文件
	data, err := os.ReadFile("Promptwords")
	if err != nil {
		return "", errors.New("读取本地提示词文件失败")
	}
	prompt = strings.TrimSpace(string(data))

	// 写 Redis 缓存，设置过期时间，比如24小时
	_ = config.RedisClient.Set(ctx, promptRedisKey, prompt, 0).Err()

	return prompt, nil
}

// SetPromptWords 写提示词到 Redis（如需要可调用）
func SetPromptWords(val string) error {
	return config.RedisClient.Set(ctx, promptRedisKey, val, 0).Err()
}