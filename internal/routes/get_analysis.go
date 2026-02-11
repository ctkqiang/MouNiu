package routes

import (
	"mouniu/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Analysis(router *gin.Engine, db *gorm.DB) {
	public := router.Group(config.API)
	{
		public.GET(config.ANALYSIS_ALL, func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}
}
