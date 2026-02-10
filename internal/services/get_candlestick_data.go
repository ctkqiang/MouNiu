package services

import (
	"bufio"
	"fmt"
	"io"
	"mouniu/internal/config"
	"mouniu/internal/database"
	"mouniu/internal/model"

	"mouniu/internal/utilities"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gorm.io/gorm"
)

func GetCandleStickData(exchange string, tickerSymbol string) (*model.CandleStickData, error) {
	var candlestickData model.CandleStickData

	questDB, err := database.GetQuestDatabaseConnection()
	if err != nil {
		return nil, fmt.Errorf("QuestDB %v", err)
	}

	migrator := questDB.Migrator()
	if !migrator.HasTable(&model.CandleStickData{}) {
		if err := migrator.CreateTable(&model.CandleStickData{}); err != nil {
			utilities.Error("CreateTable Error: %v", err)
			return nil, err
		}
	}

	tickerId := exchange + tickerSymbol
	url := config.SINA_API + tickerId

	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		utilities.Error("%s", "HTTP Request Error: "+err.Error())
		return nil, err
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	request.Header.Set("Referer", string(config.SINA_REFERER))
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	response, err := client.Do(request)
	if err != nil {
		utilities.Error("%s", "Sina API Error: "+err.Error())
		return nil, err
	}

	defer response.Body.Close()

	reader := transform.NewReader(response.Body, simplifiedchinese.GBK.NewDecoder())
	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("Read Body Error: %v", err)
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
				candlestickData.CurrentPrice = utilities.FormatStringToFloat64AndDecimalTo2(fields[1])
				candlestickData.PriceChange = utilities.FormatStringToFloat64AndDecimalTo2(fields[2])
				candlestickData.ChangePercentage = utilities.FormatStringToFloat64AndDecimalTo2(fields[3])
				candlestickData.PreviousClose = utilities.FormatStringToFloat64AndDecimalTo2(fields[4])
				candlestickData.TodayOpen = utilities.FormatStringToFloat64AndDecimalTo2(fields[5])
				candlestickData.High = utilities.FormatStringToFloat64AndDecimalTo2(fields[6])
				candlestickData.Low = utilities.FormatStringToFloat64AndDecimalTo2(fields[7])
				candlestickData.Volume = utilities.FormatStringToFloat64AndDecimalTo2(fields[8])
				candlestickData.Turnover = utilities.FormatStringToFloat64AndDecimalTo2(fields[9])
				candlestickData.PERatio = utilities.FormatStringToFloat64AndDecimalTo2(fields[10])
				candlestickData.MarketCapital = utilities.FormatStringToFloat64AndDecimalTo2(fields[11])
				candlestickData.Updatetime = fields[12]
				candlestickData.Timestamp = time.Now()
			}
		}
	}

	if candlestickData.StockName != "" {
		if err := InsertIntoTable(questDB, &candlestickData); err != nil {
			utilities.Error("Insert Error: %v", err)
		}
	}

	return &candlestickData, nil
}

func InsertIntoTable(db *gorm.DB, data *model.CandleStickData) error {
	if data.Timestamp.IsZero() {
		data.Timestamp = time.Now()
	}
	return db.Create(data).Error
}

func GetStockConcurrently(filePath string) {
	exchange := model.ExchangeHK

	symbolsFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("无法打开符号文件 [%s]: %v\n", filePath, err)
		return
	}
	defer symbolsFile.Close()

	scanner := bufio.NewScanner(symbolsFile)
	for scanner.Scan() {
		ticker := strings.TrimSpace(scanner.Text())
		if ticker == "" {
			continue
		}

		datafeed, err := GetCandleStickData(string(exchange), ticker)
		if err != nil {
			fmt.Printf("抓取 %s 出错: %v\n", ticker, err)
			continue
		}

		fmt.Println(datafeed.ToJson())
	}
}
