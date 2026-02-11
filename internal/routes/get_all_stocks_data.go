package routes

import (
	"mouniu/internal/config"
	"mouniu/internal/model" // 确保导入了你的模型包
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllStocks(router *gin.Engine, db *gorm.DB) {
	public := router.Group(config.API)
	{
		public.GET(config.STOCKS_ALL, func(c *gin.Context) {
			var stocks []model.CandleStickData

			result := db.Order("timestamp DESC").Find(&stocks)

			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "无法获取数据: " + result.Error.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, stocks)
		})

		public.GET(config.STOCKS_SYMBOL, func(c *gin.Context) {
			ticker := c.Param("ticker")
			var stocks []model.CandleStickData

			result := db.Where("股票代码 = ?", ticker).Order("timestamp DESC").Find(&stocks)

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
