package services

import (
	"fmt"
	"io"
	"mouniu/internal/config"
	"mouniu/internal/database"
	"mouniu/internal/model"
	"mouniu/internal/utilities"
	"net/http"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gorm.io/gorm"
)

func GetCandleStickData(exchange string, tickerSymbol string) (*model.CandleStickData, error) {
	var candlestickData model.CandleStickData

	database, err := database.GetQuestDatabaseConnection()
	if err != nil {
		return nil, fmt.Errorf("QuestDB %v", err)
	}

	if err := database.AutoMigrate(&model.CandleStickData{}); err != nil {
		utilities.Error("自动迁移 CandleStickData 表失败: %v", err)
		return nil, err
	}

	tickerId := exchange + tickerSymbol
	url := config.SINA_API + tickerId

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utilities.Error("%s", "创建HTTP请求失败: "+err.Error())
		return nil, err
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	request.Header.Set("Referer", string(config.SINA_REFERER))
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	response, err := client.Do(request)
	if err != nil {
		utilities.Error("%s", "请求新浪API失败: "+err.Error())
		return nil, err
	}

	defer response.Body.Close()

	reader := transform.NewReader(response.Body, simplifiedchinese.GBK.NewDecoder())
	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	body := string(bodyBytes)
	if strings.Contains(body, fmt.Sprintf("var hq_str_%s=", tickerId)) {
		start := strings.Index(body, "\"")
		end := strings.LastIndex(body, "\"")

		if start != -1 && end != -1 && end > start {
			dataStr := body[start+1 : end]
			fields := strings.Split(dataStr, ",")

			if len(fields) >= 14 {
				candlestickData.StockName = fields[0]
				candlestickData.CurrentPrice = fields[1]
				candlestickData.PriceChange = fields[2]
				candlestickData.ChangePercentage = fields[3]
				candlestickData.PreviousClose = fields[4]
				candlestickData.TodayOpen = fields[5]
				candlestickData.High = fields[6]
				candlestickData.Low = fields[7]
				candlestickData.Volume = fields[8]
				candlestickData.Turnover = fields[9]
				candlestickData.PERatio = fields[10]
				candlestickData.MarketCapital = fields[11]
				candlestickData.Updatetime = fields[12]
			}
		}
	}

	if err := InsertIntoTable(database, &candlestickData); err != nil {
		utilities.Error("插入数据失败: %v", err)
	}

	return &candlestickData, nil
}

func InsertIntoTable(db *gorm.DB, data *model.CandleStickData) error {
	if data.Timestamp.IsZero() {
		data.Timestamp = time.Now()
	}

	result := db.Create(data)

	if result.Error != nil {
		return fmt.Errorf("插入失败: %v", result.Error)
	}

	utilities.Info("已为 %s 插入 1 行数据\n", data.StockName)

	return nil
}
