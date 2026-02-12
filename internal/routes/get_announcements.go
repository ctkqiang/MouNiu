package routes

import (
	"mouniu/internal/config"
	"mouniu/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAnnouncements(router *gin.Engine, db *gorm.DB) {
	public := router.Group(config.API)
	{
		public.GET(config.ANNOUNCEMENT_ALL, GetAllAnnouncementsHandler(db))
		public.GET(config.ANNOUNCEMENT_SYMBOLS, GetAnnouncementsBySymbolHandler(db))
	}
}

// GetAllAnnouncementsHandler 获取所有股票公告
// @Summary 获取所有股票公告
// @Description 分页获取所有上市公司的最新公告信息
// @Tags 公告
// @Param page query int false "页码 (默认 1)"
// @Param pageSize query int false "每页数量 (默认 20, 最大 100)"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /announcement/all [get]
func GetAllAnnouncementsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var announcements []model.Announcement
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 20
		}
		offset := (page - 1) * pageSize
		result := db.Order("publish_date DESC").Limit(pageSize).Offset(offset).Find(&announcements)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取公告数据: " + result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data":     announcements,
			"page":     page,
			"pageSize": pageSize,
		})
	}
}

// GetAnnouncementsBySymbolHandler 获取特定股票的公告
// @Summary 获取特定股票的公告
// @Description 根据股票代码模糊查询其相关的公告信息
// @Tags 公告
// @Param symbol path string true "股票代码 (如 600519)"
// @Param page query int false "页码 (默认 1)"
// @Param pageSize query int false "每页数量 (默认 20, 最大 100)"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /announcement/{symbol} [get]
func GetAnnouncementsBySymbolHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		symbol := c.Param("symbol")
		var announcements []model.Announcement
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
		if page < 1 {
			page = 1
		}
		if pageSize < 1 || pageSize > 100 {
			pageSize = 20
		}
		offset := (page - 1) * pageSize
		result := db.Where("stock_code LIKE ?", "%"+symbol+"%").Order("publish_date DESC").Limit(pageSize).Offset(offset).Find(&announcements)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法获取指定股票的公告数据: " + result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"data":     announcements,
			"page":     page,
			"pageSize": pageSize,
		})
	}
}
