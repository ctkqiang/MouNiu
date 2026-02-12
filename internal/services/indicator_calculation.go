package services

import (
	"bufio"
	"mouniu/internal/database"
	"mouniu/internal/function"
	"mouniu/internal/model"
	"mouniu/internal/utilities"
	"os"
	"strings"
	"time"
)

// CalculateAllIndicators 为所有配置的股票计算技术指标
func CalculateAllIndicators(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		utilities.Error("无法打开符号文件 [%s]: %v", filePath, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		tickerSymbol, _, err := utilities.FormatTicker(line)
		if err != nil {
			utilities.Error("股票代码格式错误: %s", err)
			continue
		}

		// 获取股票名称（从 CandleStickData 表中获取最近的一条记录）
		db, err := database.GetQuestDatabaseConnection()
		if err != nil {
			utilities.Error("数据库连接失败: %v", err)
			continue
		}

		var latestData model.CandleStickData
		if err := db.Where("\"股票代码\" = ?", tickerSymbol).Order("timestamp desc").First(&latestData).Error; err != nil {
			utilities.Error("无法获取股票 %s 的名称: %v", tickerSymbol, err)
			continue
		}

		if err := CalculateAndStoreIndicators(tickerSymbol, latestData.StockName); err != nil {
			utilities.Error("计算 %s 指标出错: %v", tickerSymbol, err)
		}
	}
}

// CalculateAndStoreIndicators 为指定股票计算技术指标并存储到新表中
func CalculateAndStoreIndicators(stockCode, stockName string) error {
	db, err := database.GetQuestDatabaseConnection()
	if err != nil {
		return err
	}

	// 从 QuestDB 获取最近 100 条价格记录，以确保有足够的数据计算指标
	var historicalData []model.CandleStickData
	err = db.Where("股票代码 = ?", stockCode).Order("timestamp desc").Limit(100).Find(&historicalData).Error
	if err != nil {
		return err
	}

	// 如果数据少于 30 条，可能无法计算某些长周期的指标（如 MACD 26, SMA 20 等）
	if len(historicalData) < 30 {
		utilities.Info("股票 %s (%s) 历史数据不足 (%d 条)，跳过指标计算", stockName, stockCode, len(historicalData))
		return nil
	}

	// 将数据按时间顺序排列（从旧到新）
	prices := make([]float64, len(historicalData))
	for i := 0; i < len(historicalData); i++ {
		prices[i] = historicalData[len(historicalData)-1-i].CurrentPrice
	}

	// 1. 计算 MACD (12, 26, 9)
	macd, err := function.GetMACD(prices, 12, 26, 9)
	if err != nil {
		utilities.Error("计算 MACD 失败: %v", err)
	}

	// 2. 计算布林带 (20, 2)
	boll, err := function.GetBollingerBands(prices, 20, 2)
	if err != nil {
		utilities.Error("计算布林带失败: %v", err)
	}

	// 3. 计算神奇九转 (TD9)
	td9, err := function.GetTD9(prices)
	if err != nil {
		utilities.Error("计算神奇九转失败: %v", err)
	}

	// 4. 计算 RSI (14)
	rsiSeries := function.GetRSI(prices, 14)
	latestRSI := 0.0
	if len(rsiSeries) > 0 {
		latestRSI = rsiSeries[len(rsiSeries)-1]
	}

	// 5. 计算 20 日简单移动平均线 (SMA)
	sma20, err := function.GetSimpleMovingAverage(prices, 20)
	if err != nil {
		utilities.Error("计算 SMA20 失败: %v", err)
	}

	// 构造指标结果模型
	indicatorResult := model.StockIndicator{
		Timestamp:       time.Now(),
		StockCode:       stockCode,
		StockName:       stockName,
		MACD_DIF:        macd.DIF,
		MACD_DEA:        macd.DEA,
		MACD_Hist:       macd.MACD,
		BOLL_Upper:      boll.Upper,
		BOLL_Middle:     boll.Middle,
		BOLL_Lower:      boll.Lower,
		TD9_Count:       td9.Count,
		TD9_IsBuySetup:  td9.IsBuySetup,
		TD9_IsSellSetup: td9.IsSellSetup,
		RSI:             latestRSI,
		SMA_20:          sma20,
	}

	// 插入到新表 stock_indicators 中
	if err := db.Create(&indicatorResult).Error; err != nil {
		return err
	}

	utilities.Info("股票 %s (%s) 指标计算完成并存入数据库", stockName, stockCode)
	return nil
}
