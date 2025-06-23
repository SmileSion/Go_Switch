package config

import (
	"github.com/go-redis/redis/v8"
	"context"
	"github.com/BurntSushi/toml"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type ServerConfig struct {
	Port int `toml:"port"`
}

type MysqlConfig struct {
	DSN string `toml:"dsn"`
}

type LogConfig struct {
	Filepath string `toml:"filepath"`
	MaxSize    int    `toml:"max_size"`
	MaxBackups int    `toml:"max_backups"`
	MaxAge     int    `toml:"max_age"`
	Compress   bool   `toml:"compress"`
}

type RedisConfig struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

type Config struct {
	Server ServerConfig `toml:"server"`
	Mysql  MysqlConfig  `toml:"mysql"`
	Log    LogConfig    `toml:"log"`
	Redis  RedisConfig  `toml:"redis"`
}

var Conf Config

func InitConfig() {
	if _, err := toml.DecodeFile("config/config.toml", &Conf); err != nil {
		panic(err)
	}

	log.Println("配置文件加载成功：")
	log.Printf("  Server Port: %d\n", Conf.Server.Port)
	log.Printf("  MySQL DSN: %s\n", Conf.Mysql.DSN)
	log.Printf("  Log File: %s\n", Conf.Log.Filepath)
	log.Printf("  Redis Addr: %s\n", Conf.Redis.Addr)
}

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("mysql", Conf.Mysql.DSN)
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}
	if err = DB.Ping(); err != nil {
		panic("数据库不可用: " + err.Error())
	}
	log.Println("数据库连接成功")

	// 初始化建表
	if err := ensureTableExists(); err != nil {
		panic("建表失败: " + err.Error())
	}
	log.Println("数据库表初始化完成")
}

func ensureTableExists() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS ratelimit_regions (
		id INT AUTO_INCREMENT PRIMARY KEY,
		region_code VARCHAR(32) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	_, err := DB.Exec(createTableSQL)
	return err
}

var (
	RedisClient *redis.Client
	Ctx         = context.Background()
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     Conf.Redis.Addr,
		Password: Conf.Redis.Password,
		DB:       Conf.Redis.DB,
	})

	if err := RedisClient.Ping(Ctx).Err(); err != nil {
		panic("Redis连接失败: " + err.Error())
	}
	log.Println("Redis连接成功")
}
