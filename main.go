package main

import (
	"fmt"
	"log"
	_ "mouniu/docs" // 导入生成的 swagger 文档
	"mouniu/internal/crons"
	"mouniu/internal/database"
	"mouniu/internal/routes"
	"mouniu/internal/utilities"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	cron_v3 "github.com/robfig/cron/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	SymbolsFile = "internal/config/symbols.txt"
	Addr        = fmt.Sprintf(":%d", getEnvAsInt("APP_PORT", 8000))
	Port        = getEnvAsInt("APP_PORT", 8000)
)

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为整数，如果不存在或转换失败则返回默认值
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// @title 牟牛 (MouNiu) 股票分析系统 API
// @version 1.0
// @description 这是一个专业的股票数据抓取与技术指标分析系统。支持 MACD、布林带、神奇九转等多种指标。
// @termsOfService http://swagger.io/terms/

// @contact.name 钟智强
// @contact.url https://github.com/ctkqiang

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
func main() {
	gin.SetMode(gin.DebugMode)

	router := gin.Default()

	router.Use(gin.Recovery())

	cronManager := cron_v3.New()
	crons.RunStockUpdate(cronManager, SymbolsFile)
	crons.RunIndicatorUpdate(cronManager, SymbolsFile)
	crons.RunAnnouncementUpdate(cronManager, SymbolsFile)
	cronManager.Start()
	defer cronManager.Stop()

	database, err := database.GetQuestDatabaseConnection()
	if err != nil {
		utilities.Error("数据库连接失败: %v", err)
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"消息": "pong",
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	routes.GetAllStocks(router, database)
	routes.GetAnnouncements(router, database)
	routes.Analysis(router, database)

	router.Run(Addr)

	if err := router.Run(fmt.Sprintf(":%d", Port)); err != nil {
		log.Fatalf("HTTP 服务启动失败: %v", err)
	}
}
