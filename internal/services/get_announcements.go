package services

import (
	"crypto/md5"
	"fmt"
	"io"
	"mouniu/internal/config"
	"mouniu/internal/database"
	"mouniu/internal/model"
	"mouniu/internal/utilities"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gorm.io/gorm"
)

func gbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(strings.NewReader(string(s)), simplifiedchinese.GBK.NewDecoder())
	d, err := io.ReadAll(reader)

	if err != nil {
		return nil, err
	}

	return d, nil
}

func generateDeterministicID(stockCode, title, dateStr string) string {
	data := fmt.Sprintf("%s-%s-%s", stockCode, title, dateStr)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

func parseDate(dateStr string) time.Time {
	dateStr = strings.ReplaceAll(dateStr, "公告日期:", "")
	dateStr = strings.TrimSpace(dateStr)

	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}

	currentYear := time.Now().Year()
	shortFormats := []string{"01-02", "01/02"}

	for _, format := range shortFormats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return time.Date(currentYear, t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
		}
	}

	return time.Time{}
}

func storeAnnouncement(announcement model.Announcement) error {
	db, err := database.GetQuestDatabaseConnection()
	if err != nil {
		return fmt.Errorf("无法连接到数据库: %v", err)
	}

	var existing model.Announcement
	err = db.Where("announcement_id = ?", announcement.AnnouncementID).First(&existing).Error

	switch err {
	case nil:
		if updateErr := db.Model(&existing).Where("announcement_id = ?", existing.AnnouncementID).Omit("publish_date", "created_at").Updates(announcement).Error; updateErr != nil {
			return fmt.Errorf("更新公告失败: %v", updateErr)
		}
	case gorm.ErrRecordNotFound:
		if createErr := db.Create(&announcement).Error; createErr != nil {
			return fmt.Errorf("插入公告失败: %v", createErr)
		}
	default:
		return fmt.Errorf("检查公告是否存在时出错: %v", err)
	}

	return nil
}

func GetShenZhenAnnouncement(tickerSymbol string) {
	url := config.SINA_ANNNOUNCEMENT_SZ + tickerSymbol

	utilities.Info("正在爬取深圳股票 %s 的公告...", tickerSymbol)
	utilities.Debug("目标 URL: %s", url)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		utilities.Error("获取公告失败: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		utilities.Error("请求失败，状态码: %d", resp.StatusCode)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		utilities.Error("解析 HTML 失败: %v", err)
		return
	}

	count := 0
	doc.Find("table.list_table tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}

		titleLink := s.Find("th a").First()
		title := strings.TrimSpace(titleLink.Text())
		url, _ := titleLink.Attr("href")
		dateStr := s.Find("td").Last().Text()

		if title != "" && url != "" {
			announcement := model.Announcement{
				AnnouncementID: generateDeterministicID(tickerSymbol, title, dateStr),
				StockCode:      tickerSymbol,
				Title:          title,
				PublishDate:    parseDate(dateStr),
				ContentURL:     url,
				Exchange:       "SZ",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			if err := storeAnnouncement(announcement); err != nil {
				utilities.Error("存储公告失败: %v", err)
			} else {
				utilities.Info("成功存储公告: %s", title)
				count++
			}
		}
	})

	utilities.Info("共爬取并处理了 %d 条深圳公告。", count)
}

func GetHongKongAnnouncement(tickerSymbol string) {
	url := config.SINA_ANNNOUNCEMENT_HK + tickerSymbol + ".html"

	utilities.Info("正在爬取香港股票 %s 的公告...", tickerSymbol)
	utilities.Debug("目标 URL: %s", url)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		utilities.Error("获取公告失败: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		utilities.Error("请求失败，状态码: %d", resp.StatusCode)
		return
	}

	// 转换为 UTF-8
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		utilities.Error("读取响应体失败: %v", err)
		return
	}

	utf8Body, err := gbkToUtf8(bodyBytes)
	if err != nil {
		utilities.Error("编码转换失败: %v", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(utf8Body)))
	if err != nil {
		utilities.Error("解析 HTML 失败: %v", err)
		return
	}

	count := 0
	doc.Find("ul.list01 li").Each(func(i int, s *goquery.Selection) {
		titleLink := s.Find("a").First()
		title := strings.TrimSpace(titleLink.Text())
		url, _ := titleLink.Attr("href")
		dateStr := s.Find("span.rt").Text()

		if title != "" && url != "" {
			announcement := model.Announcement{
				AnnouncementID: generateDeterministicID(tickerSymbol, title, dateStr),
				StockCode:      tickerSymbol,
				Title:          title,
				PublishDate:    parseDate(dateStr),
				ContentURL:     strings.ReplaceAll(url, "01810", tickerSymbol),
				Exchange:       "HK",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			}

			if err := storeAnnouncement(announcement); err != nil {
				utilities.Error("存储公告失败: %v", err)
			} else {
				utilities.Info("成功存储公告: %s", title)
				count++
			}
		}
	})

	utilities.Info("共爬取并处理了 %d 条香港公告。", count)
}
