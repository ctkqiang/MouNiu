package database

import (
	"fmt"
	"mouniu/internal/config"
	"mouniu/internal/model"
	"mouniu/internal/utilities"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	if !questDB.Migrator().HasTable(&model.CandleStickData{}) {
		err := questDB.Migrator().CreateTable(&model.CandleStickData{})
		if err != nil {
			return nil, fmt.Errorf("初始化 QuestDB 表失败: %v", err)
		}

		utilities.Info("QuestDB 表 'candle_stick_data' 创建成功")
	}

	utilities.Info("连接QuestDB成功 |> %s", dataSourceName)

	return questDB, nil
}
