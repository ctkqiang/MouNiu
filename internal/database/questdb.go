package database

import (
	"fmt"
	"mouniu/internal/config"
	"mouniu/internal/model"
	"mouniu/internal/utilities"
	"sync"

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

	questDB, err := gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return nil, fmt.Errorf("连接QuestDB失败 |> %s", dataSourceName)
	}

	initTableOnce.Do(func() {
		if !questDB.Migrator().HasTable(&model.CandleStickData{}) {
			err := questDB.Migrator().CreateTable(&model.CandleStickData{})
			if err != nil {
				utilities.Error("初始化 QuestDB 表失败: %v", err)
				// 注意：这里不返回错误，因为连接已经建立，只是表创建失败
				// 应用可能仍然可以运行，只是无法存储数据
			} else {
				utilities.Info("QuestDB 表 'candle_stick_data' 创建成功")
			}
		}
	})

	utilities.Info("连接QuestDB成功 |> %s", dataSourceName)

	return questDB, nil
}
