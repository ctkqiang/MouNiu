package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllStocks(router *gin.Engine, db *gorm.DB) {
	public := router.Group("/api")
	{
		public.GET("/all", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}
}
