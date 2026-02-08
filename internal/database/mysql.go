package database

import (
	"fmt"
	"mouniu/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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

	_, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{})
	if err != nil {
		panic("连接 MySQL 数据库失败：" + err.Error())
	}

	return nil, nil
}
