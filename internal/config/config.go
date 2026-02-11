package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type QUESTDB_CONFIGURATION struct {
	Database string `json:"database"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     int16  `json:"port"`
}

type MYSQL_CONFIGURATION struct {
	Database string `json:"database"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     int16  `json:"port"`
}

var (
	SINA_URL     = "https://stock.finance.sina.com.cn/hkstock/quotes/"
	SINA_API     = "https://hq.sinajs.cn/list="
	SINA_REFERER = "https://finance.sina.com.cn"

	SINA_ANNNOUNCEMENT_SZ = "https://vip.stock.finance.sina.com.cn/corp/view/vCB_AllMemordDetail.php?stockid="
	SINA_ANNNOUNCEMENT_HK = "https://stock.finance.sina.com.cn/hkstock/notice/"
)

var (
	MYSQL_CONFIG   MYSQL_CONFIGURATION
	QUESTDB_CONFIG QUESTDB_CONFIGURATION
)

func init() {
	// 尝试从当前目录及上级目录查找 .env
	currDir, _ := os.Getwd()
	pathsToTry := []string{
		filepath.Join(currDir, ".env"),
		filepath.Join(currDir, "internal", "config", ".env"),
		filepath.Join(currDir, "..", ".env"),
		filepath.Join(currDir, "..", "internal", "config", ".env"),
	}

	loaded := false
	for _, p := range pathsToTry {
		if _, err := os.Stat(p); err == nil {
			if err := godotenv.Load(p); err == nil {
				loaded = true
				break
			}
		}
	}

	if !loaded {
		godotenv.Load()
	}

	MYSQL_CONFIG = MYSQL_CONFIGURATION{
		Host:     os.Getenv("MYSQL_HOST"),
		Database: os.Getenv("MYSQL_DATABASE"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Port:     3306,
	}

	QUESTDB_CONFIG = QUESTDB_CONFIGURATION{
		Host:     os.Getenv("QUESTDB_HOST"),
		Database: os.Getenv("QUESTDB_DB"),
		User:     os.Getenv("QUESTDB_USER"),
		Password: os.Getenv("QUESTDB_PASS"),
		Port:     8812,
	}
}
