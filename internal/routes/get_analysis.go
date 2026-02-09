package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Analysis(router *gin.Engine, db *gorm.DB) {
	public := router.Group("/api")
	{

		public.GET("/analysis", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}
}
