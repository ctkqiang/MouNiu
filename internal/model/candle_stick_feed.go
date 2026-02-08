package model

import (
	"encoding/json"

	"gorm.io/gorm"
)

type CandleStickData struct {
	gorm.Model
	StockName        string `json:"股票名称" gorm:"column:股票名称;index"`
	CurrentPrice     string `json:"当前股价" gorm:"column:当前股价"`
	PriceChange      string `json:"涨跌额" gorm:"column:涨跌额"`
	ChangePercentage string `json:"涨跌幅" gorm:"column:涨跌幅"`
	PreviousClose    string `json:"昨收盘" gorm:"column:昨收盘"`
	TodayOpen        string `json:"今开盘" gorm:"column:今开盘"`
	High             string `json:"最高价" gorm:"column:最高价"`
	Low              string `json:"最低价" gorm:"column:最低价"`
	Volume           string `json:"成交量" gorm:"column:成交量"`
	Turnover         string `json:"成交额" gorm:"column:成交额"`
	PERatio          string `json:"市盈率" gorm:"column:市盈率"`
	MarketCapital    string `json:"股市值" gorm:"column:股市值"`
	Updatetime       string `json:"更新时间" gorm:"column:更新时间"`
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
