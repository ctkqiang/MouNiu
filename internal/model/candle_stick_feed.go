package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

type CandleStickData struct {
	gorm.Model
	StockName        string `json:"股票名称"` // 股票名称
	CurrentPrice     string `json:"当前股价"` // 当前股价 [C]
	PriceChange      string `json:"涨跌额"`  // 涨跌额
	ChangePercentage string `json:"涨跌幅"`  // 涨跌幅
	PreviousClose    string `json:"昨收盘"`  // 昨收盘
	TodayOpen        string `json:"今开盘"`  // 今开盘 [O]
	High             string `json:"最高价"`  // 最高价 [H]
	Low              string `json:"最低价"`  // 最低价 [L]
	Volume           string `json:"成交量"`  // 成交量 [V]
	Turnover         string `json:"成交额"`  // 成交额
	PERatio          string `json:"市盈率"`  // 市盈率
	MarketCapital    string `json:"股市值"`  // 股市值
	Updatetime       string `json:"更新时间"` // 更新时间
}

func (c *CandleStickData) ToJson() (string, error) {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (c *CandleStickData) ToJsonPretty() (string, error) {
	jsonData, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (c *CandleStickData) ToString() (string, error) {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
