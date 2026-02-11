package main

import (
	"fmt"
	"log"
	"mouniu/internal/crons"
	"mouniu/internal/database"
	"mouniu/internal/routes"
	"mouniu/internal/utilities"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	cron_v3 "github.com/robfig/cron/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	SymbolsFile = "internal/config/symbols.txt"
	Addr        = ":8080"
	Port        = 8080
)

func main() {
	gin.SetMode(gin.DebugMode)

	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[VERBOSE] %s | %d | %t | %s | %s | %s | %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Latency < 500*time.Millisecond,
			param.ClientIP,
			param.Method,
			param.Path,
			param.ErrorMessage,
		)
	}))

	cronManager := cron_v3.New()

	crons.RunStockUpdate(cronManager, SymbolsFile)
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
	routes.Analysis(router, database)

	router.Run(Addr)

	if err := router.Run(fmt.Sprintf(":%d", Port)); err != nil {
		log.Fatalf("HTTP 服务启动失败: %v", err)
	}
}
