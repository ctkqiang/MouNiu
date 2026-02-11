package routes

import (
	"mouniu/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	announcements []model.Announcement
)

func GetAnnouncements(router *gin.Engine, db *gorm.DB) {
	public := router.Group("/api")
	{
		public.GET("/announcement/all", func(c *gin.Context) {

			result := db.Order("publish_date DESC").Find(&announcements)

			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "无法获取公告数据: " + result.Error.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, announcements)
		})

		public.GET("/announcement/:symbol", func(c *gin.Context) {
			symbol := c.Param("symbol")
			result := db.Where("stock_code LIKE ?", "%"+symbol+"%").Order("publish_date DESC").Find(&announcements)

			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "无法获取指定股票的公告数据: " + result.Error.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, announcements)
		})
	}
}
