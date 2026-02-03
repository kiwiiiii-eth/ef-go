package config

import (
	"fmt"
	"os"
	"time"
)

// Config 應用配置結構
type Config struct {
	Database DatabaseConfig
	App      AppConfig
	External ExternalConfig
}

// DatabaseConfig 資料庫配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// AppConfig 應用配置
type AppConfig struct {
	Port     string
	Timezone *time.Location
}

// ExternalConfig 外部API配置
type ExternalConfig struct {
	YihongAPIURL string
	TaipowerURL  string
}

// 場站ID常數
const (
	SiteNorth   = "north"
	SiteCentral = "central"
	SiteSouth   = "south"
)

// Load 載入配置
func Load() *Config {
	// 設置台灣時區 UTC+8
	location, err := time.LoadLocation("Asia/Taipei")
	if err != nil {
		location = time.FixedZone("CST", 8*60*60)
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "vpp_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		App: AppConfig{
			Port:     getEnv("PORT", "8080"),
			Timezone: location,
		},
		External: ExternalConfig{
			YihongAPIURL: getEnv("YIHONG_API_URL", "https://api.yihong-solar.com"),
			TaipowerURL:  getEnv("TAIPOWER_URL", "https://www.taipower.com.tw"),
		},
	}
}

// GetDSN 獲取資料庫連接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// getEnv 獲取環境變數，如果不存在則返回默認值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// IsValidSite 檢查場站ID是否有效
func IsValidSite(siteID string) bool {
	return siteID == SiteNorth || siteID == SiteCentral || siteID == SiteSouth
}
