package main

import (
	"fmt"
	"log"
	"mouniu/internal/crons"
	"mouniu/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"

	cron_v3 "github.com/robfig/cron/v3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	SymbolsFile = "internal/config/symbols.txt"
)

func main() {
	router := gin.Default()
	cronManager := cron_v3.New()

	crons.RunStockUpdate(cronManager, SymbolsFile)

	cronManager.Start()
	defer cronManager.Stop()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"消息": "pong",
		})
	})

	services.GetStockConcurrently(SymbolsFile)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	if err := router.Run(fmt.Sprintf(":%d", 8080)); err != nil {
		log.Fatalf("HTTP 服务启动失败: %v", err)
	}
}
