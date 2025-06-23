package model

import (
	"context"
	"edulimitrate/config"
	"errors"
)

// 使用 Redis 缓存提升性能
var ctx = context.Background()

func InsertRegionCode(code string) error {
	// 先判断 Redis 里是否存在，减少数据库访问
	exists, err := config.RedisClient.SIsMember(ctx, "ratelimit_regions", code).Result()
	if err == nil && exists {
		return errors.New("region code already exists")
	}

	// MySQL 中插入
	_, err = config.DB.Exec("INSERT IGNORE INTO ratelimit_regions (region_code) VALUES (?)", code)
	if err != nil {
		return err
	}

	// 同步写入 Redis
	return config.RedisClient.SAdd(ctx, "ratelimit_regions", code).Err()
}

func DeleteRegionCode(code string) error {
	res, err := config.DB.Exec("DELETE FROM ratelimit_regions WHERE region_code = ?", code)
	if err != nil {
		return err
	}
	// 不管数据库是否有删除，统一认为操作成功
	_, _ = res.RowsAffected()

	// 从 Redis 删除，忽略错误
	_ = config.RedisClient.SRem(ctx, "ratelimit_regions", code).Err()

	return nil
}

func ExistsRegionCode(code string) (bool, error) {
	// 优先查 Redis 缓存
	exists, err := config.RedisClient.SIsMember(ctx, "ratelimit_regions", code).Result()
	if err == nil {
		return exists, nil
	}

	// Redis 出错或未命中，回退查 MySQL
	var count int
	err = config.DB.QueryRow("SELECT COUNT(*) FROM ratelimit_regions WHERE region_code = ?", code).Scan(&count)
	if err != nil {
		return false, err
	}

	// 同步缓存，忽略错误
	if count > 0 {
		_ = config.RedisClient.SAdd(ctx, "ratelimit_regions", code).Err()
	}

	return count > 0, nil
}
