package services

import (
	"fmt"
	"io"
	"mouniu/internal/config"
	"mouniu/internal/model"
	"mouniu/internal/utilities"
	"net/http"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func GetCandleStickData(exchange string, tickerSymbol string) (*model.CandleStickData, error) {
	var candlestickData model.CandleStickData

	url := config.SINA_API + exchange + tickerSymbol
	client := &http.Client{Timeout: 10 * time.Second}
	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		utilities.Error("%s", "创建HTTP请求失败: "+err.Error())
		return nil, err
	}

	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	request.Header.Set("Referer", "https://finance.sina.com.cn")
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
	if strings.Contains(body, fmt.Sprintf("var hq_str_%s=", tickerSymbol)) {
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
				candlestickData.Volume = fields[8] + "股"
				candlestickData.Turnover = fields[9] + "元"
				candlestickData.PERatio = fields[10]
				candlestickData.MarketCapital = fields[11]
				candlestickData.Updatetime = fields[12]
			}

		}
	}

	return &candlestickData, nil
}
