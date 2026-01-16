package config

import (
	"os"
)

type Config struct {
	Port          string
	FeishuAppID   string
	FeishuSecret  string
	FeishuWikiID  string
	FeishuBaseURL string
}

func Load() *Config {
	return &Config{
		Port:          getEnv("PORT", "8080"),
		FeishuAppID:   getEnv("FEISHU_APP_ID", "xxx"),
		FeishuSecret:  getEnv("FEISHU_APP_SECRET", "xxx"),
		FeishuWikiID:  getEnv("FEISHU_WIKI_ID", "xxx"),
		FeishuBaseURL: "https://open.feishu.cn",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
