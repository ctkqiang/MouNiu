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
		public.GET(config.ANNOUNCEMENT_ALL, func(c *gin.Context) {
			var announcements []model.Announcement

			// 获取分页参数
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

			if page < 1 {
				page = 1
			}
			if pageSize < 1 || pageSize > 100 {
				pageSize = 20
			}

			offset := (page - 1) * pageSize

			result := db.Order("publish_date DESC").
				Limit(pageSize).
				Offset(offset).
				Find(&announcements)

			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "无法获取公告数据: " + result.Error.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"data":     announcements,
				"page":     page,
				"pageSize": pageSize,
			})
		})

		public.GET(config.ANNOUNCEMENT_SYMBOLS, func(c *gin.Context) {
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

			result := db.Where("stock_code LIKE ?", "%"+symbol+"%").
				Order("publish_date DESC").
				Limit(pageSize).
				Offset(offset).
				Find(&announcements)

			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "无法获取指定股票的公告数据: " + result.Error.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"data":     announcements,
				"page":     page,
				"pageSize": pageSize,
			})
		})
	}
}
