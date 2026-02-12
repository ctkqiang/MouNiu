package model

import "time"

// StockIndicator 存储每只股票计算后的各项技术指标结果
type StockIndicator struct {
	Timestamp time.Time `json:"timestamp" gorm:"column:timestamp;type:TIMESTAMP"`
	StockCode string    `json:"stock_code" gorm:"column:stock_code;type:SYMBOL"`
	StockName string    `json:"stock_name" gorm:"column:stock_name;type:SYMBOL"`

	// MACD 指标
	MACD_DIF  float64 `json:"macd_dif" gorm:"column:macd_dif"`
	MACD_DEA  float64 `json:"macd_dea" gorm:"column:macd_dea"`
	MACD_Hist float64 `json:"macd_hist" gorm:"column:macd_hist"`

	// 布林带指标
	BOLL_Upper  float64 `json:"boll_upper" gorm:"column:boll_upper"`
	BOLL_Middle float64 `json:"boll_middle" gorm:"column:boll_middle"`
	BOLL_Lower  float64 `json:"boll_lower" gorm:"column:boll_lower"`

	// 神奇九转指标
	TD9_Count       int  `json:"td9_count" gorm:"column:td9_count"`
	TD9_IsBuySetup  bool `json:"td9_is_buy_setup" gorm:"column:td9_is_buy_setup"`
	TD9_IsSellSetup bool `json:"td9_is_sell_setup" gorm:"column:td9_is_sell_setup"`

	// RSI 指标
	RSI float64 `json:"rsi" gorm:"column:rsi"`

	// 移动平均线
	SMA_20 float64 `json:"sma_20" gorm:"column:sma_20"`
}
