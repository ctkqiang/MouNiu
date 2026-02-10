package routes

import (
	"mouniu/internal/model" // 确保导入了你的模型包
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllStocks(router *gin.Engine, db *gorm.DB) {
	public := router.Group("/api")
	{
		public.GET("/all", func(c *gin.Context) {
			var stocks []model.CandleStickData

			result := db.Find(&stocks)

			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "无法获取数据: " + result.Error.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, stocks)
		})
	}
}
