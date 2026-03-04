package database

import (
	"fmt"
	"mouniu/internal/config"
	"mouniu/internal/model"
	"mouniu/internal/utilities"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	initTableOnce sync.Once
)

func GetQuestDatabaseConnection() (*gorm.DB, error) {
	dataSourceName := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=8812 sslmode=disable",
		config.QUESTDB_CONFIG.Host,
		config.QUESTDB_CONFIG.User,
		config.QUESTDB_CONFIG.Password,
		config.QUESTDB_CONFIG.Database,
	)

	var questDB *gorm.DB
	var err error
	maxRetries := 10
	retryInterval := 5 * time.Second

	// 尝试连接QuestDB，失败时重试
	for i := 0; i < maxRetries; i++ {
		questDB, err = gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})

		if err == nil {
			break
		}

		utilities.Log(utilities.ERROR, "连接QuestDB失败，%d秒后重试: %v", retryInterval/time.Second, err)
		time.Sleep(retryInterval)
	}

	if err != nil {
		return nil, fmt.Errorf("连接QuestDB失败 |> %s", dataSourceName)
	}

	initTableOnce.Do(func() {
		utilities.Info("正在初始化 QuestDB 表...")
		// 初始化 CandleStickData 表
		if !questDB.Migrator().HasTable("candle_stick_data") {
			if err := questDB.Set("gorm:table_options", "timestamp(timestamp) PARTITION BY DAY WAL").Migrator().CreateTable(&model.CandleStickData{}); err != nil {
				utilities.Error("初始化 QuestDB 表 'candle_stick_data' 失败: %v", err)
			} else {
				utilities.Info("QuestDB 表 'candle_stick_data' 创建成功")
			}
		}

		// 初始化 Announcement 表
		if !questDB.Migrator().HasTable("announcements") {
			if err := questDB.Set("gorm:table_options", "timestamp(publish_date) PARTITION BY DAY WAL").Migrator().CreateTable(&model.Announcement{}); err != nil {
				utilities.Error("初始化 QuestDB 表 'announcements' 失败: %v", err)
			} else {
				utilities.Info("QuestDB 表 'announcements' 创建成功")
			}
		}

		// 初始化 StockIndicator 表
		if !questDB.Migrator().HasTable("stock_indicators") {
			if err := questDB.Set("gorm:table_options", "timestamp(timestamp) PARTITION BY DAY WAL").Migrator().CreateTable(&model.StockIndicator{}); err != nil {
				utilities.Error("初始化 QuestDB 表 'stock_indicators' 失败: %v", err)
			} else {
				utilities.Info("QuestDB 表 'stock_indicators' 创建成功")
			}
		}
	})

	utilities.Info("连接QuestDB成功 |> %s", dataSourceName)

	return questDB, nil
}
