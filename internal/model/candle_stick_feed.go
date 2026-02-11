package model

import (
	"encoding/json"
	"time"
)

type CandleStickData struct {
	Timestamp time.Time `json:"timestamp" gorm:"column:timestamp;type:TIMESTAMP"`

	StockCode        string  `json:"股票代码" gorm:"column:股票代码;type:SYMBOL"`
	StockName        string  `json:"股票名称" gorm:"column:股票名称;type:SYMBOL"`
	CurrentPrice     float64 `json:"当前股价" gorm:"column:当前股价"`
	PriceChange      float64 `json:"涨跌额" gorm:"column:涨跌额"`
	ChangePercentage float64 `json:"涨跌幅" gorm:"column:涨跌幅"`
	PreviousClose    float64 `json:"昨收盘" gorm:"column:昨收盘"`
	TodayOpen        float64 `json:"今开盘" gorm:"column:今开盘"`
	High             float64 `json:"最高价" gorm:"column:最高价"`
	Low              float64 `json:"最低价" gorm:"column:最低价"`
	Volume           float64 `json:"成交量" gorm:"column:成交量"`
	Turnover         float64 `json:"成交额" gorm:"column:成交额"`
	PERatio          float64 `json:"市盈率" gorm:"column:市盈率"`
	MarketCapital    float64 `json:"股市值" gorm:"column:股市值"`
	Updatetime       string  `json:"更新时间" gorm:"column:更新时间"`
}

func (candleData *CandleStickData) ToJson() (string, error) {
	jsonData, err := json.Marshal(candleData)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (candleData *CandleStickData) ToJsonPretty() (string, error) {
	jsonData, err := json.MarshalIndent(candleData, "", "    ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (candleData *CandleStickData) ToString() (string, error) {
	jsonData, err := json.Marshal(candleData)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
