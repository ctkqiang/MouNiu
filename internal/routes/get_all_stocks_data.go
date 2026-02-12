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
		public.GET(config.STOCKS_ALL, GetAllStocksHandler(db))
		public.GET(config.STOCKS_SYMBOL, GetStocksByTickerHandler(db))
		public.GET(config.STOCKS_SYMBOL_CURRENT_PRICE, GetCurrentPriceByTickerHandler(db))
	}
}

// GetAllStocksHandler 获取所有股票的历史价格数据
// @Summary 获取所有股票的历史价格数据
// @Description 从数据库中检索所有股票的历史 K 线数据，按时间倒序排列
// @Tags 股票数据
// @Produce json
// @Success 200 {array} model.CandleStickData
// @Failure 500 {object} map[string]string
// @Router /stocks/all [get]
func GetAllStocksHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var stocks []model.CandleStickData
		result := db.Order("timestamp DESC").Find(&stocks)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取数据: " + result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, stocks)
	}
}

// GetStocksByTickerHandler 获取特定股票的历史价格数据
// @Summary 获取特定股票的历史价格数据
// @Description 根据股票代码获取其所有的历史 K 线数据记录
// @Tags 股票数据
// @Param ticker path string true "股票代码 (如 SH600519)"
// @Produce json
// @Success 200 {array} model.CandleStickData
// @Failure 500 {object} map[string]string
// @Router /{ticker} [get]
func GetStocksByTickerHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ticker := c.Param("ticker")
		var stocks []model.CandleStickData
		result := db.Where("\"股票代码\" = ?", ticker).Order("timestamp DESC").Find(&stocks)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取数据: " + result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, stocks)
	}
}

// GetCurrentPriceByTickerHandler 获取特定股票的当前实时价格
// @Summary 获取特定股票的当前实时价格
// @Description 获取指定股票最新的一条价格记录
// @Tags 股票数据
// @Param ticker path string true "股票代码 (如 SH600519)"
// @Produce json
// @Success 200 {array} model.CandleStickData
// @Failure 500 {object} map[string]string
// @Router /current/{ticker} [get]
func GetCurrentPriceByTickerHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		ticker := c.Param("ticker")
		var stocks []model.CandleStickData
		result := db.Where("\"股票代码\" = ?", ticker).Order("timestamp DESC").Limit(1).Find(&stocks)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取数据: " + result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, stocks)
	}
}
