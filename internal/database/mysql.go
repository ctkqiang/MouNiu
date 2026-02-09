package database

import (
	"fmt"
	"mouniu/internal/config"
	"mouniu/internal/model"
	"mouniu/internal/utilities"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	database, err := GetMYSQLConnection()
	if err != nil {
		utilities.Info("初始化数据库连接失败, %v", err)
		return
	}

	if err := database.AutoMigrate(&model.CandleStickData{}); err != nil {
		utilities.Error("自动迁移 CandleStickData 表失败: %v", err)
	}
}

// GetMYSQLConnection 根据全局配置 MYSQL_CONFIG 中的参数，建立并返回一个 GORM 的 MySQL 数据库连接实例。
// 该函数会依次校验用户名、密码、主机地址、端口和数据库名是否为空或无效；
// 若校验通过，则使用 DSN（Data Source Name）格式拼接连接字符串，并通过 gorm.Open 打开数据库连接。
// 连接成功后会自动执行 AutoMigrate，确保 model.CandleStickData 表结构已同步；
// 若任意环节出错，将记录错误日志并返回对应的错误信息。
func GetMYSQLConnection() (*gorm.DB, error) {
	if config.MYSQL_CONFIG.User == "" {
		return nil, fmt.Errorf("MySQL 用户名为空")
	}

	if config.MYSQL_CONFIG.Password == "" {
		return nil, fmt.Errorf("MySQL 密码为空")
	}

	if config.MYSQL_CONFIG.Host == "" {
		return nil, fmt.Errorf("MySQL 主机地址为空")
	}

	if config.MYSQL_CONFIG.Port == 0 {
		return nil, fmt.Errorf("MySQL 端口无效")
	}

	if config.MYSQL_CONFIG.Database == "" {
		return nil, fmt.Errorf("MySQL 数据库名为空")
	}

	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%v:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.MYSQL_CONFIG.User,
		config.MYSQL_CONFIG.Password,
		config.MYSQL_CONFIG.Host,
		config.MYSQL_CONFIG.Port,
		config.MYSQL_CONFIG.Database,
	)

	database, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		utilities.Error("%s", "连接MySQL数据库失败："+err.Error())
		return nil, err
	}

	utilities.Info("连接MySQL数据库成功")

	return database, nil
}
