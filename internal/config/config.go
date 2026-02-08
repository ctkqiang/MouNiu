package config

import (
	"mouniu/internal/utilities"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type MYSQL_CONFIGURATION struct {
	Database string `json:"database"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     int16  `json:"port"`
}

var (
	// https://stock.finance.sina.com.cn/hkstock/quotes/01810.html
	SINA_URL     = "https://stock.finance.sina.com.cn/hkstock/quotes/"
	SINA_API     = "https://hq.sinajs.cn/list="
	SINA_REFERER = "https://finance.sina.com.cn"
)

var (
	MYSQL_CONFIG MYSQL_CONFIGURATION
)

func init() {

	wd, err := os.Getwd()

	if err == nil {
		envPath := filepath.Join(wd, "internal", "config", ".env")

		if _, err := os.Stat(envPath); err == nil {
			if err := godotenv.Load(envPath); err != nil {
				utilities.Log(utilities.ERROR, "[CONFIG] WARNING!!! Failed to load [.env] File %s: %v\n", envPath, err)
			}
		} else {
			godotenv.Load()
		}
	}

	MYSQL_CONFIG = MYSQL_CONFIGURATION{
		Host:     os.Getenv("MYSQL_HOST"),
		Database: os.Getenv("MYSQL_DATABASE"),
		User:     os.Getenv("MYSQL_USER"),
		Password: os.Getenv("MYSQL_PASSWORD"),
		Port:     3306,
	}
}
