package routes

import (
	"mouniu/internal/config"
	"mouniu/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Analysis(router *gin.Engine, db *gorm.DB) {
	public := router.Group(config.API)
	{
		public.GET(config.ANALYSIS_ALL, GetAllAnalysisHandler(db))
		public.GET(config.ANALYSIS_SYMBOL, GetSymbolAnalysisHandler(db))
	}
}

// GetAllAnalysisHandler 获取所有股票的最新指标分析结果
// @Summary 获取所有股票的最新指标分析结果
// @Description 获取数据库中每只股票最新的一条指标计算记录，包括 MACD、布林带、RSI 等
// @Tags 指标分析
// @Produce json
// @Success 200 {array} model.StockIndicator
// @Failure 500 {object} map[string]string
// @Router /analysis [get]
func GetAllAnalysisHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var indicators []model.StockIndicator
		// 使用子查询获取每只股票最新的一条记录
		subQuery := db.Table("stock_indicators").Select("MAX(timestamp)").Group("stock_code")
		if err := db.Where("timestamp IN (?)", subQuery).Find(&indicators).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分析数据失败: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, indicators)
	}
}

// GetSymbolAnalysisHandler 获取特定股票的最新指标分析结果
// @Summary 获取特定股票的最新指标分析结果
// @Description 根据股票代码获取其最近一次计算的技术指标数据
// @Tags 指标分析
// @Param symbol path string true "股票代码 (如 SH600519)"
// @Produce json
// @Success 200 {object} model.StockIndicator
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /analysis/{symbol} [get]
func GetSymbolAnalysisHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		symbol := c.Param("symbol")
		var indicator model.StockIndicator
		if err := db.Where("stock_code = ?", symbol).Order("timestamp desc").First(&indicator).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "未找到该股票的分析数据"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败: " + err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, indicator)
	}
}
