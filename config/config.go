package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config はアプリ全体の設定（Laravelの config/*.php 相当）
type Config struct {
	AppPort   string `envconfig:"APP_PORT" default:"8080"`
	DBPath    string `envconfig:"DB_PATH" default:"database.sqlite"`
	RedisAddr string `envconfig:"REDIS_ADDR" default:"127.0.0.1:6379"`
}

// Load は .env を読み込んだ上で環境変数から設定を組み立てる
func Load() (*Config, error) {
	_ = godotenv.Load() // .env が無くてもエラーにしない
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, err
	}
	return &c, nil
}
