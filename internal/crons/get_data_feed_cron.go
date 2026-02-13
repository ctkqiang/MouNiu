package crons

import (
	"bufio"
	"mouniu/internal/services"
	"mouniu/internal/utilities"
	"os"
	"strings"

	"github.com/robfig/cron/v3"
)

func RunStockUpdate(c *cron.Cron, filePath string) {
	c.AddFunc("* * * * *", func() {
		services.GetStockConcurrently(filePath)
	})
}

func RunIndicatorUpdate(c *cron.Cron, filePath string) {
	c.AddFunc("*/5 * * * *", func() {
		services.CalculateAllIndicators(filePath)
	})
}

func RunAnnouncementUpdate(c *cron.Cron, filePath string) {
	c.AddFunc("0 */3 * * *", func() {
		file, err := os.Open(filePath)
		if err != nil {
			utilities.Error("无法打开符号文件 [%s]: %v", filePath, err)
			return
		}

		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasSuffix(line, ".SZ") {
				ticker, _, _ := utilities.FormatTicker(line)
				services.GetShenZhenAnnouncement(ticker)
			}

			if strings.HasSuffix(line, ".HK") {
				ticker, _, _ := utilities.FormatTicker(line)
				services.GetHongKongAnnouncement(ticker)
			}
		}
	})
}
