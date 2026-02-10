package crons

import (
	"mouniu/internal/services"

	"github.com/robfig/cron/v3"
)

func RunStockUpdate(c *cron.Cron, filePath string) {
	c.AddFunc("* * * * *", func() {
		services.GetStockConcurrently(filePath)
	})
}
